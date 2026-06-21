package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type VersionSnapshotRepo struct {
	db *sql.DB
}

func NewVersionSnapshotRepo(db *sql.DB) *VersionSnapshotRepo {
	return &VersionSnapshotRepo{db: db}
}

func (r *VersionSnapshotRepo) BulkInsertTx(ctx context.Context, tx *sql.Tx, snapshots []models.VersionFileSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO version_file_snapshots (version, path, hash) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range snapshots {
		if _, err := stmt.ExecContext(ctx, s.Version, s.Path, s.Hash); err != nil {
			return err
		}
	}
	return nil
}

// GetLatest returns the snapshot of the latest version (by created_at DESC).
// Must be called BEFORE inserting the new version row in the same transaction.
func (r *VersionSnapshotRepo) GetLatest(ctx context.Context, tx *sql.Tx) ([]models.VersionFileSnapshot, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT vs.id, vs.version, vs.path, vs.hash
		FROM version_file_snapshots vs
		JOIN versions v ON vs.version = v.version
		WHERE v.version = (SELECT v2.version FROM versions v2 ORDER BY v2.created_at DESC LIMIT 1)
		ORDER BY vs.path`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.VersionFileSnapshot
	for rows.Next() {
		var s models.VersionFileSnapshot
		if err := rows.Scan(&s.ID, &s.Version, &s.Path, &s.Hash); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (r *VersionSnapshotRepo) DeleteByVersion(ctx context.Context, version string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM version_file_snapshots WHERE version = ?", version)
	return err
}
