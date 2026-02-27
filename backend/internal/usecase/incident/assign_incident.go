package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/observability/metrics"
	"github.com/google/uuid"
)

// AssignIncidentInput representa los datos de entrada
type AssignIncidentInput struct {
	IncidentID uuid.UUID
	AssignedTo uuid.UUID
	AssignedBy uuid.UUID
}

// AssignIncidentUseCase maneja la asignacion de incidentes
type AssignIncidentUseCase struct {
	incidentRepo domain.IncidentRepository
	userRepo     domain.UserRepository
}

// NewAssignIncidentUseCase crea una nueva instancia
func NewAssignIncidentUseCase(
	incidentRepo domain.IncidentRepository,
	userRepo domain.UserRepository,
) *AssignIncidentUseCase {
	return &AssignIncidentUseCase{
		incidentRepo: incidentRepo,
		userRepo:     userRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *AssignIncidentUseCase) Execute(ctx context.Context, input AssignIncidentInput) error {
	// Validar que el usuario que asigna tiene permiso
	assignedByUser, err := uc.userRepo.GetByID(ctx, input.AssignedBy)
	if err != nil {
		return err
	}
	if assignedByUser == nil {
		return domain.ErrUserNotFound
	}
	if !assignedByUser.CanManageIncidents() {
		return domain.ErrForbidden
	}

	//Validar que el usuario asignado existe
	assignedToUser, err := uc.userRepo.GetByID(ctx, input.AssignedTo)
	if err != nil {
		return err
	}
	if assignedToUser == nil {
		return domain.ErrUserNotFound
	}

	// Obtener el incidente
	incident, err := uc.incidentRepo.GetByID(ctx, input.IncidentID)
	if err != nil {
		return err
	}
	if incident == nil {
		return domain.ErrIncidentNotFound
	}

	// Asignar el incidente
	if err := incident.Assign(input.AssignedTo); err != nil {
		return err
	}

	// Actualizar en el repositorio
	if err := uc.incidentRepo.Update(ctx, incident); err != nil {
		return err
	}

	// Registrar métrica
	metrics.RecordIncidentAssigned()

	return nil
}
