package services

import (
	"archive/tar"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"latestpack/models"
	"latestpack/repository"
)

var (
	ErrVersionExists    = errors.New("version already exists")
	ErrVersionNotFound  = errors.New("version not found")
	ErrNotLatestVersion = errors.New("only the latest version can be deleted")
)

type VersionService struct {
	db            *sql.DB
	repo          *repository.VersionRepo
	versionChange *repository.VersionChangeRepo
	snapshotRepo  *repository.VersionSnapshotRepo
	fileSvc       *FileService
	archivesDir   string
}

func NewVersionService(
	db *sql.DB,
	repo *repository.VersionRepo,
	versionChange *repository.VersionChangeRepo,
	snapshotRepo *repository.VersionSnapshotRepo,
	fileSvc *FileService,
	archivesDir string,
) *VersionService {
	os.MkdirAll(archivesDir, 0755)
	return &VersionService{
		db:            db,
		repo:          repo,
		versionChange: versionChange,
		snapshotRepo:  snapshotRepo,
		fileSvc:       fileSvc,
		archivesDir:   archivesDir,
	}
}

func (s *VersionService) List(ctx context.Context, page, pageSize int) (*models.VersionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}

	versions, total, err := s.repo.List(ctx, int64(page), int64(pageSize))
	if err != nil {
		return nil, err
	}

	items := make([]models.VersionDTO, 0, len(versions))
	for _, v := range versions {
		var notes []string
		if err := json.Unmarshal([]byte(v.Notes), &notes); err != nil {
			notes = []string{}
		}
		items = append(items, models.VersionDTO{
			Version: v.Version,
			Date:    v.Date,
			Size:    v.Size,
			Notes:   notes,
		})
	}

	return &models.VersionListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *VersionService) GetLatest(ctx context.Context) (*models.VersionDTO, error) {
	v, err := s.repo.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	var notes []string
	if err := json.Unmarshal([]byte(v.Notes), &notes); err != nil {
		notes = []string{}
	}
	return &models.VersionDTO{
		Version: v.Version,
		Date:    v.Date,
		Size:    v.Size,
		Notes:   notes,
	}, nil
}

func (s *VersionService) Create(ctx context.Context, req *models.CreateVersionRequest) (*models.VersionDTO, error) {
	existing, err := s.repo.FindByVersion(ctx, req.Version)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrVersionExists
	}

	// Walk filesystem and compute hashes (outside transaction)
	currFiles, err := s.fileSvc.WalkAndHash()
	if err != nil {
		return nil, fmt.Errorf("scan files: %w", err)
	}

	notesJSON, err := json.Marshal(req.Notes)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	v := &models.Version{
		Version:   req.Version,
		Date:      now.Format("2006-01-02"),
		Size:      "0 MB",
		Notes:     string(notesJSON),
		CreatedAt: now,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Get previous version's snapshot BEFORE inserting the new version row
	prevSnapshots, err := s.snapshotRepo.GetLatest(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("get previous snapshot: %w", err)
	}

	// Insert new version row
	if err := s.repo.CreateTx(ctx, tx, v); err != nil {
		return nil, fmt.Errorf("insert version: %w", err)
	}

	// Compute diff
	changes := computeChanges(prevSnapshots, currFiles)

	// Build and insert new snapshot
	paths := make([]string, 0, len(currFiles))
	for p := range currFiles {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	snapshots := make([]models.VersionFileSnapshot, 0, len(paths))
	for _, p := range paths {
		snapshots = append(snapshots, models.VersionFileSnapshot{
			Version: req.Version,
			Path:    p,
			Hash:    currFiles[p],
		})
	}

	if err := s.snapshotRepo.BulkInsertTx(ctx, tx, snapshots); err != nil {
		return nil, fmt.Errorf("insert snapshots: %w", err)
	}

	if len(changes) > 0 {
		if err := s.versionChange.CreateBatchTx(ctx, tx, req.Version, changes); err != nil {
			return nil, fmt.Errorf("insert changes: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	// Create tar archive with changed files (after commit so version is persisted)
	if err := s.createArchive(req.Version, changes); err != nil {
		return nil, fmt.Errorf("create archive: %w", err)
	}

	return &models.VersionDTO{
		Version: v.Version,
		Date:    v.Date,
		Size:    v.Size,
		Notes:   req.Notes,
	}, nil
}

func (s *VersionService) Delete(ctx context.Context, version string) error {
	deleted, err := s.repo.DeleteIfLatest(ctx, version)
	if err != nil {
		return err
	}
	if !deleted {
		exists, _ := s.repo.FindByVersion(ctx, version)
		if exists == nil {
			return ErrVersionNotFound
		}
		return ErrNotLatestVersion
	}
	// Clean up archive file
	os.Remove(filepath.Join(s.archivesDir, version+".tar"))
	return nil
}

func computeChanges(prevSnapshots []models.VersionFileSnapshot, currFiles map[string]string) []models.ChangeEntry {
	prevMap := make(map[string]string, len(prevSnapshots))
	for _, s := range prevSnapshots {
		prevMap[s.Path] = s.Hash
	}

	var changes []models.ChangeEntry

	// Detect added and modified files (sorted for deterministic output)
	paths := make([]string, 0, len(currFiles))
	for p := range currFiles {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		hash := currFiles[path]
		prevHash, existed := prevMap[path]
		switch {
		case !existed:
			changes = append(changes, models.ChangeEntry{Action: "add", Path: path})
		case prevHash != hash:
			changes = append(changes, models.ChangeEntry{Action: "modify", Path: path})
		}
	}

	// Detect deleted files
	for _, s := range prevSnapshots {
		if _, stillExists := currFiles[s.Path]; !stillExists {
			changes = append(changes, models.ChangeEntry{Action: "delete", Path: s.Path})
		}
	}

	return changes
}

func (s *VersionService) createArchive(version string, changes []models.ChangeEntry) error {
	// Collect add/modify paths (delete has nothing to archive)
	var filesToArchive []string
	for _, c := range changes {
		if c.Action == "add" || c.Action == "modify" {
			filesToArchive = append(filesToArchive, c.Path)
		}
	}

	if len(filesToArchive) == 0 {
		return nil
	}

	archivePath := filepath.Join(s.archivesDir, version+".tar")
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	for _, relPath := range filesToArchive {
		absPath := filepath.Join(s.fileSvc.basePath, relPath)

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(relPath)

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		file, err := os.Open(absPath)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, file); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	return nil
}
