package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/resilience"
	"github.com/google/uuid"
)

// UserRepository implementación PostgreSQL del repositorio de usuarios
type UserRepository struct {
	db             *sql.DB
	circuitBreaker *resilience.CircuitBreaker
}

// NewUserRepository crea una nueva instancia
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:             db,
		circuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
	}
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, name, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		user.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User
	err := r.circuitBreaker.Execute(func() error {
		query := `
			SELECT id, email, name, role, created_at
			FROM users
			WHERE id = $1
		`

		user = &domain.User{}
		err := r.db.QueryRowContext(ctx, query, id).Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
		)

		if err == sql.ErrNoRows {
			return domain.ErrUserNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, name, role, created_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// List obtiene una lista de usuarios
func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, name, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, role = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete elimina un usuario
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
