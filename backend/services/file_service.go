package services

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"latestpack/models"
	"latestpack/utils"
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrItemExists    = errors.New("item already exists")
	ErrInvalidPath   = errors.New("invalid path")
	ErrInvalidName   = errors.New("invalid name")
	ErrFileTooLarge  = errors.New("file too large")
)

const maxFileSize int64 = 500 << 20 // 500 MB

type FileService struct {
	basePath string
}

func NewFileService(basePath string) *FileService {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Printf("Warning: create files dir: %v", err)
	}
	return &FileService{basePath: basePath}
}

// sanitizeName extracts the base component and rejects names with path separators or empty strings.
func sanitizeName(name string) (string, error) {
	name = filepath.Base(name)
	if name == "." || name == ".." || name == "" {
		return "", ErrInvalidName
	}
	if strings.ContainsAny(name, "/\\") {
		return "", ErrInvalidName
	}
	return name, nil
}

func (s *FileService) resolvePath(virtualPath string) (string, error) {
	clean := filepath.Clean(filepath.Join(s.basePath, virtualPath))
	if !strings.HasPrefix(clean, filepath.Clean(s.basePath)+string(os.PathSeparator)) && clean != filepath.Clean(s.basePath) {
		return "", ErrInvalidPath
	}
	return clean, nil
}

func (s *FileService) List(virtualPath string) (*models.FileListResponse, error) {
	dir, err := s.resolvePath(virtualPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return &models.FileListResponse{Path: virtualPath, Items: []models.FileItem{}}, nil
		}
		return nil, err
	}

	items := make([]models.FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := models.FileItem{
			Name: entry.Name(),
			Date: info.ModTime().Format("2006-01-02"),
		}

		if entry.IsDir() {
			item.Type = "folder"
			item.Size = "-"
		} else {
			item.Type = "file"
			item.Size = utils.FormatSize(info.Size())
		}

		items = append(items, item)
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type == "folder"
		}
		return items[i].Name < items[j].Name
	})

	return &models.FileListResponse{Path: virtualPath, Items: items}, nil
}

func (s *FileService) CreateFolder(virtualPath, name string) (*models.FileItem, error) {
	parent, err := s.resolvePath(virtualPath)
	if err != nil {
		return nil, err
	}

	cleanName, err := sanitizeName(name)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(parent, cleanName)

	if _, err := os.Stat(fullPath); err == nil {
		return nil, ErrItemExists
	}

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return nil, err
	}

	return &models.FileItem{
		Name: cleanName,
		Type: "folder",
		Size: "-",
		Date: time.Now().Format("2006-01-02"),
	}, nil
}

func (s *FileService) SaveFile(virtualPath, filename string, data []byte) (*models.FileItem, error) {
	if int64(len(data)) > maxFileSize {
		return nil, ErrFileTooLarge
	}

	parent, err := s.resolvePath(virtualPath)
	if err != nil {
		return nil, err
	}

	cleanName, err := sanitizeName(filename)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(parent, cleanName)

	if err := os.MkdirAll(parent, 0755); err != nil {
		return nil, err
	}

	// Use O_EXCL to create — fails if file exists, prevents TOCTOU race
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil, ErrItemExists
		}
		return nil, err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		os.Remove(fullPath)
		return nil, err
	}

	return &models.FileItem{
		Name: cleanName,
		Type: "file",
		Size: utils.FormatSize(int64(len(data))),
		Date: time.Now().Format("2006-01-02"),
	}, nil
}

func (s *FileService) Rename(virtualPath, oldName, newName string) (*models.FileItem, error) {
	parent, err := s.resolvePath(virtualPath)
	if err != nil {
		return nil, err
	}

	cleanOld, err := sanitizeName(oldName)
	if err != nil {
		return nil, err
	}
	cleanNew, err := sanitizeName(newName)
	if err != nil {
		return nil, err
	}

	oldPath := filepath.Join(parent, cleanOld)
	newPath := filepath.Join(parent, cleanNew)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil, ErrItemNotFound
	}
	if _, err := os.Stat(newPath); err == nil {
		return nil, ErrItemExists
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return nil, err
	}

	info, statErr := os.Stat(newPath)
	item := &models.FileItem{
		Name: cleanNew,
		Date: time.Now().Format("2006-01-02"),
	}
	if statErr != nil {
		item.Type = "file"
		item.Size = "-"
	} else if info.IsDir() {
		item.Type = "folder"
		item.Size = "-"
	} else {
		item.Type = "file"
		item.Size = utils.FormatSize(info.Size())
	}

	return item, nil
}

func (s *FileService) Delete(virtualPath, name string) error {
	parent, err := s.resolvePath(virtualPath)
	if err != nil {
		return err
	}

	cleanName, err := sanitizeName(name)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(parent, cleanName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return ErrItemNotFound
	}

	log.Printf("File delete: %s (user-requested)", fullPath)
	return os.RemoveAll(fullPath)
}
