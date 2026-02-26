package incident_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListIncidentsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful list without filters", func(t *testing.T) {
		adminUser, _ := domain.NewUser("admin@test.com", "Admin", domain.RoleAdmin)
		inc1, _ := domain.NewIncident("Inc 1", "Desc 1", domain.SeverityHigh, adminUser.ID)
		inc2, _ := domain.NewIncident("Inc 2", "Desc 2", domain.SeverityMedium, adminUser.ID)
		
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				return []*domain.Incident{inc1, inc2}, nil
			},
			countFunc: func(ctx context.Context, filters domain.ListFilters) (int, error) {
				return 2, nil
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{
			Limit:  20,
			Offset: 0,
		}
		
		result, err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Incidents, 2)
		assert.Equal(t, 2, result.Total)
	})

	t.Run("default pagination", func(t *testing.T) {
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				assert.Equal(t, 20, filters.Limit)
				return []*domain.Incident{}, nil
			},
			countFunc: func(ctx context.Context, filters domain.ListFilters) (int, error) {
				return 0, nil
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{}
		
		_, err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
	})

	t.Run("max limit 100", func(t *testing.T) {
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				assert.Equal(t, 100, filters.Limit)
				return []*domain.Incident{}, nil
			},
			countFunc: func(ctx context.Context, filters domain.ListFilters) (int, error) {
				return 0, nil
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{
			Limit: 200,
		}
		
		_, err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
	})

	t.Run("list with filters", func(t *testing.T) {
		status := domain.StatusOpen
		severity := domain.SeverityHigh
		assignedTo := uuid.New()
		
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				assert.Equal(t, &status, filters.Status)
				assert.Equal(t, &severity, filters.Severity)
				assert.Equal(t, &assignedTo, filters.AssignedTo)
				return []*domain.Incident{}, nil
			},
			countFunc: func(ctx context.Context, filters domain.ListFilters) (int, error) {
				return 0, nil
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{
			Status:     &status,
			Severity:   &severity,
			AssignedTo: &assignedTo,
			Limit:      10,
			Offset:     0,
		}
		
		result, err := uc.Execute(ctx, input)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("repository list error", func(t *testing.T) {
		repoErr := errors.New("database error")
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				return nil, repoErr
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})

	t.Run("repository count error", func(t *testing.T) {
		repoErr := errors.New("database error")
		incidentRepo := &mockIncidentRepository{
			listFunc: func(ctx context.Context, filters domain.ListFilters) ([]*domain.Incident, error) {
				return []*domain.Incident{}, nil
			},
			countFunc: func(ctx context.Context, filters domain.ListFilters) (int, error) {
				return 0, repoErr
			},
		}
		
		uc := incident.NewListIncidentUseCase(incidentRepo)
		
		input := incident.ListIncidentsInput{}
		
		result, err := uc.Execute(ctx, input)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})
}
