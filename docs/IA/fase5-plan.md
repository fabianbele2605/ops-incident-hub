# Fase 5 - CI/CD Profesional - Plan de Implementación

## 🎯 Objetivo de la Fase
Automatizar calidad, seguridad y entrega continua con puertas de control, asegurando que ningún cambio llegue a producción sin pasar controles definidos.

## 📋 Entregables Esperados
1. GitHub Actions workflows configurados
2. Tests unitarios y de integración implementados
3. Escaneo de seguridad automático
4. Estrategia de deployment por entornos
5. Versionado semántico automático
6. Documentación completa de la fase

## 🔧 Implementación Paso a Paso

### **Paso 1: Tests Unitarios - Domain Layer**

#### Archivos a crear:

**1.1 `backend/internal/domain/incident_test.go`**

Tests para la entidad Incident:
- Validación de creación de incident
- Validación de campos requeridos
- Validación de estados válidos
- Validación de prioridades válidas
- Validación de transiciones de estado

**Estructura sugerida:**
```go
package domain_test

import (
	"testing"
	"time"
	
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
)

func TestNewIncident(t *testing.T) {
	// Test creación exitosa
	// Test validación de título vacío
	// Test validación de descripción vacía
	// Test validación de prioridad inválida
	// Test validación de estado inválido
}

func TestIncident_Assign(t *testing.T) {
	// Test asignación exitosa
	// Test asignación con ID vacío
}

func TestIncident_UpdateStatus(t *testing.T) {
	// Test actualización de estado
	// Test transiciones válidas
}
```

**1.2 `backend/internal/domain/user_test.go`**

Tests para la entidad User:
- Validación de creación de user
- Validación de email
- Validación de roles

---

### **Paso 2: Tests Unitarios - Use Cases**

#### Archivos a crear:

**2.1 `backend/internal/usecase/incident/create_incident_test.go`**

Tests para CreateIncidentUseCase:
- Creación exitosa de incident
- Error cuando título está vacío
- Error cuando usuario no existe
- Error cuando falla el repositorio

**Estructura sugerida:**
```go
package incident_test

import (
	"context"
	"testing"
	
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/domain"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
)

// Necesitarás crear mocks de los repositorios
type mockIncidentRepository struct {
	createFunc func(ctx context.Context, incident *domain.Incident) error
}

func (m *mockIncidentRepository) Create(ctx context.Context, incident *domain.Incident) error {
	return m.createFunc(ctx, incident)
}

func TestCreateIncidentUseCase_Execute(t *testing.T) {
	// Test creación exitosa
	// Test error de validación
	// Test error de repositorio
}
```

**2.2 `backend/internal/usecase/incident/assign_incident_test.go`**

Tests para AssignIncidentUseCase

**2.3 `backend/internal/usecase/incident/list_incidents_test.go`**

Tests para ListIncidentUseCase

---

### **Paso 3: Tests de Integración - Repositorios**

#### Archivos a crear:

**3.1 `backend/internal/infrastructure/postgres/incident_repository_test.go`**

Tests de integración con base de datos real:
- Setup de base de datos de test
- Tests de CRUD completo
- Tests de filtros y paginación
- Cleanup después de tests

**Estructura sugerida:**
```go
package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/infrastructure/postgres"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Conectar a base de datos de test
	// Ejecutar migraciones
	// Retornar conexión
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	// Limpiar datos de test
	// Cerrar conexión
}

func TestIncidentRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	
	repo := postgres.NewIncidentRepository(db)
	
	// Test creación
}
```

**3.2 `backend/internal/infrastructure/postgres/user_repository_test.go`**

Tests de integración para UserRepository

---

### **Paso 4: Tests de API - Handlers**

#### Archivos a crear:

**4.1 `backend/internal/api/handler/incident_handler_test.go`**

Tests de handlers HTTP:
- Test de endpoint Create
- Test de endpoint List
- Test de endpoint Assign
- Test de validación de input
- Test de manejo de errores

**Estructura sugerida:**
```go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
)

func TestIncidentHandler_Create(t *testing.T) {
	// Setup handler con mocks
	// Crear request HTTP
	// Ejecutar handler
	// Verificar response
}
```

**4.2 `backend/internal/api/handler/health_handler_test.go`**

Tests para health check handlers

---

### **Paso 5: GitHub Actions - CI Workflow**

#### Archivo a crear:

**5.1 `.github/workflows/ci.yml`**

Workflow de Integración Continua que se ejecuta en cada PR:

```yaml
name: CI

on:
  pull_request:
    branches: [develop, main]
  push:
    branches: [develop, main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          working-directory: backend

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: ops_incident_hub_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run tests
        working-directory: backend
        env:
          DB_HOST: localhost
          DB_PORT: 5432
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: ops_incident_hub_test
          DB_SSLMODE: disable
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./backend/coverage.out

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build
        working-directory: backend
        run: go build -v ./cmd/api

  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker image
        run: docker build -t ops-incident-hub:test ./backend
      
      - name: Test Docker image
        run: |
          docker run -d --name test-container ops-incident-hub:test
          sleep 5
          docker logs test-container
```

---

### **Paso 6: GitHub Actions - Security Workflow**

#### Archivo a crear:

**6.1 `.github/workflows/security.yml`**

Workflow de escaneo de seguridad:

```yaml
name: Security

on:
  pull_request:
    branches: [develop, main]
  push:
    branches: [develop, main]
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday

jobs:
  gosec:
    name: Go Security
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out results.sarif ./backend/...'
      
      - name: Upload SARIF file
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif

  trivy:
    name: Docker Security
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build image
        run: docker build -t ops-incident-hub:${{ github.sha }} ./backend
      
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'ops-incident-hub:${{ github.sha }}'
          format: 'sarif'
          output: 'trivy-results.sarif'
      
      - name: Upload Trivy results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

  dependency-review:
    name: Dependency Review
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
      - uses: actions/checkout@v4
      
      - name: Dependency Review
        uses: actions/dependency-review-action@v4
```

