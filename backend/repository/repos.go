package repository

import "database/sql"

type Repositories struct {
	DB             *sql.DB
	User           *UserRepo
	Version        *VersionRepo
	Stats          *StatsRepo
	Channel        *ChannelRepo
	VersionChange  *VersionChangeRepo
	VersionSnapshot *VersionSnapshotRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		DB:              db,
		User:            NewUserRepo(db),
		Version:         NewVersionRepo(db),
		Stats:           NewStatsRepo(db),
		Channel:         NewChannelRepo(db),
		VersionChange:   NewVersionChangeRepo(db),
		VersionSnapshot: NewVersionSnapshotRepo(db),
	}
}
