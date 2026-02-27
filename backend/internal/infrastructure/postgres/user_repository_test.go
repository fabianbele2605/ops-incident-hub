package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	user, err := domain.NewUser("test@example.com", "Test User", domain.RoleAdmin)
	require.NoError(t, err)

	err = repo.Create(ctx, user)
	require.NoError(t, err)

	// Verificar que se creó
	retrieved, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.Role, retrieved.Role)
}

func TestUserRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Crear usuario
	user, _ := domain.NewUser("test@example.com", "Test User", domain.RoleOperator)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Obtener por ID
	retrieved, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	retrieved, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Equal(t, domain.ErrUserNotFound, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Crear usuario
	user, _ := domain.NewUser("unique@example.com", "Test User", domain.RoleViewer)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Obtener por email
	retrieved, err := repo.GetByEmail(ctx, "unique@example.com")
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, "unique@example.com", retrieved.Email)
}

func TestUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Crear usuario
	user, _ := domain.NewUser("test@example.com", "Original Name", domain.RoleViewer)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Actualizar
	user.Name = "Updated Name"
	user.Role = domain.RoleOperator
	err = repo.Update(ctx, user)
	require.NoError(t, err)

	// Verificar actualización
	retrieved, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
	assert.Equal(t, domain.RoleOperator, retrieved.Role)
}

func TestUserRepository_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Crear varios usuarios
	for i := 1; i <= 3; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		user, _ := domain.NewUser(email, "User", domain.RoleViewer)
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Listar
	users, err := repo.List(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 3)
}

func TestUserRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Crear usuario
	user, _ := domain.NewUser("delete@example.com", "To Delete", domain.RoleViewer)
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Eliminar
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Verificar que no existe
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Equal(t, domain.ErrUserNotFound, err)
}
