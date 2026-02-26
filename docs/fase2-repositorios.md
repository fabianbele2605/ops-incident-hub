# Fase 2 - Interfaces de Repositorios

## 1. ¿Qué son las interfaces de repositorios?

Son **contratos** que definen cómo acceder a los datos, sin especificar la implementación (PostgreSQL, MongoDB, etc).

**Ventajas:**
- Independencia de la base de datos
- Fácil de testear (mocks)
- Cumple con Clean Architecture

---

## 2. Archivo: `backend/internal/domain/repository.go`

**Contenido que debes escribir:**

```go
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
```

**Explicación:**

- **context.Context:** Para cancelación y timeouts
- **Punteros en filtros:** Permiten distinguir entre "no filtrar" y "filtrar por valor cero"
- **Métodos CRUD:** Create, Read (GetByID, List), Update, Delete
- **Métodos específicos:** GetByEmail para usuarios, Count para paginación
- **ListFilters:** Struct para filtros de búsqueda

---

## 3. Validar compilación

```bash
cd /mnt/n/ubuntu/azureSenior/backend
go build ./internal/domain/...
```

**Si no hay errores, está perfecto** ✅

---

## 4. Próximos pasos

Una vez creado `repository.go`:

1. Crear casos de uso (use cases)
2. Implementar repositorios con PostgreSQL
3. Crear handlers HTTP

---

## 5. Características de las interfaces

✅ **Independientes de implementación:** No saben si es PostgreSQL, MongoDB, etc  
✅ **Testeables:** Fácil crear mocks  
✅ **Context-aware:** Soportan cancelación y timeouts  
✅ **Paginación:** List con Limit/Offset  
✅ **Filtros flexibles:** Punteros permiten filtros opcionales  

---

**Crea el archivo `backend/internal/domain/repository.go` con el contenido de arriba.**
