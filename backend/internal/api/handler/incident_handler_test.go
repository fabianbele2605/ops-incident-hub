package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentHandler_Create(t *testing.T) {
	t.Run("success - creates incident with valid data", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.CreateIncidentInput) (*domain.Incident, error) {
				inc, _ := domain.NewIncident(input.Title, input.Description, input.Severity, input.CreatedBy)
				return inc, nil
			},
		}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		reqBody := map[string]interface{}{
			"title":       "Test Incident",
			"description": "Test Description",
			"severity":    "high",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Act
		h.Create(rec, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Test Incident", response["title"])
		assert.Equal(t, "open", response["status"])
	})

	t.Run("error - invalid request body", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader([]byte("invalid json")))
		rec := httptest.NewRecorder()

		// Act
		h.Create(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error - validation error from use case", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.CreateIncidentInput) (*domain.Incident, error) {
				return nil, domain.ErrInvalidTitle
			},
		}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		reqBody := map[string]interface{}{
			"title":       "",
			"description": "Test",
			"severity":    "high",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Act
		h.Create(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestIncidentHandler_Assign(t *testing.T) {
	t.Run("success - assigns incident to user", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.AssignIncidentInput) error {
				return nil
			},
		}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		incidentID := uuid.New()
		assignedTo := uuid.New()

		reqBody := map[string]interface{}{
			"assigned_to": assignedTo.String(),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents/"+incidentID.String()+"/assign", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// Configurar mux vars
		req = mux.SetURLVars(req, map[string]string{"id": incidentID.String()})

		// Act
		h.Assign(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error - invalid incident ID", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		reqBody := map[string]interface{}{
			"assigned_to": uuid.New().String(),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents/invalid-id/assign", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})

		// Act
		h.Assign(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error - incident not found", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.AssignIncidentInput) error {
				return domain.ErrIncidentNotFound
			},
		}
		mockList := &MockListIncidentsUseCase{}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		incidentID := uuid.New()
		assignedTo := uuid.New()

		reqBody := map[string]interface{}{
			"assigned_to": assignedTo.String(),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/incidents/"+incidentID.String()+"/assign", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": incidentID.String()})

		// Act
		h.Assign(rec, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestIncidentHandler_List(t *testing.T) {
	t.Run("success - returns paginated incidents", func(t *testing.T) {
		// Arrange
		inc1, _ := domain.NewIncident("Incident 1", "Description 1", domain.SeverityHigh, uuid.New())
		inc2, _ := domain.NewIncident("Incident 2", "Description 2", domain.SeverityMedium, uuid.New())

		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.ListIncidentsInput) (*incident.ListIncidentsOutput, error) {
				return &incident.ListIncidentsOutput{
					Incidents: []*domain.Incident{inc1, inc2},
					Total:     2,
				}, nil
			},
		}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/incidents", nil)
		rec := httptest.NewRecorder()

		// Act
		h.List(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, float64(2), response["total"])
		assert.NotNil(t, response["data"])
	})

	t.Run("error - use case returns error", func(t *testing.T) {
		// Arrange
		mockCreate := &MockCreateIncidentUseCase{}
		mockAssign := &MockAssignIncidentUseCase{}
		mockList := &MockListIncidentsUseCase{
			ExecuteFunc: func(ctx context.Context, input incident.ListIncidentsInput) (*incident.ListIncidentsOutput, error) {
				return nil, domain.ErrIncidentNotFound
			},
		}

		h := handler.NewIncidentHandler(mockCreate, mockAssign, mockList)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/incidents", nil)
		rec := httptest.NewRecorder()

		// Act
		h.List(rec, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
