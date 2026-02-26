# Fase 2 - Modelo de Dominio

## 1. Entidades del Dominio

Las entidades representan los conceptos principales del negocio. Son independientes de la infraestructura.

### 1.1 Incident (Incidente)

**Propósito:** Representa un incidente operativo que necesita ser gestionado.

**Archivo a crear:** `backend/internal/domain/incident.go`

**Contenido que debes escribir:**

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Severity representa la severidad de un incidente
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Status representa el estado de un incidente
type Status string

const (
	StatusOpen       Status = "open"
	StatusAssigned   Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

// Incident representa un incidente operativo
type Incident struct {
	ID          uuid.UUID
	Title       string
	Description string
	Severity    Severity
	Status      Status
	AssignedTo  *uuid.UUID // Puntero porque puede ser nil
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ResolvedAt  *time.Time // Puntero porque puede ser nil
	Metadata    map[string]interface{}
}

// NewIncident crea un nuevo incidente
func NewIncident(title, description string, severity Severity, createdBy uuid.UUID) (*Incident, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if !isValidSeverity(severity) {
		return nil, ErrInvalidSeverity
	}

	now := time.Now()
	return &Incident{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Severity:    severity,
		Status:      StatusOpen,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    make(map[string]interface{}),
	}, nil
}

// Assign asigna el incidente a un usuario
func (i *Incident) Assign(userID uuid.UUID) error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	i.AssignedTo = &userID
	i.Status = StatusAssigned
	i.UpdatedAt = time.Now()
	return nil
}

// StartProgress marca el incidente como en progreso
func (i *Incident) StartProgress() error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	if i.AssignedTo == nil {
		return ErrIncidentNotAssigned
	}
	i.Status = StatusInProgress
	i.UpdatedAt = time.Now()
	return nil
}

// Resolve marca el incidente como resuelto
func (i *Incident) Resolve() error {
	if i.Status == StatusClosed {
		return ErrIncidentClosed
	}
	now := time.Now()
	i.Status = StatusResolved
	i.ResolvedAt = &now
	i.UpdatedAt = now
	return nil
}

// Close cierra el incidente
func (i *Incident) Close() error {
	if i.Status != StatusResolved {
		return ErrIncidentNotResolved
	}
	i.Status = StatusClosed
	i.UpdatedAt = time.Now()
	return nil
}

// IsOpen verifica si el incidente está abierto
func (i *Incident) IsOpen() bool {
	return i.Status == StatusOpen || i.Status == StatusAssigned || i.Status == StatusInProgress
}

// isValidSeverity valida si la severidad es válida
func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow:
		return true
	default:
		return false
	}
}
```

**Explicación:**

- **Tipos personalizados:** `Severity` y `Status` son strings con valores constantes
- **Struct Incident:** Contiene todos los campos del incidente
- **Punteros:** `AssignedTo` y `ResolvedAt` son punteros porque pueden ser nil
- **Constructor:** `NewIncident` crea un incidente con validaciones
- **Métodos de negocio:** `Assign`, `StartProgress`, `Resolve`, `Close` encapsulan la lógica
- **Sin dependencias:** No depende de base de datos ni frameworks

---

### 1.2 Errors (Errores del Dominio)

**Propósito:** Definir errores específicos del dominio.

**Archivo a crear:** `backend/internal/domain/errors.go`

**Contenido que debes escribir:**

```go
package domain

import "errors"

var (
	// Incident errors
	ErrInvalidTitle          = errors.New("invalid incident title")
	ErrInvalidDescription    = errors.New("invalid incident description")
	ErrInvalidSeverity       = errors.New("invalid incident severity")
	ErrIncidentNotFound      = errors.New("incident not found")
	ErrIncidentClosed        = errors.New("incident is closed")
	ErrIncidentNotAssigned   = errors.New("incident is not assigned")
	ErrIncidentNotResolved   = errors.New("incident is not resolved")

	// User errors
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidName           = errors.New("invalid name")
	ErrInvalidRole           = errors.New("invalid role")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")

	// Generic errors
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInternalServer        = errors.New("internal server error")
)
```

**Explicación:**

- **Errores de dominio:** Específicos del negocio, no técnicos
- **Reutilizables:** Se usan en toda la aplicación
- **Claros:** Nombres descriptivos que explican el problema

---

### 1.3 User (Usuario)

**Propósito:** Representa un usuario del sistema.

**Archivo a crear:** `backend/internal/domain/user.go`

**Contenido que debes escribir:**

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role representa el rol de un usuario
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

// User representa un usuario del sistema
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Role      Role
	CreatedAt time.Time
}

// NewUser crea un nuevo usuario
func NewUser(email, name string, role Role) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: time.Now(),
	}, nil
}

// CanManageIncidents verifica si el usuario puede gestionar incidentes
func (u *User) CanManageIncidents() bool {
	return u.Role == RoleAdmin || u.Role == RoleOperator
}

// CanViewIncidents verifica si el usuario puede ver incidentes
func (u *User) CanViewIncidents() bool {
	return true // Todos los roles pueden ver
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// isValidRole valida si el rol es válido
func isValidRole(r Role) bool {
	switch r {
	case RoleAdmin, RoleOperator, RoleViewer:
		return true
	default:
		return false
	}
}
```

**Explicación:**

- **Role:** Define los permisos del usuario
- **Métodos de autorización:** `CanManageIncidents`, `CanViewIncidents`, `IsAdmin`
- **Validaciones:** En el constructor `NewUser`
- **Simple:** Solo lo esencial para el dominio

---

## 2. Orden de Creación

Crea los archivos en este orden:

1. **`backend/internal/domain/errors.go`** (primero, porque los otros lo usan)
2. **`backend/internal/domain/incident.go`**
3. **`backend/internal/domain/user.go`**

---

## 3. Dependencias Necesarias

Antes de crear los archivos, necesitas inicializar el módulo Go:

**Comandos que debes ejecutar:**

```bash
cd backend

# Inicializar módulo Go
go mod init github.com/fabianbele2605/ops-incident-hub/backend

# Agregar dependencia de UUID
go get github.com/google/uuid
```

**Explicación:**
- `go mod init`: Crea el archivo `go.mod` que gestiona dependencias
- `go get`: Descarga la librería de UUID de Google

---

## 4. Validación

Después de crear los archivos, valida que compilan:

```bash
cd backend

# Verificar que no hay errores de sintaxis
go build ./internal/domain/...

# Si todo está bien, no debería mostrar errores
```

---

## 5. Próximos Pasos

Una vez creados estos archivos:

1. Crear interfaces de repositorios
2. Crear casos de uso
3. Crear handlers HTTP

---

## 6. Características del Modelo de Dominio

✅ **Sin dependencias externas:** Solo usa tipos de Go estándar y UUID  
✅ **Lógica de negocio encapsulada:** Los métodos validan y cambian estado  
✅ **Inmutabilidad parcial:** Los campos no se modifican directamente  
✅ **Validaciones:** En constructores y métodos  
✅ **Errores tipados:** Errores específicos del dominio  

---

**¿Listo para crear estos archivos? Te guío paso a paso.**
