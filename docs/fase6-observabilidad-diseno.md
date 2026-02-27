# Fase 6 - Observabilidad y Operación: Diseño e Implementación

## 📋 Información General

**Fase:** 6 - Observabilidad y Operación  
**Objetivo:** Tener visibilidad completa de salud, rendimiento y fallos  
**Principio clave:** "Lo que no se observa, no se puede operar"

## 🎯 Objetivos de la Fase

1. **Logging estructurado:** Logs en formato JSON con contexto completo
2. **Métricas:** Indicadores de negocio y técnicos para monitoreo
3. **Trazabilidad:** Seguimiento de requests a través de capas
4. **Alertas:** Detección temprana de problemas con umbrales claros

## 📐 Arquitectura de Observabilidad

### Los 3 Pilares de Observabilidad

```
┌─────────────────────────────────────────────────────┐
│              OBSERVABILIDAD                         │
├─────────────────┬─────────────────┬─────────────────┤
│     LOGS        │    METRICS      │     TRACES      │
│  (Qué pasó)     │  (Cuánto/Cuándo)│  (Dónde/Cómo)   │
└─────────────────┴─────────────────┴─────────────────┘
```

### Stack Tecnológico Propuesto

Para este proyecto usaremos herramientas nativas de Go y estándares de la industria:

1. **Logging:** `slog` (Go 1.21+ standard library)
2. **Metrics:** Prometheus + `promhttp` (Go client)
3. **Traces:** OpenTelemetry (opcional para fase avanzada)
4. **Visualization:** Prometheus (queries) + logs estructurados

**Justificación:**
- `slog` es el estándar oficial de Go desde 1.21
- Prometheus es el estándar de facto para métricas en cloud-native
- No requiere dependencias externas pesadas
- Fácil integración con Azure Monitor en el futuro

## 🔍 Paso 1: Logging Estructurado

### Objetivo
Implementar logging estructurado con `slog` en todas las capas de la aplicación.

### Características del Logger

**Formato:** JSON estructurado
```json
{
  "time": "2025-02-27T10:30:00Z",
  "level": "INFO",
  "msg": "incident created",
  "incident_id": "123e4567-e89b-12d3-a456-426614174000",
  "severity": "high",
  "created_by": "user-123",
  "duration_ms": 45
}
```

**Niveles de log:**
- `DEBUG`: Información detallada para debugging
- `INFO`: Eventos normales de operación
- `WARN`: Situaciones anómalas pero manejables
- `ERROR`: Errores que requieren atención

**Contexto obligatorio:**
- `request_id`: ID único por request HTTP
- `user_id`: Usuario que ejecuta la acción (si aplica)
- `operation`: Nombre de la operación (create_incident, assign_incident, etc.)
- `duration_ms`: Tiempo de ejecución

### Estructura de Implementación

```
backend/internal/observability/
├── logger/
│   ├── logger.go          # Configuración de slog
│   ├── context.go         # Helpers para contexto
│   └── middleware.go      # Middleware HTTP para logging
```

### Qué Implementar

#### 1. `backend/internal/observability/logger/logger.go`

**Responsabilidad:** Configurar y crear instancia de logger

**Funciones principales:**
- `New(env string) *slog.Logger` - Crear logger según entorno
- `WithContext(ctx context.Context, logger *slog.Logger) context.Context` - Agregar logger a contexto
- `FromContext(ctx context.Context) *slog.Logger` - Obtener logger desde contexto

**Configuración por entorno:**
- `development`: Logs en formato texto, nivel DEBUG
- `production`: Logs en formato JSON, nivel INFO

#### 2. `backend/internal/observability/logger/middleware.go`

**Responsabilidad:** Middleware HTTP para logging automático

**Funcionalidades:**
- Generar `request_id` único por request
- Loggear inicio de request (método, path, IP)
- Loggear fin de request (status, duration)
- Agregar logger con contexto al request context
- Capturar panics y loggear como ERROR

**Estructura del log:**
```json
{
  "time": "2025-02-27T10:30:00Z",
  "level": "INFO",
  "msg": "http request completed",
  "request_id": "req-abc123",
  "method": "POST",
  "path": "/api/v1/incidents",
  "status": 201,
  "duration_ms": 45,
  "ip": "192.168.1.1"
}
```

#### 3. Integración en Capas

**Handler Layer:**
```go
// Obtener logger desde contexto
logger := logger.FromContext(r.Context())

// Loggear con contexto
logger.Info("creating incident",
    "title", req.Title,
    "severity", req.Severity,
    "created_by", req.CreatedBy,
)
```

