package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		"SELECT username, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.Username, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)",
		user.Username, user.PasswordHash, user.CreatedAt,
	)
	return err
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (r *UserRepo) UpdateUsername(ctx context.Context, oldUsername, newUsername string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET username = ? WHERE username = ?",
		newUsername, oldUsername,
	)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, username, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET password_hash = ? WHERE username = ?",
		passwordHash, username,
	)
	return err
}
