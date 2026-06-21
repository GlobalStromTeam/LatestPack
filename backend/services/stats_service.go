package services

import (
	"context"
	"math"
	"time"

	"latestpack/models"
	"latestpack/repository"
)

type StatsService struct {
	repo *repository.StatsRepo
}

func NewStatsService(repo *repository.StatsRepo) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	todayStats, err := s.repo.FindByDate(ctx, today)
	if err != nil {
		return nil, err
	}
	yesterdayStats, err := s.repo.FindByDate(ctx, yesterday)
	if err != nil {
		return nil, err
	}

	resp := &models.StatsResponse{}

	if todayStats != nil {
		resp.Launches.Value = todayStats.Launches
		resp.Updates.Value = todayStats.Updates
		resp.Traffic.Value = todayStats.Traffic
		resp.Traffic.Unit = "MB"
	}

	if yesterdayStats != nil && todayStats != nil {
		resp.Launches.Change = calcChange(todayStats.Launches, yesterdayStats.Launches)
		resp.Updates.Change = calcChange(todayStats.Updates, yesterdayStats.Updates)
		resp.Traffic.Change = calcChange(todayStats.Traffic, yesterdayStats.Traffic)
	}

	return resp, nil
}

func calcChange(today, yesterday int64) float64 {
	if yesterday == 0 {
		return 0
	}
	return math.Round(float64(today-yesterday)/float64(yesterday)*1000) / 10
}
