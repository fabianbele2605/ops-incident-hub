package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations ejecuta las migraciones de base de datos
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Crear tabla de migraciones si no existe
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Obtener migraciones aplicadas
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Obtener archivos de migración
	files, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Aplicar migraciones pendientes
	for _, file := range files {
		if !strings.HasSuffix(file, ".up.sql") {
			continue
		}

		name := strings.TrimSuffix(filepath.Base(file), ".up.sql")
		if applied[name] {
			continue
		}

		if err := applyMigration(db, file, name); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", name, err)
		}

		fmt.Printf("Applied migration: %s\n", name)
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT name FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, rows.Err()
}

func getMigrationFiles(path string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(path, "*.sql"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func applyMigration(db *sql.DB, file, name string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback() // Ignore error as commit will be called if successful
	}()

	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (name) VALUES ($1)", name); err != nil {
		return err
	}

	return tx.Commit()
}