**Use Case Layer:**
```go
// Recibir contexto y loggear operaciones
func (uc *CreateIncidentUseCase) Execute(ctx context.Context, input CreateIncidentInput) (*Incident, error) {
    logger := logger.FromContext(ctx)
    
    logger.Debug("validating input")
    // ... validación
    
    logger.Info("creating incident in repository")
    // ... crear incidente
    
    logger.Info("incident created successfully", "incident_id", incident.ID)
    return incident, nil
}
```

**Repository Layer:**
```go
// Loggear queries y errores
logger.Debug("executing query", "query", "INSERT INTO incidents...")
if err != nil {
    logger.Error("database error", "error", err, "operation", "create_incident")
}
```

### Criterios de Validación

- [x] Logger configurado con `slog`
- [x] Middleware HTTP genera `request_id`
- [x] Logs en formato JSON en producción
- [x] Logs incluyen contexto completo (request_id, user_id, operation)
- [x] Todas las capas usan logger desde contexto
- [x] Errores loggeados con nivel ERROR
- [x] Operaciones importantes loggeadas con nivel INFO

## 📊 Paso 2: Métricas con Prometheus

### Objetivo
Exponer métricas de negocio y técnicas en formato Prometheus.

### Tipos de Métricas

**1. Métricas de Negocio:**
- `incidents_total` (Counter): Total de incidentes creados por severidad
- `incidents_by_status` (Gauge): Incidentes actuales por estado
- `incident_resolution_time_seconds` (Histogram): Tiempo de resolución
- `incidents_assigned_total` (Counter): Total de asignaciones

**2. Métricas Técnicas:**
- `http_requests_total` (Counter): Total de requests HTTP por endpoint y status
- `http_request_duration_seconds` (Histogram): Duración de requests
- `database_queries_total` (Counter): Total de queries por operación
- `database_query_duration_seconds` (Histogram): Duración de queries

**3. Métricas de Sistema:**
- Go runtime metrics (automáticas con Prometheus client)
- Goroutines, memoria, GC, etc.

### Estructura de Implementación

```
backend/internal/observability/
├── metrics/
│   ├── metrics.go         # Definición de métricas
│   ├── collector.go       # Collectors personalizados
│   └── middleware.go      # Middleware HTTP para métricas
```

### Qué Implementar

#### 1. `backend/internal/observability/metrics/metrics.go`

**Responsabilidad:** Definir y registrar métricas

**Métricas a crear:**
```go
// Business metrics
var (
    IncidentsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "incidents_total",
            Help: "Total number of incidents created",
        },
        []string{"severity"},
    )
    
    IncidentsByStatus = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "incidents_by_status",
            Help: "Current number of incidents by status",
        },
        []string{"status"},
    )
    
    IncidentResolutionTime = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name: "incident_resolution_time_seconds",
            Help: "Time taken to resolve incidents",
            Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
        },
    )
)

// Technical metrics
var (
    HTTPRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
        },
        []string{"method", "endpoint"},
    )
)
```

**Funciones:**
- `Register()` - Registrar todas las métricas en Prometheus
- `RecordIncidentCreated(severity string)` - Helper para incrementar contador
- `RecordIncidentResolved(duration time.Duration)` - Helper para histogram

#### 2. `backend/internal/observability/metrics/middleware.go`

**Responsabilidad:** Middleware HTTP para métricas automáticas

**Funcionalidades:**
- Incrementar `http_requests_total` por cada request
- Medir duración con `http_request_duration_seconds`
- Normalizar endpoints (evitar cardinalidad alta)

#### 3. `backend/internal/api/handler/metrics_handler.go`

**Responsabilidad:** Endpoint `/metrics` para Prometheus

**Implementación:**
```go
func NewMetricsHandler() http.Handler {
    return promhttp.Handler()
}
```

#### 4. Integración en Use Cases

**CreateIncidentUseCase:**
```go
func (uc *CreateIncidentUseCase) Execute(ctx context.Context, input CreateIncidentInput) (*Incident, error) {
    // ... crear incidente
    
    // Registrar métrica
    metrics.RecordIncidentCreated(string(incident.Severity))
    
    return incident, nil
}
```

**AssignIncidentUseCase:**
```go
func (uc *AssignIncidentUseCase) Execute(ctx context.Context, input AssignIncidentInput) error {
    // ... asignar incidente
    
    // Registrar métrica
    metrics.IncidentsAssignedTotal.Inc()
    
    return nil
}
```

### Endpoint de Métricas

**URL:** `GET /metrics`

