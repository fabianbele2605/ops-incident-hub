# Arquitectura Consolidada - Ops Incident Hub

## 1. Visión General

### Propósito del Sistema
Plataforma para gestión de incidentes operativos con seguimiento completo del ciclo de vida, diseñada con estándares de nivel senior para producción.

### Arquitectura de Alto Nivel
- **Patrón:** Clean Architecture (Hexagonal)
- **Estilo:** Monolito modular
- **Comunicación:** REST API
- **Persistencia:** PostgreSQL

### Principios Arquitectónicos
1. **Separación de responsabilidades** - Capas bien definidas
2. **Independencia de frameworks** - Domain no depende de infraestructura
3. **Testabilidad** - Cada capa es testeable independientemente
4. **Mantenibilidad** - Código limpio y documentado
5. **Escalabilidad** - Preparado para crecimiento

## 2. Stack Tecnológico

### Backend
- **Lenguaje:** Go 1.24
- **Arquitectura:** Clean Architecture
- **Base de datos:** PostgreSQL 15
- **Migraciones:** golang-migrate
- **Router:** gorilla/mux
- **UUID:** google/uuid

### Infraestructura
- **Contenedores:** Docker
- **Orquestación:** Docker Compose
- **CI/CD:** GitHub Actions
- **VCS:** Git + GitHub

### Observabilidad
- **Logging:** slog (Go standard library)
- **Métricas:** Prometheus (prometheus/client_golang)
- **Tracing:** Request ID propagation
- **Health Checks:** Custom implementation

### Seguridad
- **Headers:** OWASP recommendations
- **CORS:** Origin validation con wildcards
- **Rate Limiting:** golang.org/x/time/rate
- **Input Validation:** Manual validation

### Resiliencia
- **Circuit Breaker:** Custom implementation
- **Retry:** Exponential backoff
- **Graceful Shutdown:** Signal handling
- **Connection Pool:** sql.DB optimizado

## 3. Capas de Arquitectura

### Domain Layer (`internal/domain/`)
**Responsabilidad:** Lógica de negocio pura

**Componentes:**
- **Entidades:**
  - `Incident` - Incidente con estados y severidad
  - `User` - Usuario con roles
- **Errores tipados:**
  - `ErrIncidentNotFound`
  - `ErrUserNotFound`
  - `ErrForbidden`
- **Interfaces:**
  - `IncidentRepository`
  - `UserRepository`

**Características:**
- Sin dependencias externas
- Reglas de negocio centralizadas
- Validaciones en constructores

### Use Case Layer (`internal/usecase/`)
**Responsabilidad:** Orquestación de lógica de negocio

**Casos de Uso:**
- `CreateIncidentUseCase` - Crear incidente
- `AssignIncidentUseCase` - Asignar incidente
- `ListIncidentUseCase` - Listar incidentes

**Características:**
- Depende solo de domain
- Coordina repositorios
- Registra métricas de negocio

### Infrastructure Layer (`internal/infrastructure/`)
**Responsabilidad:** Implementaciones concretas

**Componentes:**
- **PostgreSQL:**
  - `IncidentRepository` - CRUD de incidentes
  - `UserRepository` - CRUD de usuarios
  - `Database` - Connection pooling
  - `Migrations` - Migraciones automáticas

**Características:**
- Implementa interfaces de domain
- Circuit breaker integrado
- Manejo de errores SQL

### API Layer (`internal/api/`)
**Responsabilidad:** Exposición HTTP

**Componentes:**
- **Handlers:**
  - `IncidentHandler` - Endpoints de incidentes
  - `HealthHandler` - Health checks
- **Middlewares:**
  - `SecurityHeadersMiddleware` - Headers OWASP
  - `CORSMiddleware` - CORS validation
  - `RateLimitMiddleware` - Rate limiting
  - `MetricsMiddleware` - Prometheus metrics
  - `LoggingMiddleware` - Request logging
- **Router:** Configuración de rutas

**Características:**
- DTOs para request/response
- Validación de entrada
- Manejo de errores HTTP

## 4. Decisiones Técnicas Clave

### Fase 2: Clean Architecture
**Decisión:** Implementar Clean Architecture

