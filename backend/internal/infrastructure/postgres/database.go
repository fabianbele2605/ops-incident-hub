package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewDatabase crea una nueva conexión a PostgreSQL
func NewDatabase(config DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configurar connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)       // Máximo de conexiones abiertas
	db.SetMaxIdleConns(config.MaxIdleConns)       // Máximo de conexiones idle
	db.SetConnMaxLifetime(config.ConnMaxLifetime) // Tiempo de vida máximo
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // Tiempo idle máximo

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
