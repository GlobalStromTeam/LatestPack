package seed

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"latestpack/models"
	"latestpack/repository"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser(repo *repository.UserRepo) error {
	ctx := context.Background()
	count, err := repo.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Println("Seeding admin user (admin / admin123)")
	return repo.Create(ctx, &models.User{
		Username:     "admin",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	})
}

func SeedLocalChannel(repo *repository.ChannelRepo) error {
	ctx := context.Background()
	existing, err := repo.FindByID(ctx, "ch_local")
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	configJSON, _ := json.Marshal(models.ChannelConfig{})
	now := time.Now()
	log.Println("Seeding local channel")
	return repo.Create(ctx, &models.Channel{
		ID:        "ch_local",
		Name:      "本地存储",
		Type:      "local",
		Enabled:   true,
		Weight:    0,
		Config:    string(configJSON),
		CreatedAt: now,
		UpdatedAt: now,
	})
}
