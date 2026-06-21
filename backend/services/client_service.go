package services

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"latestpack/models"
	"latestpack/repository"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type ClientService struct {
	versionRepo       *repository.VersionRepo
	versionChangeRepo *repository.VersionChangeRepo
	archivesDir       string
}

func NewClientService(versionRepo *repository.VersionRepo, versionChangeRepo *repository.VersionChangeRepo, archivesDir string) *ClientService {
	return &ClientService{
		versionRepo:       versionRepo,
		versionChangeRepo: versionChangeRepo,
		archivesDir:       archivesDir,
	}
}

func (s *ClientService) GetLatest(ctx context.Context) (*models.ClientLatestResponse, error) {
	v, err := s.versionRepo.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return &models.ClientLatestResponse{
		Version:   v.Version,
		Timestamp: v.CreatedAt.UnixMilli(),
	}, nil
}

func (s *ClientService) GetUpdates(ctx context.Context, from string) (*models.ClientUpdatesResponse, error) {
	var allChanges []models.VersionChange
	var err error

	if from == "" {
		allChanges, err = s.versionChangeRepo.ListAll(ctx)
	} else {
		// If the from version doesn't exist in DB, fall back to returning all changes
		existing, findErr := s.versionRepo.FindByVersion(ctx, from)
		if findErr != nil || existing == nil {
			allChanges, err = s.versionChangeRepo.ListAll(ctx)
		} else {
			allChanges, err = s.versionChangeRepo.ListAfterVersion(ctx, from)
		}
	}
	if err != nil {
		return nil, err
	}

	versionMap := make(map[string][]models.ChangeEntry)
	var versionOrder []string

	for _, c := range allChanges {
		if _, exists := versionMap[c.Version]; !exists {
			versionOrder = append(versionOrder, c.Version)
		}
		versionMap[c.Version] = append(versionMap[c.Version], models.ChangeEntry{
			Action: c.Action,
			Path:   c.Path,
		})
	}

	versions := make([]models.ClientUpdateVersion, 0, len(versionOrder))
	for _, ver := range versionOrder {
		v, err := s.versionRepo.FindByVersion(ctx, ver)
		if err != nil || v == nil {
			continue
		}
		versions = append(versions, models.ClientUpdateVersion{
			Version:   ver,
			Timestamp: v.CreatedAt.UnixMilli(),
			Changes:   versionMap[ver],
		})
	}

	return &models.ClientUpdatesResponse{Versions: versions}, nil
}

// findInTar opens a tar archive and returns a reader for the specified file entry.
// Returns the tar header (for size/mod time) and a reader positioned at the file content.
func findInTar(archivePath, targetPath string) (*tar.Header, *tar.Reader, *os.File, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, nil, nil, err
	}

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			f.Close()
			return nil, nil, nil, ErrFileNotFound
		}
		if err != nil {
			f.Close()
			return nil, nil, nil, err
		}
		if hdr.Name == targetPath {
			return hdr, tr, f, nil
		}
	}
}

func (s *ClientService) DownloadFile(ctx context.Context, version string, relPath string, w http.ResponseWriter, r *http.Request) error {
	archivePath := filepath.Join(s.archivesDir, version+".tar")
	hdr, tr, f, err := findInTar(archivePath, relPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fileSize := hdr.Size

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	w.Header().Set("Accept-Ranges", "bytes")

	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		return s.serveRangeFromTar(hdr, tr, f, fileSize, rangeHeader, w)
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.CopyN(w, tr, fileSize)
	return err
}

func (s *ClientService) serveRangeFromTar(hdr *tar.Header, tr *tar.Reader, f *os.File, fileSize int64, rangeHeader string, w http.ResponseWriter) error {
	// We need to re-read from the tar to support Range.
	// The tar reader is sequential, so for range requests we reopen the tar,
	// seek to the file entry, then skip to the range start.
	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return errors.New("invalid range")
	}

	var start, end int64
	var err error

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("invalid range start")
	}

	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return errors.New("invalid range end")
		}
	}

	if start >= fileSize || end >= fileSize || start > end {
		w.Header().Set("Content-Range", "bytes */"+strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	contentLength := end - start + 1

	// The tar reader is already positioned at the file content.
	// Skip 'start' bytes, then copy contentLength bytes.
	if _, err := io.CopyN(io.Discard, tr, start); err != nil {
		return err
	}

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.WriteHeader(http.StatusPartialContent)

	_, err = io.CopyN(w, tr, contentLength)
	return err
}

func (s *ClientService) HeadFile(ctx context.Context, version string, relPath string, w http.ResponseWriter) error {
	archivePath := filepath.Join(s.archivesDir, version+".tar")
	hdr, _, f, err := findInTar(archivePath, relPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(hdr.Size, 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
	return nil
}
