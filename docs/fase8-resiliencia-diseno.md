# Fase 8 - Resiliencia y Continuidad: Diseño e Implementación

## 📋 Información General

**Fase:** 8 - Resiliencia y Continuidad  
**Objetivo:** Garantizar continuidad operativa ante fallos técnicos y errores humanos  
**Principio clave:** "Lo que no se puede recuperar, no está listo para producción"

## 🎯 Objetivos de la Fase

1. **Graceful Shutdown:** Apagado limpio sin pérdida de requests en proceso
2. **Circuit Breaker:** Protección contra fallos en cascada
3. **Retry Policies:** Reintentos inteligentes con backoff exponencial
4. **Health Checks Avanzados:** Detección temprana de degradación
5. **Database Resilience:** Manejo de desconexiones y reconexión automática

## 🔄 Alcance Adaptado

**Nota:** Esta fase se enfoca en resiliencia de aplicación. Los backups de PostgreSQL y disaster recovery se implementarán cuando se despliegue en Azure con servicios gestionados.

### Implementación Actual (Local/Docker)
- ✅ Graceful shutdown con signal handling
- ✅ Circuit breaker para dependencias externas
- ✅ Retry policies con backoff exponencial
- ✅ Health checks con degradación gradual
- ✅ Database connection pooling y retry

### Implementación Futura (Azure)
- ⏳ Azure Database for PostgreSQL con backups automáticos
- ⏳ Geo-replication para disaster recovery
- ⏳ Azure Site Recovery
- ⏳ Backup policies con retención configurable

## 🛡️ Paso 1: Graceful Shutdown

### Objetivo
Apagar la aplicación de forma limpia, permitiendo que requests en proceso terminen antes de cerrar.

### Problema Actual
Cuando se detiene el contenedor (SIGTERM), la aplicación se cierra inmediatamente, potencialmente cortando requests en proceso.

### Solución: Signal Handling

#### Qué Implementar

**`backend/cmd/api/main.go`** - Modificar para capturar señales

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // ... configuración existente ...
    
    // Crear servidor HTTP
    server := &http.Server{
        Addr:         cfg.Server.Address,
        Handler:      router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }
    
    // Canal para errores del servidor
    serverErrors := make(chan error, 1)
    
    // Iniciar servidor en goroutine
    go func() {
        logger.Info("starting server", "address", cfg.Server.Address)
        serverErrors <- server.ListenAndServe()
    }()
    
    // Canal para señales del sistema
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
    
    // Esperar señal o error
    select {
    case err := <-serverErrors:
        logger.Error("server error", "error", err)
        os.Exit(1)
        
    case sig := <-shutdown:
        logger.Info("shutdown signal received", "signal", sig)
        
        // Timeout para graceful shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        // Intentar shutdown graceful
        if err := server.Shutdown(ctx); err != nil {
            logger.Error("graceful shutdown failed", "error", err)
            server.Close()
        }
        
        // Cerrar conexión de base de datos
        if err := db.Close(); err != nil {
            logger.Error("database close error", "error", err)
        }
        
        logger.Info("server stopped gracefully")
    }
}
```

### Comportamiento Esperado

1. **SIGTERM recibido** → Servidor deja de aceptar nuevas conexiones
2. **Requests en proceso** → Se les da 30 segundos para terminar
3. **Timeout alcanzado** → Cierre forzado
4. **Database** → Conexiones cerradas limpiamente

### Validación

```bash
# Terminal 1: Iniciar servidor
docker-compose up

# Terminal 2: Enviar requests continuos
while true; do curl http://localhost:8081/health; sleep 0.5; done

# Terminal 3: Enviar SIGTERM
docker-compose stop

# Verificar en logs:
# - "shutdown signal received"
# - Requests en proceso terminan
# - "server stopped gracefully"
```

## 🔌 Paso 2: Circuit Breaker Pattern

### Objetivo
Prevenir fallos en cascada cuando una dependencia externa falla repetidamente.

### Concepto

**Estados del Circuit Breaker:**
- **Closed:** Funcionamiento normal, requests pasan
- **Open:** Demasiados fallos, requests fallan rápido sin intentar
- **Half-Open:** Periodo de prueba, algunos requests pasan para verificar recuperación

### Estructura de Implementación

```
backend/internal/resilience/
├── circuitbreaker.go    # Circuit breaker implementation
└── circuitbreaker_test.go
```

### Qué Implementar

#### 1. `backend/internal/resilience/circuitbreaker.go`

```go
package resilience

