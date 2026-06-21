package services

import (
	"context"
	"errors"
	"time"

	"latestpack/repository"
	"latestpack/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCreds      = errors.New("invalid credentials")
	ErrUsernameExists    = errors.New("username already exists")
	ErrPasswordTooShort  = errors.New("password must be at least 6 characters")
	ErrCurrentPassword   = errors.New("current password is incorrect")
)

type AuthService struct {
	repo      *repository.UserRepo
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(repo *repository.UserRepo, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCreds
	}

	token, err := utils.GenerateToken(username, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) UpdateUsername(ctx context.Context, currentUsername, newUsername string) (string, error) {
	existing, err := s.repo.FindByUsername(ctx, newUsername)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrUsernameExists
	}

	if err := s.repo.UpdateUsername(ctx, currentUsername, newUsername); err != nil {
		return "", err
	}
	return newUsername, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, username, currentPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrCurrentPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, username, string(hash))
}
