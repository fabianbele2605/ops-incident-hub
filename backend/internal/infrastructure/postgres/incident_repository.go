package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/resilience"
)

// IncidentRepository implementación PostgreSQL del repositorio de incidentes
type IncidentRepository struct {
	db *sql.DB
	circuitBreaker  *resilience.CircuitBreaker
}

// NewIncidentRepository crea una nueva instancia
func NewIncidentRepository(db *sql.DB) *IncidentRepository {
	return &IncidentRepository{
		db: db,
		circuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
}
}

// Create crea un nuevo incidente
func (r *IncidentRepository) Create(ctx context.Context, incident *domain.Incident) error {
	return r.circuitBreaker.Execute(func() error {
		metadata, err := json.Marshal(incident.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		query := `
			INSERT INTO incidents (
				id, title, description, severity, status, 
				assigned_to, created_by, created_at, updated_at, 
				resolved_at, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = r.db.ExecContext(
			ctx, query,
			incident.ID,
			incident.Title,
			incident.Description,
			incident.Severity,
			incident.Status,
			incident.AssignedTo,
			incident.CreatedBy,
			incident.CreatedAt,
			incident.UpdatedAt,
			incident.ResolvedAt,
			metadata,
		)

		if err != nil {
			return fmt.Errorf("failed to create incident: %w", err)
		}

		return nil
	})
}

// GetByID obtiene un incidente por su ID
func (r *IncidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	var incident *domain.Incident
	err := r.circuitBreaker.Execute(func() error {
		query := `
			SELECT id, title, description, severity, status, 
			       assigned_to, created_by, created_at, updated_at, 
			       resolved_at, metadata
			FROM incidents
			WHERE id = $1
		`

		incident = &domain.Incident{}
		var metadata []byte

		err := r.db.QueryRowContext(ctx, query, id).Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Status,
			&incident.AssignedTo,
			&incident.CreatedBy,
			&incident.CreatedAt,
			&incident.UpdatedAt,
			&incident.ResolvedAt,
			&metadata,
		)

		if err == sql.ErrNoRows {
			return domain.ErrIncidentNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to get incident: %w", err)
		}

		if err := json.Unmarshal(metadata, &incident.Metadata); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return incident, nil
}

// List obtiene una lista paginada de incidentes
func (r *IncidentRepository) List(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
	query := `
		SELECT id, title, description, severity, status, 
		       assigned_to, created_by, created_at, updated_at, 
		       resolved_at, metadata
		FROM incidents
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filters.Status)
		argPos++
	}

	if filters.Severity != nil {
		query += fmt.Sprintf(" AND severity = $%d", argPos)
		args = append(args, *filters.Severity)
		argPos++
	}

	if filters.AssignedTo != nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", argPos)
		args = append(args, *filters.AssignedTo)
		argPos++
	}

	if filters.CreatedBy != nil {
		query += fmt.Sprintf(" AND created_by = $%d", argPos)
		args = append(args, *filters.CreatedBy)
		argPos++
	}

	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list incidents: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	incidents := []*domain.Incident{}
	for rows.Next() {
		incident := &domain.Incident{}
		var metadata []byte

		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Status,
			&incident.AssignedTo,
			&incident.CreatedBy,
			&incident.CreatedAt,
			&incident.UpdatedAt,
			&incident.ResolvedAt,
			&metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}

		if err := json.Unmarshal(metadata, &incident.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating incidents: %w", err)
	}

	return incidents, nil
}

// Update actualiza un incidente existente
func (r *IncidentRepository) Update(ctx context.Context, incident *domain.Incident) error {
	return r.circuitBreaker.Execute(func() error {
		metadata, err := json.Marshal(incident.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		query := `
			UPDATE incidents
			SET title = $2, description = $3, severity = $4, status = $5,
			    assigned_to = $6, updated_at = $7, resolved_at = $8, metadata = $9
			WHERE id = $1
		`

		result, err := r.db.ExecContext(
			ctx, query,
			incident.ID,
			incident.Title,
			incident.Description,
			incident.Severity,
			incident.Status,
			incident.AssignedTo,
			incident.UpdatedAt,
			incident.ResolvedAt,
			metadata,
		)

		if err != nil {
			return fmt.Errorf("failed to update incident: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			return domain.ErrIncidentNotFound
		}

		return nil
	})
}

// Delete elimina un incidente (soft delete)
func (r *IncidentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM incidents WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete incident: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrIncidentNotFound
	}

	return nil
}

// Count cuenta el total de incidentes con filtros
func (r *IncidentRepository) Count(ctx context.Context, filters domain.ListFilters) (int, error) {
	query := `SELECT COUNT(*) FROM incidents WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filters.Status)
		argPos++
	}

	if filters.Severity != nil {
		query += fmt.Sprintf(" AND severity = $%d", argPos)
		args = append(args, *filters.Severity)
		argPos++
	}

	if filters.AssignedTo != nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", argPos)
		args = append(args, *filters.AssignedTo)
		argPos++
	}

	if filters.CreatedBy != nil {
		query += fmt.Sprintf(" AND created_by = $%d", argPos)
		args = append(args, *filters.CreatedBy)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count incidents: %w", err)
	}

	return count, nil
}
