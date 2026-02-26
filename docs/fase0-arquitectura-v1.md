# Arquitectura v1 - Ops Incident Hub

## 1. Visión General

Sistema de gestión de incidentes operativos con arquitectura cloud-native en Azure, diseñado para alta disponibilidad, observabilidad y seguridad.

## 2. Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│                         USUARIOS                                 │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Azure Front Door (CDN)                        │
└────────────┬───────────────────────────────────┬────────────────┘
             │                                   │
             │                                   │
             ▼                                   ▼
┌──────────────────────────┐      ┌──────────────────────────────┐
│  Azure Static Web Apps   │      │   Azure Container Apps       │
│  (Frontend - TypeScript) │      │   (Backend API - Go)         │
│                          │      │                              │
│  - React SPA             │      │  - REST API                  │
│  - CDN global            │      │  - Clean Architecture        │
│  - SSL automático        │      │  - Auto-scaling              │
└──────────────────────────┘      └───────────┬──────────────────┘
                                              │
                                              │
                    ┌─────────────────────────┼─────────────────────┐
                    │                         │                     │
                    ▼                         ▼                     ▼
        ┌────────────────────┐   ┌────────────────────┐  ┌──────────────────┐
        │ Azure Database for │   │  Azure Key Vault   │  │ Azure Monitor +  │
        │    PostgreSQL      │   │                    │  │  App Insights    │
        │                    │   │  - Secrets         │  │                  │
        │  - Incidents       │   │  - Certificates    │  │  - Logs          │
        │  - Users           │   │  - Connection str  │  │  - Metrics       │
        │  - Audit Log       │   │                    │  │  - Traces        │
        └────────────────────┘   └────────────────────┘  │  - Alerts        │
                                                          └──────────────────┘
                    │
                    │ (Private Endpoint)
                    │
        ┌────────────────────────────────────┐
        │      Azure Virtual Network         │
        │  - Subnets                         │
        │  - NSG (Network Security Groups)   │
        │  - Private DNS                     │
        └────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Azure Container Registry                      │
│  - Imágenes Docker del backend                                  │
│  - Escaneo de vulnerabilidades                                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                       Azure AD B2C                               │
│  - Autenticación de usuarios                                    │
│  - OAuth 2.0 / OIDC                                             │
└─────────────────────────────────────────────────────────────────┘
```

## 3. Componentes Principales

### 3.1 Frontend (TypeScript + React)
**Responsabilidades:**
- Interfaz de usuario para gestión de incidentes
- Dashboard con métricas en tiempo real
- Autenticación con Azure AD B2C
- Comunicación con API REST

**Tecnologías:**
- React 18+
- TypeScript 5+
- React Query (gestión de estado servidor)
- Tailwind CSS (estilos)

### 3.2 Backend API (Go)
**Responsabilidades:**
- Lógica de negocio de incidentes
- Validación y autorización
- Persistencia de datos
- Auditoría de cambios
- Exposición de métricas

**Estructura (Clean Architecture):**
```
backend/
├── cmd/api/              # Entry point
├── internal/
│   ├── domain/           # Entidades y lógica de negocio
│   │   ├── incident.go
│   │   └── user.go
│   ├── usecase/          # Casos de uso
│   │   ├── create_incident.go
│   │   └── assign_incident.go
│   ├── infrastructure/   # Implementaciones concretas
│   │   ├── postgres/     # Repositorios
│   │   └── auth/         # Azure AD
│   └── api/              # HTTP handlers
│       └── handlers/
└── pkg/                  # Utilidades compartidas
```

### 3.3 Base de Datos (PostgreSQL)
**Modelo de datos inicial:**

```sql
-- Tabla principal de incidentes
incidents (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  severity VARCHAR(20), -- critical, high, medium, low
  status VARCHAR(20),   -- open, assigned, in_progress, resolved, closed
  assigned_to UUID,
  created_by UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  resolved_at TIMESTAMP,
  metadata JSONB        -- Campos flexibles
)

-- Usuarios
users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255),
  role VARCHAR(50),     -- admin, operator, viewer
  created_at TIMESTAMP NOT NULL
)

-- Auditoría (inmutable)
audit_log (
  id BIGSERIAL PRIMARY KEY,
  entity_type VARCHAR(50),  -- incident, user
  entity_id UUID,
  action VARCHAR(50),       -- created, updated, deleted
  changed_by UUID,
  changes JSONB,            -- Diff de cambios
  timestamp TIMESTAMP NOT NULL
)

-- Comentarios en incidentes
incident_comments (
  id UUID PRIMARY KEY,
  incident_id UUID REFERENCES incidents(id),
  user_id UUID REFERENCES users(id),
  comment TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
)
```

## 4. Flujo de Datos

### 4.1 Crear Incidente
```
Usuario → Frontend → Azure AD (auth) → Backend API → PostgreSQL
                                      ↓
                                  Audit Log
                                      ↓
                                Azure Monitor
```

### 4.2 Consultar Dashboard
```
Usuario → Frontend → Backend API → PostgreSQL (agregaciones)
                                 ↓
                            Cache (futuro)
```

## 5. Endpoints API Principales

```
POST   /api/v1/incidents              # Crear incidente
GET    /api/v1/incidents              # Listar incidentes (paginado)
GET    /api/v1/incidents/:id          # Obtener incidente
PATCH  /api/v1/incidents/:id          # Actualizar incidente
POST   /api/v1/incidents/:id/assign   # Asignar incidente
POST   /api/v1/incidents/:id/resolve  # Resolver incidente
POST   /api/v1/incidents/:id/comments # Agregar comentario

GET    /api/v1/metrics                # Métricas del dashboard
GET    /api/v1/health                 # Health check
GET    /api/v1/ready                  # Readiness check
```

## 6. Seguridad

### 6.1 Autenticación
- Azure AD B2C con OAuth 2.0
- JWT tokens con expiración de 1 hora
- Refresh tokens para renovación

### 6.2 Autorización
- RBAC basado en roles:
  - **Admin**: CRUD completo
  - **Operator**: Crear, actualizar, asignar
  - **Viewer**: Solo lectura

### 6.3 Red
- Private endpoints para PostgreSQL
- NSG restrictivos
- HTTPS obligatorio
- CORS configurado

## 7. Observabilidad

### 7.1 Logs
- Formato JSON estructurado
- Niveles: DEBUG, INFO, WARN, ERROR
- Correlación con trace_id

### 7.2 Métricas
- Request rate, latency, errors (RED)
- Incidentes por severidad
- Tiempo promedio de resolución
- Usuarios activos

### 7.3 Alertas
- API latency > 500ms
- Error rate > 5%
- Database connections > 80%
- Incidentes críticos sin asignar > 10 min

## 8. Escalabilidad

### 8.1 Backend
- Horizontal scaling en Container Apps
- Min: 1 instancia
- Max: 10 instancias
- Trigger: CPU > 70% o requests > 100/s

### 8.2 Base de Datos
- Vertical scaling según carga
- Read replicas (futuro)
- Connection pooling

## 9. Resiliencia

- Health checks cada 30s
- Readiness checks antes de recibir tráfico
- Graceful shutdown (30s timeout)
- Circuit breaker para dependencias externas
- Retry con backoff exponencial

## 10. Próximos Pasos

- [ ] Validar arquitectura con stakeholders
- [ ] Definir SLOs (Service Level Objectives)
- [ ] Crear prototipos de UI
- [ ] Iniciar Fase 1: Fundación del repositorio