import (
    "errors"
    "sync"
    "time"
)

var (
    ErrCircuitOpen = errors.New("circuit breaker is open")
)

// State representa el estado del circuit breaker
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

// CircuitBreaker implementa el patrón circuit breaker
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    
    mu           sync.RWMutex
    state        State
    failures     int
    lastFailTime time.Time
}

// NewCircuitBreaker crea un nuevo circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        StateClosed,
    }
}

// Execute ejecuta una función con protección de circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
    // Verificar estado
    if !cb.canExecute() {
        return ErrCircuitOpen
    }
    
    // Ejecutar función
    err := fn()
    
    // Registrar resultado
    cb.recordResult(err)
    
    return err
}

// canExecute verifica si se puede ejecutar la función
func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        // Verificar si es tiempo de intentar de nuevo
        if time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.mu.Unlock()
            cb.mu.RLock()
            return true
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

// recordResult registra el resultado de la ejecución
func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        // Éxito: resetear
        cb.failures = 0
        cb.state = StateClosed
    }
}

// State retorna el estado actual
func (cb *CircuitBreaker) State() State {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}
```

### Uso en Repositorios

**Ejemplo en `incident_repository.go`:**

```go
type IncidentRepository struct {
    db             *sql.DB
    circuitBreaker *resilience.CircuitBreaker
}

func NewIncidentRepository(db *sql.DB) *IncidentRepository {
    return &IncidentRepository{
        db:             db,
        circuitBreaker: resilience.NewCircuitBreaker(5, 30*time.Second),
    }
}

func (r *IncidentRepository) Create(ctx context.Context, incident *domain.Incident) error {
    return r.circuitBreaker.Execute(func() error {
        // Lógica de creación existente
        query := `INSERT INTO incidents ...`
        _, err := r.db.ExecContext(ctx, query, ...)
        return err
    })
}
```

## 🔄 Paso 3: Retry Policies con Backoff Exponencial

### Objetivo
Reintentar operaciones fallidas de forma inteligente, con esperas crecientes entre intentos.

### Estructura de Implementación

```
backend/internal/resilience/
├── retry.go           # Retry logic con backoff
└── retry_test.go
```

### Qué Implementar

#### 1. `backend/internal/resilience/retry.go`

```go
package resilience

import (
    "context"
    "errors"
    "math"
    "time"
)

// RetryConfig configuración de retry
type RetryConfig struct {
    MaxAttempts int
    InitialWait time.Duration
    MaxWait     time.Duration
    Multiplier  float64
}

// DefaultRetryConfig configuración por defecto
func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxAttempts: 3,
        InitialWait: 100 * time.Millisecond,
        MaxWait:     5 * time.Second,
        Multiplier:  2.0,
    }
}

// IsRetryable función que determina si un error es reintentable
type IsRetryable func(error) bool

