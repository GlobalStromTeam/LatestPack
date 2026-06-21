package services

import (
	"context"
	"errors"
	"time"

	"latestpack/repository"
	"latestpack/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCreds = errors.New("invalid credentials")

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
