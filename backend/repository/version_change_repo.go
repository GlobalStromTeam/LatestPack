package repository

import (
	"context"
	"database/sql"

	"latestpack/models"
)

type VersionChangeRepo struct {
	db *sql.DB
}

func NewVersionChangeRepo(db *sql.DB) *VersionChangeRepo {
	return &VersionChangeRepo{db: db}
}

func (r *VersionChangeRepo) Create(ctx context.Context, vc *models.VersionChange) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO version_changes (version, action, path) VALUES (?, ?, ?)",
		vc.Version, vc.Action, vc.Path,
	)
	return err
}

func (r *VersionChangeRepo) CreateBatch(ctx context.Context, version string, changes []models.ChangeEntry) error {
	if len(changes) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO version_changes (version, action, path) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, c := range changes {
		if _, err := stmt.ExecContext(ctx, version, c.Action, c.Path); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *VersionChangeRepo) CreateBatchTx(ctx context.Context, tx *sql.Tx, version string, changes []models.ChangeEntry) error {
	if len(changes) == 0 {
		return nil
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO version_changes (version, action, path) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range changes {
		if _, err := stmt.ExecContext(ctx, version, c.Action, c.Path); err != nil {
			return err
		}
	}
	return nil
}

func (r *VersionChangeRepo) ListByVersion(ctx context.Context, version string) ([]models.ChangeEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT action, path FROM version_changes WHERE version = ? ORDER BY id",
		version,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []models.ChangeEntry
	for rows.Next() {
		var c models.ChangeEntry
		if err := rows.Scan(&c.Action, &c.Path); err != nil {
			return nil, err
		}
		changes = append(changes, c)
	}
	return changes, nil
}

// ListAfterVersion returns all versions with their changes that were created after the given version.
// It uses created_at ordering to determine version sequence.
func (r *VersionChangeRepo) ListAfterVersion(ctx context.Context, afterVersion string) ([]models.VersionChange, error) {
	// Use COALESCE so that if the from version is not found, we return all changes
	query := `
		SELECT vc.id, vc.version, vc.action, vc.path
		FROM version_changes vc
		JOIN versions v ON vc.version = v.version
		WHERE v.created_at > COALESCE((SELECT created_at FROM versions WHERE version = ?), '1970-01-01')
		ORDER BY v.created_at ASC, vc.id ASC`

	rows, err := r.db.QueryContext(ctx, query, afterVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []models.VersionChange
	for rows.Next() {
		var c models.VersionChange
		if err := rows.Scan(&c.ID, &c.Version, &c.Action, &c.Path); err != nil {
			return nil, err
		}
		changes = append(changes, c)
	}
	return changes, nil
}

// ListAll returns all version changes ordered by version creation time.
func (r *VersionChangeRepo) ListAll(ctx context.Context) ([]models.VersionChange, error) {
	query := `
		SELECT vc.id, vc.version, vc.action, vc.path
		FROM version_changes vc
		JOIN versions v ON vc.version = v.version
		ORDER BY v.created_at ASC, vc.id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []models.VersionChange
	for rows.Next() {
		var c models.VersionChange
		if err := rows.Scan(&c.ID, &c.Version, &c.Action, &c.Path); err != nil {
			return nil, err
		}
		changes = append(changes, c)
	}
	return changes, nil
}

func (r *VersionChangeRepo) DeleteByVersion(ctx context.Context, version string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM version_changes WHERE version = ?", version)
	return err
}