// Retry ejecuta una función con reintentos y backoff exponencial
func Retry(ctx context.Context, cfg RetryConfig, isRetryable IsRetryable, fn func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
        // Ejecutar función
        err := fn()
        
        // Si no hay error, retornar éxito
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Si no es reintentable, retornar error inmediatamente
        if isRetryable != nil && !isRetryable(err) {
            return err
        }
        
        // Si es el último intento, no esperar
        if attempt == cfg.MaxAttempts-1 {
            break
        }
        
        // Calcular tiempo de espera con backoff exponencial
        wait := calculateBackoff(attempt, cfg)
        
        // Esperar con respeto al contexto
        select {
        case <-time.After(wait):
            // Continuar con siguiente intento
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return lastErr
}

// calculateBackoff calcula el tiempo de espera con backoff exponencial
func calculateBackoff(attempt int, cfg RetryConfig) time.Duration {
    wait := float64(cfg.InitialWait) * math.Pow(cfg.Multiplier, float64(attempt))
    
    if wait > float64(cfg.MaxWait) {
        wait = float64(cfg.MaxWait)
    }
    
    return time.Duration(wait)
}

// IsTemporaryError verifica si un error es temporal
func IsTemporaryError(err error) bool {
    // Errores temporales comunes
    var tempErr interface{ Temporary() bool }
    if errors.As(err, &tempErr) {
        return tempErr.Temporary()
    }
    
    // Agregar más casos según necesidad
    return false
}
```

### Uso en Casos de Uso

**Ejemplo en `create_incident.go`:**

```go
func (uc *CreateIncidentUseCase) Execute(ctx context.Context, input CreateIncidentInput) (*domain.Incident, error) {
    // ... validaciones ...
    
    var incident *domain.Incident
    
    // Retry con backoff exponencial
    err := resilience.Retry(
        ctx,
        resilience.DefaultRetryConfig(),
        resilience.IsTemporaryError,
        func() error {
            var err error
            incident, err = uc.repo.Create(ctx, newIncident)
            return err
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    return incident, nil
}
```

## 🏥 Paso 4: Health Checks Avanzados

### Objetivo
Detectar degradación gradual del sistema antes de fallos completos.

### Mejoras a Implementar

#### 1. Health Check con Dependencias

**`backend/internal/api/handler/health_handler.go`** - Mejorar health checks

```go
package handler

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
)

type HealthHandler struct {
    db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

// HealthResponse respuesta de health check
type HealthResponse struct {
    Status       string            `json:"status"`
    Version      string            `json:"version"`
    Dependencies map[string]string `json:"dependencies"`
    Timestamp    time.Time         `json:"timestamp"`
}

// Health endpoint básico
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:       "healthy",
        Version:      "1.0.0",
        Dependencies: make(map[string]string),
        Timestamp:    time.Now(),
    }
    
    // Verificar base de datos
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    if err := h.db.PingContext(ctx); err != nil {
        response.Status = "unhealthy"
        response.Dependencies["database"] = "down"
        w.WriteHeader(http.StatusServiceUnavailable)
    } else {
        response.Dependencies["database"] = "up"
        w.WriteHeader(http.StatusOK)
    }
    
    json.NewEncoder(w).Encode(response)
}

// Live liveness probe (proceso vivo)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// Ready readiness probe (listo para recibir tráfico)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
    // Verificar que todas las dependencias estén listas
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    if err := h.db.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Database not ready"))
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ready"))
}
```

#### 2. Actualizar Router

**`backend/internal/api/router.go`:**

```go
// Health checks
healthHandler := handler.NewHealthHandler(db)
router.HandleFunc("/health", healthHandler.Health).Methods("GET")
router.HandleFunc("/health/live", healthHandler.Live).Methods("GET")
router.HandleFunc("/health/ready", healthHandler.Ready).Methods("GET")
```

## 🗄️ Paso 5: Database Resilience

### Objetivo
Manejar desconexiones de base de datos y reconectar automáticamente.

### Mejoras a Implementar

#### 1. Connection Pooling Optimizado

**`backend/internal/infrastructure/database/postgres.go`:**

```go
package database

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq"
)

type Config struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func Connect(cfg Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
    )
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configurar connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)       // Máximo de conexiones abiertas
    db.SetMaxIdleConns(cfg.MaxIdleConns)       // Máximo de conexiones idle
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime) // Tiempo de vida máximo
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // Tiempo idle máximo
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return db, nil
}

