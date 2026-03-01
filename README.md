# 🚨 Ops Incident Hub

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.24-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Coverage](https://img.shields.io/badge/coverage-85%25-yellowgreen)
![Maturity](https://img.shields.io/badge/maturity-94%25%20SENIOR-success)
![PRs](https://img.shields.io/badge/PRs-34%20merged-blue)

> Plataforma profesional para gestión de incidentes operativos con arquitectura limpia, observabilidad completa y resiliencia robusta.

## 📋 Descripción

**Ops Incident Hub** es un sistema de gestión de incidentes diseñado con estándares de nivel senior para entornos de producción. Implementa Clean Architecture, observabilidad completa con Prometheus, resiliencia con Circuit Breaker y Retry Policies, y seguridad siguiendo recomendaciones OWASP.

### Valor de Negocio

- **Gestión centralizada** de incidentes operativos
- **Seguimiento completo** del ciclo de vida (Open → In Progress → Resolved → Closed)
- **Asignación inteligente** de incidentes a equipos
- **Observabilidad en tiempo real** con métricas y logs estructurados
- **Alta disponibilidad** con patrones de resiliencia

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                         API Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Handlers   │  │  Middleware  │  │    Router    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                       Use Case Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Create    │  │    Assign    │  │     List     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                       Domain Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Incident   │  │     User     │  │  Interfaces  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  PostgreSQL  │  │Circuit Breaker│ │    Retry     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### Patrones Implementados

- **Clean Architecture** - Separación de responsabilidades en 4 capas
- **Repository Pattern** - Abstracción de persistencia
- **Circuit Breaker** - Protección contra fallos en cascada
- **Retry Pattern** - Recuperación automática con backoff exponencial
- **Middleware Chain** - Security, CORS, Rate Limiting, Metrics, Logging

## 🛠️ Tech Stack

| Categoría | Tecnología | Versión |
|-----------|-----------|---------|
| **Backend** | Go | 1.24 |
| **Database** | PostgreSQL | 15 |
| **Observability** | Prometheus + slog | - |
| **Containers** | Docker + Docker Compose | - |
| **CI/CD** | GitHub Actions | - |
| **Testing** | Go testing + Testify | - |
| **Security** | gosec + govulncheck + gitleaks | - |

## 🚀 Quick Start

### Prerrequisitos

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (solo para desarrollo)

### Levantar el Sistema (3 comandos)

```bash
# 1. Clonar el repositorio
git clone https://github.com/fabianbele2605/ops-incident-hub.git
cd ops-incident-hub

# 2. Levantar servicios con Docker Compose
docker-compose up -d

# 3. Verificar que está funcionando
curl http://localhost:8080/health
```

**¡Listo!** El sistema está corriendo en `http://localhost:8080`

### Endpoints Disponibles

```bash
# Health Check
curl http://localhost:8080/health

# Crear Incidente
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database connection timeout",
    "description": "Production DB not responding",
    "severity": "high"
  }'

# Listar Incidentes
curl http://localhost:8080/api/v1/incidents

# Asignar Incidente
curl -X POST http://localhost:8080/api/v1/incidents/{id}/assign \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-uuid"}'

# Métricas Prometheus
curl http://localhost:8080/metrics
```

## ✨ Features Principales

### Gestión de Incidentes
- ✅ Crear incidentes con severidad (low, medium, high, critical)
- ✅ Asignar incidentes a usuarios
- ✅ Listar incidentes con filtros
- ✅ Estados del ciclo de vida (open, in_progress, resolved, closed)

### Observabilidad
- ✅ **Structured Logging** con slog (JSON en producción)
- ✅ **Métricas de negocio** (incidents_total, incidents_assigned_total)
- ✅ **Métricas técnicas** (http_requests_total, http_request_duration_seconds)
- ✅ **Request ID tracing** (UUID único por request)
- ✅ **Health Checks** (/health, /health/live, /health/ready)

### Seguridad
- ✅ **Security Headers OWASP** (7 headers implementados)
- ✅ **CORS** con validación de origen y wildcard support
- ✅ **Rate Limiting** (10 req/s por IP con burst de 20)
- ✅ **Input Validation** en todos los endpoints

### Resiliencia
- ✅ **Circuit Breaker** (5 fallos → 30s timeout)
- ✅ **Retry Policies** (3 intentos, backoff exponencial 100ms-5s)
- ✅ **Graceful Shutdown** (30s timeout)
- ✅ **Connection Pooling** optimizado (25 open, 5 idle)

### CI/CD
- ✅ **Pipeline automatizado** (lint, test, security, build)
- ✅ **Security Scanning** (gosec, govulncheck, gitleaks)
- ✅ **6 checks automáticos** en cada PR
- ✅ **100% success rate** en PRs mergeados

## 📁 Estructura del Proyecto

```
ops-incident-hub/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Entry point
│   ├── internal/
│   │   ├── domain/                  # Entidades y reglas de negocio
│   │   │   ├── incident.go
│   │   │   ├── user.go
│   │   │   └── errors.go
│   │   ├── usecase/                 # Casos de uso
│   │   │   ├── create_incident.go
│   │   │   ├── assign_incident.go
│   │   │   └── list_incident.go
│   │   ├── infrastructure/          # Implementaciones
│   │   │   └── postgres/
│   │   │       ├── incident_repository.go
│   │   │       ├── user_repository.go
│   │   │       └── database.go
│   │   ├── api/                     # HTTP Layer
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── router.go
│   │   ├── observability/           # Logging y métricas
│   │   │   ├── logger/
│   │   │   └── metrics/
│   │   └── resilience/              # Circuit breaker, retry
│   │       ├── circuitbreaker.go
│   │       └── retry.go
│   ├── migrations/                  # Database migrations
│   └── tests/                       # Tests (unit, integration)
├── docs/                            # Documentación técnica
│   ├── arquitectura-consolidada.md
│   ├── metricas-madurez.md
│   ├── backlog-mejoras.md
│   ├── guia-operacion.md
│   └── roadmap-tecnico.md
├── .github/
│   └── workflows/
│       └── ci.yml                   # GitHub Actions pipeline
├── docker-compose.yml               # Orquestación local
├── Dockerfile                       # Multi-stage build
└── README.md                        # Este archivo
```

## 📚 Documentación

### Documentación Técnica
- [Arquitectura Consolidada](docs/arquitectura-consolidada.md) - Decisiones técnicas y patrones
- [Métricas de Madurez](docs/metricas-madurez.md) - Evaluación 94% SENIOR
- [Guía de Operación](docs/guia-operacion.md) - Deployment y troubleshooting
- [Roadmap Técnico](docs/roadmap-tecnico.md) - Planificación Q2-Q4 2025
- [Backlog de Mejoras](docs/backlog-mejoras.md) - 10 mejoras priorizadas

### Guías de Deployment
- [Deployment Guide](docs/deployment-guide.md) - Local, AWS, Azure

### Demo y Presentación
- [Demo Script](docs/demo/demo-script.md) - Script para demo en vivo
- [Presentación Ejecutiva](docs/presentacion-ejecutiva.md) - 10 slides para portafolio

## 🧪 Testing

### Ejecutar Tests

```bash
# Unit tests
go test ./backend/internal/domain/... -v
go test ./backend/internal/usecase/... -v

# Integration tests
go test ./backend/internal/infrastructure/postgres/... -v

# Todos los tests
go test ./... -v

# Con coverage
go test ./... -cover
```

### Cobertura
- **51 test cases**
- **887+ líneas de tests**
- **Cobertura** en capas críticas (domain, use case, infrastructure)

## 🔒 Seguridad

### Security Headers (OWASP)
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### Security Scanning
- **gosec** - Static analysis security scanner
- **govulncheck** - Vulnerability detection
- **gitleaks** - Secrets detection

## 📊 Métricas y Monitoreo

### Prometheus Metrics

```bash
# Métricas de negocio
incidents_total{severity="high"} 42
incidents_assigned_total 38

# Métricas técnicas
http_requests_total{method="POST",endpoint="/api/v1/incidents",status="201"} 42
http_request_duration_seconds_bucket{le="0.1"} 95
```

### Health Checks

```bash
# Health completo (con dependencias)
GET /health
{
  "status": "healthy",
  "timestamp": "2025-03-01T10:00:00Z",
  "checks": {
    "database": "healthy"
  }
}

# Liveness probe
GET /health/live
200 OK

# Readiness probe
GET /health/ready
200 OK (o 503 si no está listo)
```

## 🔧 Configuración

### Variables de Entorno

```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=ops_incident_hub
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Security
ALLOWED_ORIGINS=https://example.com,https://*.example.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

Ver [Guía de Operación](docs/guia-operacion.md) para configuración completa.

## 🚀 Deployment

### Local (Docker Compose)
```bash
docker-compose up -d
```

### AWS (ECS + RDS)
```bash
cd infrastructure/terraform/aws
terraform init
terraform apply
```

### Azure (Container Apps + PostgreSQL)
```bash
cd infrastructure/terraform/azure
terraform init
terraform apply
```

Ver [Deployment Guide](docs/deployment-guide.md) para instrucciones detalladas.

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Por favor lee:

- [CONTRIBUTING.md](CONTRIBUTING.md) - Guía de contribución
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) - Código de conducta

### Proceso de Contribución

1. Fork el repositorio
2. Crea una rama feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📈 Roadmap

### Q2 2025
- ✅ Container hardening (non-root user)
- ✅ E2E tests con Testcontainers
- ✅ Input validation avanzada

### Q3 2025
- 🔄 Distributed tracing (OpenTelemetry)
- 🔄 Grafana dashboards
- 🔄 API documentation (OpenAPI)

### Q4 2025
- 📅 Frontend (TypeScript + React)
- 📅 Azure deployment
- 📅 Performance optimization

Ver [Roadmap Técnico](docs/roadmap-tecnico.md) para detalles completos.

## 📊 Estado del Proyecto

- **Madurez técnica:** 94% (SENIOR)
- **Fases completadas:** 10/10 (100%)
- **PRs mergeados:** 34
- **Documentación:** 30+ documentos técnicos
- **Tests:** 51 test cases, 887+ líneas

## 📄 Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 👤 Autor

**Fabián Bele**

- GitHub: [@fabianbele2605](https://github.com/fabianbele2605)
- LinkedIn: [Fabián Bele](https://linkedin.com/in/fabianbele)

## 🙏 Agradecimientos

- Clean Architecture por Robert C. Martin
- Go community por las excelentes librerías
- OWASP por las recomendaciones de seguridad
- Prometheus por el sistema de métricas

---

⭐ Si este proyecto te resulta útil, considera darle una estrella en GitHub!