**Formato de respuesta:**
```
# HELP incidents_total Total number of incidents created
# TYPE incidents_total counter
incidents_total{severity="critical"} 5
incidents_total{severity="high"} 12
incidents_total{severity="medium"} 23
incidents_total{severity="low"} 8

# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="POST",endpoint="/api/v1/incidents",status="201"} 48
http_requests_total{method="GET",endpoint="/api/v1/incidents",status="200"} 156
```

### Criterios de Validación

- [x] Métricas de negocio definidas y registradas
- [x] Métricas técnicas definidas y registradas
- [x] Endpoint `/metrics` expuesto
- [x] Middleware HTTP registra métricas automáticamente
- [x] Use cases incrementan métricas de negocio
- [x] Formato Prometheus válido

## 🔗 Paso 3: Request ID y Trazabilidad

### Objetivo
Implementar trazabilidad de requests a través de todas las capas.

### Concepto

Cada request HTTP recibe un `request_id` único que se propaga:
1. Generado en middleware HTTP
2. Agregado al contexto
3. Incluido en todos los logs
4. Incluido en headers de respuesta
5. Propagado a servicios externos (futuro)

### Estructura de Implementación

```
backend/internal/observability/
├── tracing/
│   ├── request_id.go      # Generación y manejo de request_id
│   └── middleware.go      # Middleware para request_id
```

### Qué Implementar

#### 1. `backend/internal/observability/tracing/request_id.go`

**Responsabilidad:** Generar y manejar request_id

**Funciones:**
```go
// GenerateRequestID genera un request_id único
func GenerateRequestID() string {
    return fmt.Sprintf("req-%s", uuid.New().String()[:8])
}

// WithRequestID agrega request_id al contexto
func WithRequestID(ctx context.Context, requestID string) context.Context

// GetRequestID obtiene request_id desde contexto
func GetRequestID(ctx context.Context) string

// RequestIDFromHeader extrae request_id de header HTTP (si existe)
func RequestIDFromHeader(r *http.Request) string
```

#### 2. `backend/internal/observability/tracing/middleware.go`

**Responsabilidad:** Middleware para request_id

**Funcionalidades:**
- Extraer `X-Request-ID` de header (si existe)
- Generar nuevo request_id si no existe
- Agregar request_id al contexto
- Agregar request_id a response header
- Agregar request_id al logger

### Integración

**En main.go:**
```go
// Orden de middlewares (importante)
router.Use(tracing.RequestIDMiddleware())  // 1. Request ID
router.Use(logger.LoggingMiddleware())     // 2. Logging (usa request_id)
router.Use(metrics.MetricsMiddleware())    // 3. Metrics
```

**En logs:**
```json
{
  "time": "2025-02-27T10:30:00Z",
  "level": "INFO",
  "msg": "creating incident",
  "request_id": "req-abc12345",
  "incident_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

**En response headers:**
```
X-Request-ID: req-abc12345
```

### Criterios de Validación

- [x] Request ID generado por request
- [x] Request ID en contexto
- [x] Request ID en todos los logs
- [x] Request ID en response headers
- [x] Request ID extraído de header si existe

## 🚨 Paso 4: Health Checks Mejorados

### Objetivo
Mejorar health checks existentes con más información y métricas.

### Health Checks Actuales

Ya tenemos:
- `/health` - Health check básico
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe

### Mejoras a Implementar

#### 1. Agregar Información Detallada

**`/health` response mejorado:**
```json
{
  "status": "healthy",
  "timestamp": "2025-02-27T10:30:00Z",
  "version": "1.0.0",
  "uptime_seconds": 3600,
  "checks": {
    "database": {
      "status": "healthy",
      "response_time_ms": 5
    },
    "memory": {
      "status": "healthy",
      "alloc_mb": 45,
      "sys_mb": 72
    }
  }
}
```

#### 2. Agregar Métricas de Health

```go
var (
    HealthCheckTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "health_check_total",
            Help: "Total number of health checks",
        },
        []string{"endpoint", "status"},
    )
    
    HealthCheckDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "health_check_duration_seconds",
            Help: "Health check duration",
        },
        []string{"endpoint"},
    )
)
```

### Criterios de Validación

- [x] Health checks retornan información detallada
- [x] Health checks incluyen timestamp y uptime
- [x] Health checks verifican database
- [x] Health checks exponen métricas

## 📦 Resumen de Archivos a Crear

### Nuevos Archivos

```
backend/internal/observability/
├── logger/
│   ├── logger.go          # Configuración de slog
│   ├── context.go         # Helpers para contexto
│   └── middleware.go      # Middleware HTTP
├── metrics/
│   ├── metrics.go         # Definición de métricas
│   ├── helpers.go         # Helper functions
│   └── middleware.go      # Middleware HTTP
└── tracing/
    ├── request_id.go      # Manejo de request_id
    └── middleware.go      # Middleware HTTP
