package repository

import "database/sql"

type Repositories struct {
	User    *UserRepo
	Version *VersionRepo
	Stats   *StatsRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepo(db),
		Version: NewVersionRepo(db),
		Stats:   NewStatsRepo(db),
	}
}
