# Fase 2 - Handlers HTTP (API REST)

## 1. ¿Qué son los handlers?

Son la **capa de entrada HTTP** que expone los casos de uso como endpoints REST.

**Responsabilidades:**
- Recibir requests HTTP
- Validar entrada
- Llamar casos de uso
- Formatear respuestas
- Manejar errores HTTP

---

## 2. Estructura de carpetas

```
backend/internal/api/
├── handler/
│   └── incident_handler.go
├── middleware/
│   └── cors.go
├── dto/
│   ├── incident_dto.go
│   └── response.go
└── router.go
```

---

## 3. Crear carpetas

```bash
mkdir -p backend/internal/api/handler
mkdir -p backend/internal/api/middleware
mkdir -p backend/internal/api/dto
```

---

## 4. Dependencias necesarias

```bash
cd backend
go get github.com/gorilla/mux
```

**Explicación:** Gorilla Mux es un router HTTP popular para Go.

---

## 5. Archivo 1: `backend/internal/api/dto/response.go`

```go
package dto

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	TotalPages int         `json:"total_pages"`
}
```

**Explicación:** DTOs (Data Transfer Objects) para respuestas HTTP.

---

## 6. Archivo 2: `backend/internal/api/dto/incident_dto.go`

```go
package dto

import (
	"time"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
)

// CreateIncidentRequest representa la petición para crear un incidente
type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// AssignIncidentRequest representa la petición para asignar un incidente
type AssignIncidentRequest struct {
	AssignedTo string `json:"assigned_to"`
}

// IncidentResponse representa la respuesta de un incidente
type IncidentResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// ToIncidentResponse convierte una entidad Incident a DTO
func ToIncidentResponse(incident *domain.Incident) *IncidentResponse {
	response := &IncidentResponse{
		ID:          incident.ID.String(),
		Title:       incident.Title,
		Description: incident.Description,
		Severity:    string(incident.Severity),
		Status:      string(incident.Status),
		CreatedBy:   incident.CreatedBy.String(),
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
		ResolvedAt:  incident.ResolvedAt,
	}

	if incident.AssignedTo != nil {
		assignedTo := incident.AssignedTo.String()
		response.AssignedTo = &assignedTo
	}

	return response
}

// ToIncidentResponseList convierte una lista de incidentes a DTOs
func ToIncidentResponseList(incidents []*domain.Incident) []*IncidentResponse {
	responses := make([]*IncidentResponse, len(incidents))
	for i, incident := range incidents {
		responses[i] = ToIncidentResponse(incident)
	}
	return responses
}
```

**Explicación:** DTOs para requests y responses de incidentes.

---

## 7. Archivo 3: `backend/internal/api/handler/incident_handler.go`

```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/dto"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// IncidentHandler maneja las peticiones HTTP de incidentes
type IncidentHandler struct {
	createUseCase *incident.CreateIncidentUseCase
	assignUseCase *incident.AssignIncidentUseCase
	listUseCase   *incident.ListIncidentsUseCase
}

// NewIncidentHandler crea una nueva instancia
func NewIncidentHandler(
	createUseCase *incident.CreateIncidentUseCase,
	assignUseCase *incident.AssignIncidentUseCase,
	listUseCase *incident.ListIncidentsUseCase,
) *IncidentHandler {
	return &IncidentHandler{
		createUseCase: createUseCase,
		assignUseCase: assignUseCase,
		listUseCase:   listUseCase,
	}
}

// Create maneja POST /api/v1/incidents
func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Obtener usuario del contexto (autenticación)
	createdBy := uuid.New() // Temporal

	input := incident.CreateIncidentInput{
		Title:       req.Title,
		Description: req.Description,
		Severity:    domain.Severity(req.Severity),
		CreatedBy:   createdBy,
	}

	createdIncident, err := h.createUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, dto.ToIncidentResponse(createdIncident))
}

// Assign maneja POST /api/v1/incidents/:id/assign
func (h *IncidentHandler) Assign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	incidentID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	var req dto.AssignIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	assignedTo, err := uuid.Parse(req.AssignedTo)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid assigned_to ID")
		return
	}

	// TODO: Obtener usuario del contexto
	assignedBy := uuid.New() // Temporal

	input := incident.AssignIncidentInput{
		IncidentID: incidentID,
		AssignedTo: assignedTo,
		AssignedBy: assignedBy,
	}

	if err := h.assignUseCase.Execute(r.Context(), input); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "Incident assigned successfully",
	})
}

// List maneja GET /api/v1/incidents
func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parsear query params para filtros
	input := incident.ListIncidentsInput{
		Limit:  20,
		Offset: 0,
	}

	output, err := h.listUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	totalPages := (output.Total + input.Limit - 1) / input.Limit

	response := dto.PaginatedResponse{
		Data:       dto.ToIncidentResponseList(output.Incidents),
		Total:      output.Total,
		Limit:      input.Limit,
		Offset:     input.Offset,
		TotalPages: totalPages,
	}

	respondJSON(w, http.StatusOK, response)
}

// respondJSON envía una respuesta JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError envía una respuesta de error
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, dto.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// handleUseCaseError maneja errores de casos de uso
func handleUseCaseError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidTitle, domain.ErrInvalidDescription, domain.ErrInvalidSeverity:
		respondError(w, http.StatusBadRequest, err.Error())
	case domain.ErrIncidentNotFound, domain.ErrUserNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case domain.ErrForbidden:
		respondError(w, http.StatusForbidden, err.Error())
	case domain.ErrIncidentClosed, domain.ErrIncidentNotAssigned, domain.ErrIncidentNotResolved:
		respondError(w, http.StatusConflict, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
```

**Explicación:** Handlers que exponen los casos de uso como endpoints HTTP.

---

## 8. Archivo 4: `backend/internal/api/router.go`

```go
package api

import (
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
	"github.com/gorilla/mux"
)

// SetupRouter configura las rutas de la API
func SetupRouter(incidentHandler *handler.IncidentHandler) *mux.Router {
	router := mux.NewRouter()

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Incidents routes
	api.HandleFunc("/incidents", incidentHandler.Create).Methods("POST")
	api.HandleFunc("/incidents", incidentHandler.List).Methods("GET")
	api.HandleFunc("/incidents/{id}/assign", incidentHandler.Assign).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
```

**Explicación:** Configuración de rutas HTTP.

---

## 9. Orden de creación

1. Instalar dependencia: `go get github.com/gorilla/mux`
2. Crear carpetas
3. Crear `dto/response.go`
4. Crear `dto/incident_dto.go`
5. Crear `handler/incident_handler.go`
6. Crear `router.go`

---

## 10. Validar compilación

```bash
cd /mnt/n/ubuntu/azureSenior/backend
go build ./internal/api/...
```

---

**Empieza instalando la dependencia y creando las carpetas.**
