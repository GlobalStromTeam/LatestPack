package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type VersionRepo struct {
	db *sql.DB
}

func NewVersionRepo(db *sql.DB) *VersionRepo {
	return &VersionRepo{db: db}
}

func (r *VersionRepo) List(ctx context.Context, page, pageSize int64) ([]models.Version, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM versions").Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx,
		"SELECT version, date, size, notes, created_at FROM versions ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var versions []models.Version
	for rows.Next() {
		var v models.Version
		if err := rows.Scan(&v.Version, &v.Date, &v.Size, &v.Notes, &v.CreatedAt); err != nil {
			return nil, 0, err
		}
		versions = append(versions, v)
	}
	return versions, total, nil
}

func (r *VersionRepo) FindByVersion(ctx context.Context, version string) (*models.Version, error) {
	var v models.Version
	err := r.db.QueryRowContext(ctx,
		"SELECT version, date, size, notes, created_at FROM versions WHERE version = ?",
		version,
	).Scan(&v.Version, &v.Date, &v.Size, &v.Notes, &v.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VersionRepo) GetLatest(ctx context.Context) (*models.Version, error) {
	var v models.Version
	err := r.db.QueryRowContext(ctx,
		"SELECT version, date, size, notes, created_at FROM versions ORDER BY created_at DESC LIMIT 1",
	).Scan(&v.Version, &v.Date, &v.Size, &v.Notes, &v.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VersionRepo) Create(ctx context.Context, v *models.Version) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO versions (version, date, size, notes, created_at) VALUES (?, ?, ?, ?, ?)",
		v.Version, v.Date, v.Size, v.Notes, v.CreatedAt,
	)
	return err
}

func (r *VersionRepo) DeleteIfLatest(ctx context.Context, version string) (bool, error) {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM versions WHERE version = ? AND version = (SELECT version FROM versions ORDER BY created_at DESC LIMIT 1)",
		version,
	)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

func (r *VersionRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM versions").Scan(&count)
	return count, err
}

