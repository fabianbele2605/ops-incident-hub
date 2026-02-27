package domain

import (
	"context"

	"github.com/google/uuid"
)

// IncidentRepository define el contrato para acceso a datos de incidentes
type IncidentRepository interface {
	// Create crea un nuevo incidente
	Create(ctx context.Context, incident *Incident) error

	// GetByID obtiene un incidente por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*Incident, error)

	// List obtiene una lista paginada de incidentes
	List(ctx context.Context, filters ListFilters) ([]*Incident, error)

	// Update actualiza un incidente existente
	Update(ctx context.Context, incident *Incident) error

	// Delete elimina un incidente (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Count cuenta el total de incidentes con filtros
	Count(ctx context.Context, filters ListFilters) (int, error)
}

// UserRepository define el contrato para acceso a datos de usuarios
type UserRepository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *User) error

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail obtiene un usuario por su email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List obtiene una lista de usuarios
	List(ctx context.Context) ([]*User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *User) error

	// Delete elimina un usuario
	Delete(ctx context.Context, id uuid.UUID) error
}

// ListFilters define los filtros para listar incidentes
type ListFilters struct {
	Status     *Status
	Severity   *Severity
	AssignedTo *uuid.UUID
	CreatedBy  *uuid.UUID
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}
