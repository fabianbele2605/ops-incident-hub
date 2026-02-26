# Fase 2 - Casos de Uso (Use Cases)

## 1. ¿Qué son los casos de uso?

Son la **lógica de aplicación** que orquesta las entidades del dominio y los repositorios para cumplir con los requisitos del negocio.

**Características:**
- Independientes de la interfaz (HTTP, CLI, etc)
- Usan las entidades del dominio
- Usan las interfaces de repositorios
- Contienen validaciones de negocio

---

## 2. Estructura de carpetas

```
backend/internal/usecase/
├── incident/
│   ├── create_incident.go
│   ├── assign_incident.go
│   ├── resolve_incident.go
│   └── list_incidents.go
└── user/
    ├── create_user.go
    └── get_user.go
```

---

## 3. Crear carpetas

```bash
mkdir -p backend/internal/usecase/incident
mkdir -p backend/internal/usecase/user
```

---

## 4. Archivo 1: `backend/internal/usecase/incident/create_incident.go`

```go
package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

// CreateIncidentInput representa los datos de entrada para crear un incidente
type CreateIncidentInput struct {
	Title       string
	Description string
	Severity    domain.Severity
	CreatedBy   uuid.UUID
}

// CreateIncidentUseCase maneja la creación de incidentes
type CreateIncidentUseCase struct {
	incidentRepo domain.IncidentRepository
	userRepo     domain.UserRepository
}

// NewCreateIncidentUseCase crea una nueva instancia del caso de uso
func NewCreateIncidentUseCase(
	incidentRepo domain.IncidentRepository,
	userRepo domain.UserRepository,
) *CreateIncidentUseCase {
	return &CreateIncidentUseCase{
		incidentRepo: incidentRepo,
		userRepo:     userRepo,
	}
}

// Execute ejecuta el caso de uso de crear incidente
func (uc *CreateIncidentUseCase) Execute(ctx context.Context, input CreateIncidentInput) (*domain.Incident, error) {
	// Validar que el usuario existe
	user, err := uc.userRepo.GetByID(ctx, input.CreatedBy)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Validar que el usuario puede crear incidentes
	if !user.CanManageIncidents() {
		return nil, domain.ErrForbidden
	}

	// Crear la entidad de incidente
	incident, err := domain.NewIncident(
		input.Title,
		input.Description,
		input.Severity,
		input.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	// Guardar en el repositorio
	if err := uc.incidentRepo.Create(ctx, incident); err != nil {
		return nil, err
	}

	return incident, nil
}
```

**Explicación:**
- **Input struct:** Define los datos necesarios
- **Constructor:** Inyecta dependencias (repositorios)
- **Execute:** Lógica del caso de uso
  1. Valida que el usuario existe
  2. Valida permisos
  3. Crea la entidad
  4. Guarda en el repositorio

---

## 5. Archivo 2: `backend/internal/usecase/incident/assign_incident.go`

```go
package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

// AssignIncidentInput representa los datos de entrada
type AssignIncidentInput struct {
	IncidentID uuid.UUID
	AssignedTo uuid.UUID
	AssignedBy uuid.UUID
}

// AssignIncidentUseCase maneja la asignación de incidentes
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
	// Validar que el usuario que asigna tiene permisos
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

	// Validar que el usuario asignado existe
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
	return uc.incidentRepo.Update(ctx, incident)
}
```

---

## 6. Archivo 3: `backend/internal/usecase/incident/list_incidents.go`

```go
package incident

import (
	"context"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

// ListIncidentsInput representa los filtros de búsqueda
type ListIncidentsInput struct {
	Status     *domain.Status
	Severity   *domain.Severity
	AssignedTo *uuid.UUID
	Limit      int
	Offset     int
}

// ListIncidentsOutput representa el resultado
type ListIncidentsOutput struct {
	Incidents []*domain.Incident
	Total     int
}

// ListIncidentsUseCase maneja el listado de incidentes
type ListIncidentsUseCase struct {
	incidentRepo domain.IncidentRepository
}

// NewListIncidentsUseCase crea una nueva instancia
func NewListIncidentsUseCase(incidentRepo domain.IncidentRepository) *ListIncidentsUseCase {
	return &ListIncidentsUseCase{
		incidentRepo: incidentRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListIncidentsUseCase) Execute(ctx context.Context, input ListIncidentsInput) (*ListIncidentsOutput, error) {
	// Valores por defecto para paginación
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	// Crear filtros
	filters := domain.ListFilters{
		Status:     input.Status,
		Severity:   input.Severity,
		AssignedTo: input.AssignedTo,
		Limit:      input.Limit,
		Offset:     input.Offset,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}

	// Obtener incidentes
	incidents, err := uc.incidentRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Contar total
	total, err := uc.incidentRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &ListIncidentsOutput{
		Incidents: incidents,
		Total:     total,
	}, nil
}
```

---

## 7. Orden de creación

1. Crear carpetas
2. Crear `create_incident.go`
3. Crear `assign_incident.go`
4. Crear `list_incidents.go`

---

## 8. Validar compilación

```bash
cd /mnt/n/ubuntu/azureSenior/backend
go build ./internal/usecase/...
```

---

## 9. Características de los casos de uso

✅ **Independientes de infraestructura:** No saben de HTTP, DB, etc  
✅ **Inyección de dependencias:** Reciben repositorios en el constructor  
✅ **Validaciones de negocio:** Verifican permisos y reglas  
✅ **Context-aware:** Soportan cancelación  
✅ **Testeables:** Fácil crear mocks de repositorios  

---

**Crea las carpetas y los 3 archivos con el contenido de arriba.**
