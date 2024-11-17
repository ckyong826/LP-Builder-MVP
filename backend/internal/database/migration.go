package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a single database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// List of all migrations
var migrations = []Migration{
	{
		Version:     1,
		Description: "Create users table",
		Up: `
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE
			);
			CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
			CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
		`,
		Down: `
			DROP TABLE IF EXISTS users;
		`,
	},
	{
		Version:     2,
		Description: "Create templates table",
		Up: `
			CREATE TABLE IF NOT EXISTS templates (
				id BIGSERIAL PRIMARY KEY,
				original_url TEXT NOT NULL,
				html_path TEXT,
				file_paths TEXT,
				status VARCHAR(20) NOT NULL DEFAULT 'pending',
				error_message TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE,
				CONSTRAINT templates_status_check 
					CHECK (status IN ('pending', 'complete', 'failed', 'in_progress'))
			);
			CREATE INDEX IF NOT EXISTS idx_templates_status ON templates(status);
			CREATE INDEX IF NOT EXISTS idx_templates_created_at ON templates(created_at);
			CREATE INDEX IF NOT EXISTS idx_templates_deleted_at ON templates(deleted_at);
		`,
		Down: `
			DROP TABLE IF EXISTS templates;
		`,
	},
}

// Migrator handles database migrations
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := m.db.ExecContext(ctx, query)
	return err
}

// getAppliedMigrations returns a map of already applied migrations
func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	rows, err := m.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

// MigrateUp applies all pending migrations
func (m *Migrator) MigrateUp(ctx context.Context) error {
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if !applied[migration.Version] {
			log.Printf("Applying migration %d: %s", migration.Version, migration.Description)

			tx, err := m.db.BeginTx(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to start transaction: %w", err)
			}

			// Apply migration
			if _, err := tx.ExecContext(ctx, migration.Up); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}

			// Record migration
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO schema_migrations (version, description) VALUES ($1, $2)",
				migration.Version, migration.Description); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
			}

			log.Printf("Successfully applied migration %d", migration.Version)
		}
	}

	return nil
}

// MigrateDown reverts all migrations
func (m *Migrator) MigrateDown(ctx context.Context) error {
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Reverse order for down migrations
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if applied[migration.Version] {
			log.Printf("Reverting migration %d: %s", migration.Version, migration.Description)

			tx, err := m.db.BeginTx(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to start transaction: %w", err)
			}

			// Revert migration
			if _, err := tx.ExecContext(ctx, migration.Down); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to revert migration %d: %w", migration.Version, err)
			}

			// Remove migration record
			if _, err := tx.ExecContext(ctx,
				"DELETE FROM schema_migrations WHERE version = $1",
				migration.Version); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit revert of migration %d: %w", migration.Version, err)
			}

			log.Printf("Successfully reverted migration %d", migration.Version)
		}
	}

	return nil
}