**Justificación:**
- Separación clara de responsabilidades
- Testabilidad mejorada
- Independencia de frameworks
- Mantenibilidad a largo plazo

**Impacto:** Base sólida para crecimiento

### Fase 3: PostgreSQL con Migraciones
**Decisión:** PostgreSQL con golang-migrate

**Justificación:**
- ACID compliance
- Relaciones complejas
- Migraciones versionadas
- Rollback automático

**Impacto:** Integridad de datos garantizada

### Fase 4: Docker Multi-Stage Build
**Decisión:** Dockerfile optimizado con scratch

**Justificación:**
- Imagen de ~10MB (95% reducción)
- Seguridad mejorada
- Deploy más rápido
- Menos superficie de ataque

**Impacto:** Eficiencia en producción

### Fase 5: Testing Strategy
**Decisión:** Tests en 3 niveles

**Justificación:**
- Unit tests para domain y use cases
- Integration tests para repositories
- Mocks compartidos para handlers
- CI/CD automatizado

**Impacto:** Confianza en cambios

### Fase 6: slog para Logging
**Decisión:** Usar slog en lugar de librerías externas

**Justificación:**
- Standard library (Go 1.21+)
- Sin dependencias externas
- Performance comparable a zap
- Soporte oficial

**Impacto:** Logs estructurados sin overhead

### Fase 7: Security Middleware Stack
**Decisión:** Middlewares de seguridad en orden específico

**Justificación:**
- Security → CORS → RateLimit → Metrics → Logging
- Cada middleware tiene responsabilidad única
- Orden importa para funcionalidad correcta

**Impacto:** Seguridad en capas

### Fase 8: Circuit Breaker en Repositorios
**Decisión:** Circuit breaker a nivel de repositorio

**Justificación:**
- Protege operaciones de DB
- Falla rápido cuando DB está caída
- Recuperación automática
- Thread-safe

**Impacto:** Resiliencia ante fallos de DB

## 5. Flujo de Datos

### Request Flow
```
HTTP Request
    ↓
[Security Middleware] → Headers OWASP
    ↓
[CORS Middleware] → Origin validation
    ↓
[RateLimit Middleware] → 10 req/s per IP
    ↓
[Metrics Middleware] → Prometheus metrics
    ↓
[Logging Middleware] → Request ID + logs
    ↓
[Router] → Route matching
    ↓
[Handler] → Parse request, validate
    ↓
[Use Case] → Business logic
    ↓
[Repository] → Circuit breaker → Retry → Database
    ↓
[Response] → JSON + status code
```

### Database Flow
```
Use Case
    ↓
Repository.Create()
    ↓
Circuit Breaker.Execute()
    ↓
[If Circuit CLOSED]
    ↓
Retry.Execute() (3 attempts, exponential backoff)
    ↓
db.ExecContext() → PostgreSQL
    ↓
[Success] → Circuit stays CLOSED
[Failure] → Circuit counts failure
[5 failures] → Circuit OPEN (30s)
```

## 6. Patrones Implementados

### Repository Pattern
**Ubicación:** `internal/domain/` (interfaces), `internal/infrastructure/postgres/` (implementación)

**Propósito:** Abstracción de persistencia

**Beneficios:**
- Testabilidad (mocks fáciles)
- Cambio de DB sin afectar domain
- Separación de responsabilidades

### Use Case Pattern
**Ubicación:** `internal/usecase/`

**Propósito:** Orquestación de lógica de negocio

**Beneficios:**
- Casos de uso explícitos
- Reutilización de lógica
- Testing independiente

### Middleware Pattern
**Ubicación:** `internal/api/middleware/`

**Propósito:** Cross-cutting concerns

**Beneficios:**
- Separación de responsabilidades
- Reutilización
- Composición flexible

### Circuit Breaker Pattern
**Ubicación:** `internal/resilience/circuitbreaker.go`

**Propósito:** Protección contra fallos en cascada

**Beneficios:**
- Falla rápido
- Recuperación automática
- Protección de recursos

### Retry Pattern
**Ubicación:** `internal/resilience/retry.go`

**Propósito:** Manejo de fallos temporales

**Beneficios:**
- Recuperación automática
- Backoff exponencial
- Context-aware