```

### Archivos a Modificar

```
backend/
├── cmd/api/main.go                        # Integrar observabilidad
├── internal/api/router.go                 # Agregar middlewares y /metrics
├── internal/api/handler/health_handler.go # Mejorar health checks
├── internal/api/handler/incident_handler.go # Agregar logging
├── internal/usecase/incident/*.go         # Agregar logging y métricas
└── internal/infrastructure/postgres/*.go  # Agregar logging
```

## 🔄 Orden de Implementación Recomendado

### Paso 1: Logging Estructurado (PR #21)
1. Crear `internal/observability/logger/`
2. Implementar logger con slog
3. Crear middleware de logging
4. Integrar en main.go
5. Agregar logging en handlers
6. Agregar logging en use cases
7. Agregar logging en repositories

**Validación:** Logs en formato JSON con request_id

### Paso 2: Request ID y Trazabilidad (PR #22)
1. Crear `internal/observability/tracing/`
2. Implementar generación de request_id
3. Crear middleware de request_id
4. Integrar con logger
5. Agregar request_id a response headers

**Validación:** Request ID en logs y headers

### Paso 3: Métricas Prometheus (PR #23)
1. Crear `internal/observability/metrics/`
2. Definir métricas de negocio
3. Definir métricas técnicas
4. Crear middleware de métricas
5. Agregar endpoint `/metrics`
6. Integrar métricas en use cases

**Validación:** Endpoint `/metrics` funcional

### Paso 4: Health Checks Mejorados (PR #24)
1. Mejorar `health_handler.go`
2. Agregar información detallada
3. Agregar métricas de health
4. Agregar verificación de componentes

**Validación:** Health checks con información completa

## ✅ Criterios de Éxito de la Fase

Al finalizar Fase 6, el proyecto debe tener:

### Logging
- [x] Logs estructurados en formato JSON
- [x] Request ID en todos los logs
- [x] Niveles de log apropiados (DEBUG, INFO, WARN, ERROR)
- [x] Contexto completo en logs (user_id, operation, duration)
- [x] Logs en todas las capas (handler, use case, repository)

### Métricas
- [x] Endpoint `/metrics` expuesto
- [x] Métricas de negocio (incidents_total, resolution_time, etc.)
- [x] Métricas técnicas (http_requests, duration, etc.)
- [x] Métricas de sistema (Go runtime)
- [x] Formato Prometheus válido

### Trazabilidad
- [x] Request ID único por request
- [x] Request ID propagado en contexto
- [x] Request ID en logs
- [x] Request ID en response headers

### Health Checks
- [x] Health checks con información detallada
- [x] Verificación de componentes (database, memory)
- [x] Métricas de health checks

### Operación
- [x] Capacidad de debugging con logs
- [x] Capacidad de monitoreo con métricas
- [x] Capacidad de trazabilidad con request_id
- [x] Detección temprana de problemas

## 🎯 Entregables de la Fase

1. **Código:**
   - Implementación completa de observabilidad
   - Tests unitarios de componentes de observabilidad
   - Integración en todas las capas

2. **Documentación:**
   - Este documento de diseño
   - Documentación de implementación por paso
   - Guía de operación (cómo usar logs y métricas)
   - Documento de resumen de fase

3. **Evidencia:**
   - Logs estructurados funcionando
   - Endpoint `/metrics` con datos reales
   - Request ID en logs y headers
   - Health checks mejorados

## 📚 Referencias Técnicas

- [Go slog documentation](https://pkg.go.dev/log/slog)
- [Prometheus Go client](https://github.com/prometheus/client_golang)
- [Prometheus best practices](https://prometheus.io/docs/practices/naming/)
- [The Twelve-Factor App - Logs](https://12factor.net/logs)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)

## 🚀 Próximos Pasos

Una vez completada esta fase, estarás listo para:
- **Fase 7:** Seguridad integral (secrets, identidades, hardening)
- **Fase 8:** Resiliencia (backups, disaster recovery)
- Integrar con Azure Monitor (Application Insights)
- Agregar dashboards con Grafana
- Configurar alertas basadas en métricas

---

**Documento creado:** 27 de febrero de 2025  
**Fase:** 6 - Observabilidad y Operación  
**Estado:** Diseño completo - Listo para implementación  
**Siguiente acción:** Developer implementa Paso 1 (Logging Estructurado)