---

### **Paso 7: GitHub Actions - CD Workflow**

#### Archivo a crear:

**7.1 `.github/workflows/cd.yml`**

Workflow de Continuous Deployment:

```yaml
name: CD

on:
  push:
    branches: [main]
    tags:
      - 'v*'

jobs:
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment:
      name: staging
      url: https://staging.ops-incident-hub.com
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build and push Docker image
        run: |
          echo "Building image for staging..."
          docker build -t ops-incident-hub:staging ./backend
          # Push to registry (configurar después)
      
      - name: Deploy to staging
        run: |
          echo "Deploying to staging environment..."
          # Deployment steps (configurar después)

  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/v')
    environment:
      name: production
      url: https://ops-incident-hub.com
    needs: [deploy-staging]
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build and push Docker image
        run: |
          echo "Building image for production..."
          docker build -t ops-incident-hub:${{ github.ref_name }} ./backend
          # Push to registry
      
      - name: Deploy to production
        run: |
          echo "Deploying to production environment..."
          # Deployment steps
```

---

### **Paso 8: Configuración de Linting**

#### Archivos a crear:

**8.1 `backend/.golangci.yml`**

Configuración de golangci-lint:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck
    - gosec
    - gocyclo
    - dupl
    - misspell
    - unparam
    - unconvert
    - goconst
    - goimports

linters-settings:
  gocyclo:
    min-complexity: 15
  dupl:
    threshold: 100
  goconst:
    min-len: 3
    min-occurrences: 3

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

---

### **Paso 9: Makefile para Tests**

#### Archivo a crear/actualizar:

**9.1 `Makefile`**

Agregar comandos de testing:

```makefile
.PHONY: test test-unit test-integration test-coverage lint

test: ## Ejecutar todos los tests
	cd backend && go test -v -race ./...

test-unit: ## Ejecutar solo tests unitarios
	cd backend && go test -v -race -short ./...

test-integration: ## Ejecutar tests de integración
	cd backend && go test -v -race -run Integration ./...

test-coverage: ## Generar reporte de cobertura
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

lint: ## Ejecutar linter
	cd backend && golangci-lint run ./...

ci: lint test ## Ejecutar CI localmente
	@echo "CI checks passed!"
```

---

### **Paso 10: Configuración de Dependencias**

#### Archivo a actualizar:

**10.1 `backend/go.mod`**

Agregar dependencias de testing:

```bash
go get -u github.com/stretchr/testify
go mod tidy
```

---

## 📊 Orden de Implementación Recomendado

1. ✅ **Paso 10**: Agregar dependencias de testing
2. ✅ **Paso 1**: Tests unitarios de domain
3. ✅ **Paso 2**: Tests unitarios de use cases (con mocks)
4. ✅ **Paso 3**: Tests de integración de repositorios
5. ✅ **Paso 4**: Tests de handlers HTTP
6. ✅ **Paso 9**: Actualizar Makefile con comandos de test
7. ✅ **Paso 8**: Configurar golangci-lint
8. ✅ **Paso 5**: GitHub Actions - CI workflow
9. ✅ **Paso 6**: GitHub Actions - Security workflow
10. ✅ **Paso 7**: GitHub Actions - CD workflow (básico)

---

## ✅ Criterios de Validación

Antes de considerar el paso completado, verificar:

### Tests
- [ ] Tests unitarios de domain pasan
- [ ] Tests unitarios de use cases pasan
- [ ] Tests de integración pasan
- [ ] Tests de handlers pasan
- [ ] Cobertura de código > 70%
- [ ] Todos los tests pasan con `-race` flag

### CI/CD
- [ ] Workflow de CI se ejecuta en PRs
- [ ] Workflow de security escanea vulnerabilidades
- [ ] Linter pasa sin errores
- [ ] Docker build exitoso en CI
- [ ] Branch protection rules configuradas

### Calidad
- [ ] No hay errores de linting
- [ ] No hay vulnerabilidades críticas
- [ ] Código formateado correctamente
- [ ] Tests son mantenibles y claros

---

## 🎯 Resultado Esperado

Al finalizar este paso, tendrás:
- ✅ Suite completa de tests (unitarios + integración + API)
- ✅ CI/CD automatizado con GitHub Actions
- ✅ Escaneo de seguridad automático
- ✅ Cobertura de código > 70%
- ✅ Linting automático
- ✅ Branch protection con checks obligatorios
- ✅ Base para deployment automatizado

---

## 📝 Notas Importantes

### Sobre Tests
- Usar `testing` package estándar de Go
- Usar `testify/assert` para assertions más claras
- Crear mocks manualmente (simple y claro)
- Tests de integración usan base de datos real
- Cleanup después de cada test

### Sobre CI/CD
- CI se ejecuta en cada PR
- Security scan semanal + en cada PR
- CD solo en main y tags
- Usar GitHub Environments para staging/prod
- Secrets en GitHub Secrets (no en código)

### Sobre Cobertura
- Objetivo: > 70% de cobertura
- Priorizar código crítico (use cases, domain)
- No obsesionarse con 100%
- Cobertura es métrica, no objetivo

---

## 📚 Próximos Pasos

Una vez completada la implementación:
1. Informarme para crear la documentación de la fase
2. Crear PR con los cambios
3. Verificar que CI pasa
4. Mergear a develop
5. Continuar con deployment a entornos reales

---

**Fecha de creación:** 26 de febrero de 2025  
**Fase:** 5 - CI/CD Profesional  
**Estado:** Planificado
