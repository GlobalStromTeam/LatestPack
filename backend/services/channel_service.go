package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"latestpack/models"
	"latestpack/repository"
)

var (
	ErrChannelNotFound    = errors.New("channel not found")
	ErrChannelNameExists  = errors.New("channel name already exists")
	ErrCannotDeleteLocal  = errors.New("cannot delete local channel")
	ErrCannotCreateLocal  = errors.New("cannot create local channel")
	ErrLocalChannelLocked = errors.New("local channel only allows enabled changes")
)

type ChannelService struct {
	repo *repository.ChannelRepo
}

func NewChannelService(repo *repository.ChannelRepo) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) List(ctx context.Context) ([]models.ChannelDTO, error) {
	channels, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]models.ChannelDTO, 0, len(channels))
	for _, ch := range channels {
		dto, err := toChannelDTO(ch)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (s *ChannelService) Create(ctx context.Context, req *models.CreateChannelRequest) (*models.ChannelDTO, error) {
	if req.Type == "local" {
		return nil, ErrCannotCreateLocal
	}
	if req.Type != "webdav" && req.Type != "s3" {
		return nil, fmt.Errorf("invalid channel type: %s", req.Type)
	}

	existing, err := s.repo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrChannelNameExists
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	weight := 50
	if req.Weight != nil {
		weight = *req.Weight
	}

	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := fmt.Sprintf("ch_%s_%d", req.Type, now.UnixMilli())

	ch := &models.Channel{
		ID:        id,
		Name:      req.Name,
		Type:      req.Type,
		Enabled:   enabled,
		Weight:    weight,
		Config:    string(configJSON),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, ch); err != nil {
		return nil, err
	}

	dto, err := toChannelDTO(*ch)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (s *ChannelService) Update(ctx context.Context, id string, req *models.UpdateChannelRequest) (*models.ChannelDTO, error) {
	ch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ch == nil {
		return nil, ErrChannelNotFound
	}

	if ch.Type == "local" {
		if req.Name != nil || req.Weight != nil || req.Config != nil {
			return nil, ErrLocalChannelLocked
		}
		if req.Enabled != nil {
			ch.Enabled = *req.Enabled
		}
		ch.UpdatedAt = time.Now()
		if err := s.repo.UpdateEnabled(ctx, ch.ID, ch.Enabled, ch.UpdatedAt); err != nil {
			return nil, err
		}
	} else {
		if req.Name != nil {
			existing, err := s.repo.FindByName(ctx, *req.Name)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != id {
				return nil, ErrChannelNameExists
			}
			ch.Name = *req.Name
		}
		if req.Enabled != nil {
			ch.Enabled = *req.Enabled
		}
		if req.Weight != nil {
			ch.Weight = *req.Weight
		}
		if req.Config != nil {
			configJSON, err := json.Marshal(req.Config)
			if err != nil {
				return nil, err
			}
			ch.Config = string(configJSON)
		}
		ch.UpdatedAt = time.Now()
		if err := s.repo.Update(ctx, ch); err != nil {
			return nil, err
		}
	}

	dto, err := toChannelDTO(*ch)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (s *ChannelService) Delete(ctx context.Context, id string) error {
	ch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ch == nil {
		return ErrChannelNotFound
	}
	if ch.Type == "local" {
		return ErrCannotDeleteLocal
	}
	return s.repo.Delete(ctx, id)
}

func toChannelDTO(ch models.Channel) (models.ChannelDTO, error) {
	var config models.ChannelConfig
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		config = models.ChannelConfig{}
	}
	return models.ChannelDTO{
		ID:      ch.ID,
		Name:    ch.Name,
		Type:    ch.Type,
		Enabled: ch.Enabled,
		Weight:  ch.Weight,
		Config:  config,
	}, nil
}
