package services

import (
	"context"
	"encoding/json"
	"errors"
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
	repo *repository.VersionRepo
}

func NewVersionService(repo *repository.VersionRepo) *VersionService {
	return &VersionService{repo: repo}
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

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
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
	return nil
}
