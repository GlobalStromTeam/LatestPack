package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type StatsRepo struct {
	db *sql.DB
}

func NewStatsRepo(db *sql.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) FindByDate(ctx context.Context, date string) (*models.DailyStats, error) {
	var s models.DailyStats
	err := r.db.QueryRowContext(ctx,
		"SELECT date, launches, updates, traffic FROM stats WHERE date = ?",
		date,
	).Scan(&s.Date, &s.Launches, &s.Updates, &s.Traffic)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StatsRepo) Upsert(ctx context.Context, s *models.DailyStats) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO stats (date, launches, updates, traffic) VALUES (?, ?, ?, ?) ON CONFLICT(date) DO UPDATE SET launches=excluded.launches, updates=excluded.updates, traffic=excluded.traffic",
		s.Date, s.Launches, s.Updates, s.Traffic,
	)
	return err
}
