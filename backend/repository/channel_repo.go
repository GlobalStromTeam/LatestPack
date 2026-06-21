package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type ChannelRepo struct {
	db *sql.DB
}

func NewChannelRepo(db *sql.DB) *ChannelRepo {
	return &ChannelRepo{db: db}
}

func (r *ChannelRepo) GetAll(ctx context.Context) ([]models.Channel, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, type, enabled, weight, config, created_at, updated_at FROM channels ORDER BY weight ASC, created_at ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var ch models.Channel
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Type, &ch.Enabled, &ch.Weight, &ch.Config, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

func (r *ChannelRepo) FindByID(ctx context.Context, id string) (*models.Channel, error) {
	var ch models.Channel
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, type, enabled, weight, config, created_at, updated_at FROM channels WHERE id = ?",
		id,
	).Scan(&ch.ID, &ch.Name, &ch.Type, &ch.Enabled, &ch.Weight, &ch.Config, &ch.CreatedAt, &ch.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (r *ChannelRepo) FindByName(ctx context.Context, name string) (*models.Channel, error) {
	var ch models.Channel
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, type, enabled, weight, config, created_at, updated_at FROM channels WHERE name = ?",
		name,
	).Scan(&ch.ID, &ch.Name, &ch.Type, &ch.Enabled, &ch.Weight, &ch.Config, &ch.CreatedAt, &ch.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (r *ChannelRepo) Create(ctx context.Context, ch *models.Channel) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO channels (id, name, type, enabled, weight, config, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		ch.ID, ch.Name, ch.Type, ch.Enabled, ch.Weight, ch.Config, ch.CreatedAt, ch.UpdatedAt,
	)
	return err
}

func (r *ChannelRepo) Update(ctx context.Context, ch *models.Channel) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE channels SET name = ?, enabled = ?, weight = ?, config = ?, updated_at = ? WHERE id = ?",
		ch.Name, ch.Enabled, ch.Weight, ch.Config, ch.UpdatedAt, ch.ID,
	)
	return err
}

func (r *ChannelRepo) UpdateEnabled(ctx context.Context, id string, enabled bool, updatedAt interface{}) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE channels SET enabled = ?, updated_at = ? WHERE id = ?",
		enabled, updatedAt, id,
	)
	return err
}

func (r *ChannelRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM channels WHERE id = ?", id)
	return err
}