// DefaultConfig configuración por defecto
func DefaultConfig() Config {
    return Config{
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 5 * time.Minute,
        ConnMaxIdleTime: 1 * time.Minute,
    }
}
```

#### 2. Actualizar Config

**`backend/internal/config/config.go`:**

```go
type DatabaseConfig struct {
    Host            string        `env:"DB_HOST" envDefault:"localhost"`
    Port            int           `env:"DB_PORT" envDefault:"5432"`
    User            string        `env:"DB_USER" envDefault:"postgres"`
    Password        string        `env:"DB_PASSWORD" envDefault:"postgres"`
    Database        string        `env:"DB_NAME" envDefault:"ops_incident_hub"`
    MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
    MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
    ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
    ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"1m"`
}
```

## 📦 Resumen de Archivos a Crear/Modificar

### Archivos a Crear

```
backend/internal/resilience/
├── circuitbreaker.go          # Circuit breaker pattern
├── circuitbreaker_test.go     # Tests de circuit breaker
├── retry.go                   # Retry con backoff exponencial
└── retry_test.go              # Tests de retry

backend/internal/infrastructure/database/
└── postgres.go                # Connection pooling optimizado
```

### Archivos a Modificar

```
backend/
├── cmd/api/main.go                         # Graceful shutdown
├── internal/config/config.go               # Database config extendida
├── internal/api/handler/health_handler.go  # Health checks mejorados
├── internal/api/router.go                  # Nuevos endpoints de health
└── internal/infrastructure/repository/     # Integrar circuit breaker y retry
    ├── incident_repository.go
    └── user_repository.go
```

## 🔄 Orden de Implementación Recomendado

### Paso 1: Graceful Shutdown (PR #27)
1. Modificar `main.go` con signal handling
2. Probar con docker-compose stop
3. Verificar logs de shutdown limpio

**Validación:** Servidor se detiene sin cortar requests

### Paso 2: Circuit Breaker (PR #28)
1. Crear `resilience/circuitbreaker.go`
2. Agregar tests unitarios
3. Integrar en repositorios
4. Probar con database down

**Validación:** Fallos rápidos cuando circuit está open

### Paso 3: Retry Policies (PR #29)
1. Crear `resilience/retry.go`
2. Agregar tests unitarios
3. Integrar en casos de uso
4. Probar con fallos temporales

**Validación:** Reintentos automáticos con backoff

### Paso 4: Health Checks Avanzados (PR #30)
1. Mejorar `health_handler.go`
2. Agregar verificación de dependencias
3. Actualizar router
4. Probar con database down

**Validación:** Health checks reflejan estado real

### Paso 5: Database Resilience (PR #31)
1. Crear `database/postgres.go`
2. Configurar connection pooling
3. Actualizar config
4. Probar con desconexiones

**Validación:** Reconexión automática funciona

## ✅ Criterios de Éxito de la Fase

Al finalizar Fase 8, el proyecto debe tener:

### Graceful Shutdown
- [x] Signal handling (SIGTERM, SIGINT)
- [x] Timeout configurable (30s)
- [x] Cierre limpio de conexiones
- [x] Logs de shutdown

### Circuit Breaker
- [x] Estados: Closed, Open, Half-Open
- [x] Configuración de thresholds
- [x] Fallos rápidos cuando open
- [x] Recuperación automática

### Retry Policies
- [x] Backoff exponencial
- [x] Máximo de intentos configurable
- [x] Respeto a contexto
- [x] Detección de errores temporales

### Health Checks
- [x] /health con dependencias
- [x] /health/live (liveness)
- [x] /health/ready (readiness)
- [x] Timeouts configurados

### Database Resilience
- [x] Connection pooling optimizado
- [x] Configuración de límites
- [x] Timeouts de conexión
- [x] Reconexión automática

## 🎯 Entregables de la Fase

1. **Código:**
   - Graceful shutdown
   - Circuit breaker
   - Retry policies
   - Health checks avanzados
   - Database resilience

2. **Documentación:**
   - Este documento de diseño
   - Documentación de implementación
   - Runbook de operación
   - Documento de resumen

3. **Evidencia:**
   - Shutdown limpio sin pérdida de requests
   - Circuit breaker funcionando
   - Reintentos con backoff
   - Health checks con dependencias
   - Connection pooling optimizado

## 📚 Referencias Técnicas

- [Graceful Shutdown in Go](https://golang.org/pkg/os/signal/)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [Exponential Backoff](https://en.wikipedia.org/wiki/Exponential_backoff)
- [Database Connection Pooling](https://go.dev/doc/database/manage-connections)
- [Health Check Patterns](https://microservices.io/patterns/observability/health-check-api.html)

## 🚀 Próximos Pasos

Una vez completada esta fase, estarás listo para:
- **Fase 9:** Gobierno y costos
- **Fase 10:** Cierre profesional y portafolio
- Despliegue en Azure con alta disponibilidad

## 📝 Backlog de Mejoras (Fuera de Alcance)

### Container Hardening (Fase 7 pendiente)
- Non-root user en Dockerfile
- Read-only filesystem
- Security options

### Input Validation Avanzada (Fase 7 pendiente)
- Validación de UUIDs en handlers
- Sanitización de strings
- Validación de enums

### Azure Integration (Futuro)
- Azure Database for PostgreSQL con backups
- Geo-replication
- Azure Site Recovery
- Backup policies

---

**Documento creado:** 1 de marzo de 2025  
**Fase:** 8 - Resiliencia y Continuidad  
**Estado:** Diseño completo - Listo para implementación  
**Siguiente acción:** Developer implementa Paso 1 (Graceful Shutdown)
