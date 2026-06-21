package seed

import (
	"context"
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
