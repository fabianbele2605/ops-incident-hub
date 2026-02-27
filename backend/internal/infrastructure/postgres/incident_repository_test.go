package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=ops_incident_hub_test sslmode=disable"
	
	db, err := postgres.NewDatabase(postgres.DatabaseConfig{
		DSN:             dsn,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 300,
	})
	require.NoError(t, err)
	
	// Ejecutar migraciones
	err = postgres.RunMigrations(db, "../../../migrations")
	require.NoError(t, err)
	
	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE incidents, users CASCADE")
	require.NoError(t, err)
	
	err = db.Close()
	require.NoError(t, err)
}

func TestIncidentRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	ctx := context.Background()
	
	// Crear usuario primero
	user, _ := domain.NewUser("creator@example.com", "Creator", domain.RoleAdmin)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	
	incident, err := domain.NewIncident("Test Incident", "Test Description", domain.SeverityHigh, user.ID)
	require.NoError(t, err)
	
	err = repo.Create(ctx, incident)
	require.NoError(t, err)
	
	// Verificar que se creó
	retrieved, err := repo.GetByID(ctx, incident.ID)
	require.NoError(t, err)
	assert.Equal(t, incident.Title, retrieved.Title)
	assert.Equal(t, incident.Description, retrieved.Description)
	assert.Equal(t, incident.Severity, retrieved.Severity)
}

func TestIncidentRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	ctx := context.Background()
	
	// Crear usuario
	user, _ := domain.NewUser("creator@example.com", "Creator", domain.RoleAdmin)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	
	// Crear incidente
	incident, _ := domain.NewIncident("Test", "Desc", domain.SeverityMedium, user.ID)
	err = repo.Create(ctx, incident)
	require.NoError(t, err)
	
	// Obtener por ID
	retrieved, err := repo.GetByID(ctx, incident.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, incident.ID, retrieved.ID)
}

func TestIncidentRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	ctx := context.Background()
	
	retrieved, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Equal(t, domain.ErrIncidentNotFound, err)
}

func TestIncidentRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	ctx := context.Background()
	
	// Crear usuario
	user, _ := domain.NewUser("creator@example.com", "Creator", domain.RoleAdmin)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	
	// Crear incidente
	incident, _ := domain.NewIncident("Original", "Desc", domain.SeverityLow, user.ID)
	err = repo.Create(ctx, incident)
	require.NoError(t, err)
	
	// Actualizar
	incident.Title = "Updated"
	err = repo.Update(ctx, incident)
	require.NoError(t, err)
	
	// Verificar actualización
	retrieved, err := repo.GetByID(ctx, incident.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", retrieved.Title)
}

func TestIncidentRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	ctx := context.Background()
	
	// Crear usuario
	user, _ := domain.NewUser("creator@example.com", "Creator", domain.RoleAdmin)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	
	// Crear varios incidentes
	for i := 0; i < 3; i++ {
		incident, _ := domain.NewIncident("Incident", "Desc", domain.SeverityHigh, user.ID)
		err := repo.Create(ctx, incident)
		require.NoError(t, err)
	}
	
	// Listar
	filters := domain.ListFilters{
		Limit:  10,
		Offset: 0,
	}
	incidents, err := repo.List(ctx, filters)
	require.NoError(t, err)
	assert.Len(t, incidents, 3)
}

func TestIncidentRepository_Count(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	ctx := context.Background()
	
	// Crear usuario
	user, _ := domain.NewUser("creator@example.com", "Creator", domain.RoleAdmin)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	
	// Crear incidentes
	for i := 0; i < 5; i++ {
		incident, _ := domain.NewIncident("Incident", "Desc", domain.SeverityHigh, user.ID)
		err := repo.Create(ctx, incident)
		require.NoError(t, err)
	}
	
	// Contar
	filters := domain.ListFilters{}
	count, err := repo.Count(ctx, filters)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}