## 7. Integraciones

### PostgreSQL
- **Versión:** 15
- **Driver:** lib/pq
- **Connection Pool:** 25 open, 5 idle
- **Migraciones:** golang-migrate

### Prometheus
- **Endpoint:** `/metrics`
- **Métricas de negocio:** incidents_total, incidents_assigned_total
- **Métricas técnicas:** http_requests_total, http_request_duration_seconds

### GitHub Actions
- **Pipeline:** lint → test → security → build → docker-build
- **Triggers:** push, pull_request
- **Checks:** 6 validaciones automáticas

## 8. Configuración

### Variables de Entorno
```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_SHUTDOWN_TIMEOUT=30s
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=ops_incident_hub
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
DB_CONN_MAX_IDLE_TIME=1m

# Security
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Configuración por Entorno
- **Development:** Logs en texto, debug habilitado
- **Staging:** Logs en JSON, métricas habilitadas
- **Production:** Logs en JSON, todas las protecciones activas

### Secrets Management
- Variables de entorno para desarrollo
- Azure Key Vault para producción (futuro)

## 9. Métricas y Observabilidad

### Logging
- **Formato:** JSON (producción), texto (desarrollo)
- **Niveles:** DEBUG, INFO, WARN, ERROR
- **Request ID:** UUID único por request
- **Contexto:** Propagación vía context.Context

### Métricas
- **incidents_total** - Counter por severity
- **incidents_assigned_total** - Counter
- **http_requests_total** - Counter por method, endpoint, status
- **http_request_duration_seconds** - Histogram

### Health Checks
- `/health` - Estado completo con dependencias (JSON)
- `/health/live` - Liveness probe (200 OK)
- `/health/ready` - Readiness probe (200 OK / 503)

## 10. Seguridad

### Headers OWASP
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security: max-age=31536000
- Content-Security-Policy: default-src 'self'
- Referrer-Policy: strict-origin-when-cross-origin
- Permissions-Policy: geolocation=(), microphone=(), camera=()

### CORS
- Validación de origen
- Wildcard support (*.example.com)
- Preflight handling (OPTIONS)

### Rate Limiting
- 10 requests/segundo por IP
- Burst de 20 requests
- Cleanup automático cada 5 minutos

## 11. Resiliencia

### Circuit Breaker
- Max failures: 5
- Reset timeout: 30 segundos
- Estados: Closed → Open → Half-Open

### Retry Policies
- Max attempts: 3
- Initial wait: 100ms
- Max wait: 5s
- Multiplier: 2.0 (exponencial)

### Graceful Shutdown
- Timeout: 30 segundos
- Signal handling: SIGTERM, SIGINT
- Cierre ordenado de recursos

### Connection Pooling
- Max open: 25 conexiones
- Max idle: 5 conexiones
- Lifetime: 5 minutos
- Idle time: 1 minuto

## 12. Testing

### Unit Tests
- Domain entities
- Use cases
- Mocks para repositorios

### Integration Tests
- Repositorios con PostgreSQL real
- Testcontainers (futuro)

### Coverage
- 51 test cases
- 887+ líneas de tests
- Cobertura en capas críticas

## 13. CI/CD

### Pipeline
1. **Lint** - golangci-lint
2. **Test** - go test
3. **Security** - gosec, govulncheck, gitleaks
4. **Build** - go build
5. **Docker Build** - docker build

### Checks
- 6 validaciones automáticas
- 100% success rate
- Bloqueo de merge si falla

## 14. Deployment

### Docker
- Multi-stage build
- Imagen base: scratch
- Tamaño: ~10MB
- Non-root user (futuro)

### Docker Compose
- PostgreSQL + API
- Health checks configurados
- Restart policies
- Resource limits

## 15. Evolución Futura

### Próximos Pasos
- Container hardening
- E2E tests
- Distributed tracing
- Azure deployment

### Mejoras Planificadas
- Frontend (TypeScript + React)
- API documentation (OpenAPI)
- Performance optimization
- Multi-tenancy

---

**Documento creado:** 1 de marzo de 2025  
**Versión:** 1.0  
**Estado:** Consolidado hasta Fase 8
