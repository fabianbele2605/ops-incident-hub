package domain_test

import (
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIncident(t *testing.T) {
	createdBy := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		incident, err := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		require.NoError(t, err)
		assert.NotNil(t, incident)
		assert.NotEqual(t, uuid.Nil, incident.ID)
		assert.Equal(t, "Test Incident", incident.Title)
		assert.Equal(t, "Test Description", incident.Description)
		assert.Equal(t, domain.SeverityHigh, incident.Severity)
		assert.Equal(t, domain.StatusOpen, incident.Status)
		assert.Equal(t, createdBy, incident.CreatedBy)
		assert.Nil(t, incident.AssignedTo)
		assert.NotZero(t, incident.CreatedAt)
		assert.NotZero(t, incident.UpdatedAt)
		assert.Nil(t, incident.ResolvedAt)
		assert.NotNil(t, incident.Metadata)
	})

	t.Run("empty title", func(t *testing.T) {
		incident, err := domain.NewIncident(
			"",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		assert.Error(t, err)
		assert.Nil(t, incident)
		assert.Equal(t, domain.ErrInvalidTitle, err)
	})

	t.Run("empty description", func(t *testing.T) {
		incident, err := domain.NewIncident(
			"Test Incident",
			"",
			domain.SeverityHigh,
			createdBy,
		)

		assert.Error(t, err)
		assert.Nil(t, incident)
		assert.Equal(t, domain.ErrInvalidDescription, err)
	})

	t.Run("invalid severity", func(t *testing.T) {
		incident, err := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.Severity("invalid"),
			createdBy,
		)

		assert.Error(t, err)
		assert.Nil(t, incident)
		assert.Equal(t, domain.ErrInvalidSeverity, err)
	})
}

func TestIncident_Assign(t *testing.T) {
	createdBy := uuid.New()
	assignedTo := uuid.New()

	t.Run("successful assignment", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		err := incident.Assign(assignedTo)

		require.NoError(t, err)
		assert.NotNil(t, incident.AssignedTo)
		assert.Equal(t, assignedTo, *incident.AssignedTo)
		assert.Equal(t, domain.StatusAssigned, incident.Status)
	})

	t.Run("cannot assign closed incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()
		_ = incident.Close()

		err := incident.Assign(uuid.New())

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentClosed, err)
	})
}

func TestIncident_StartProgress(t *testing.T) {
	createdBy := uuid.New()
	assignedTo := uuid.New()

	t.Run("successful start progress", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)

		err := incident.StartProgress()

		require.NoError(t, err)
		assert.Equal(t, domain.StatusInProgress, incident.Status)
	})

	t.Run("cannot start progress if not assigned", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		err := incident.StartProgress()

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentNotAssigned, err)
	})

	t.Run("cannot start progress if closed", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()
		_ = incident.Close()

		err := incident.StartProgress()

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentClosed, err)
	})
}

func TestIncident_Resolve(t *testing.T) {
	createdBy := uuid.New()
	assignedTo := uuid.New()

	t.Run("successful resolve", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()

		err := incident.Resolve()

		require.NoError(t, err)
		assert.Equal(t, domain.StatusResolved, incident.Status)
		assert.NotNil(t, incident.ResolvedAt)
	})

	t.Run("cannot resolve closed incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()
		_ = incident.Close()

		err := incident.Resolve()

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentClosed, err)
	})
}

func TestIncident_Close(t *testing.T) {
	createdBy := uuid.New()
	assignedTo := uuid.New()

	t.Run("successful close", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()

		err := incident.Close()

		require.NoError(t, err)
		assert.Equal(t, domain.StatusClosed, incident.Status)
	})

	t.Run("cannot close if not resolved", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		err := incident.Close()

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIncidentNotResolved, err)
	})
}

func TestIncident_IsOpen(t *testing.T) {
	createdBy := uuid.New()
	assignedTo := uuid.New()

	t.Run("open incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)

		assert.True(t, incident.IsOpen())
	})

	t.Run("assigned incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)

		assert.True(t, incident.IsOpen())
	})

	t.Run("in progress incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()

		assert.True(t, incident.IsOpen())
	})

	t.Run("resolved incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()

		assert.False(t, incident.IsOpen())
	})

	t.Run("closed incident", func(t *testing.T) {
		incident, _ := domain.NewIncident(
			"Test Incident",
			"Test Description",
			domain.SeverityHigh,
			createdBy,
		)
		_ = incident.Assign(assignedTo)
		_ = incident.StartProgress()
		_ = incident.Resolve()
		_ = incident.Close()

		assert.False(t, incident.IsOpen())
	})
}
