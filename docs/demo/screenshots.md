# Screenshots - Ops Incident Hub

## Guía de Capturas de Pantalla

Este documento describe las capturas de pantalla necesarias para documentar el proyecto en portafolio, README y presentaciones.

---

## 1. Arquitectura y Código

### 1.1 Estructura del Proyecto

**Descripción:** Vista del árbol de directorios mostrando Clean Architecture

**Ubicación:** IDE (VS Code / IntelliJ)

**Qué capturar:**
```
backend/internal/
├── domain/           # ← Resaltar
├── usecase/          # ← Resaltar
├── infrastructure/   # ← Resaltar
└── api/              # ← Resaltar
```

**Uso:** README, presentación (slide de arquitectura)

**Placeholder:** `[SCREENSHOT: project-structure.png]`

---

### 1.2 Domain Layer - Entidades

**Descripción:** Código de entidad `Incident` mostrando validaciones

**Ubicación:** `backend/internal/domain/incident.go`

**Qué capturar:**
- Struct `Incident` con campos
- Constructor `NewIncident()` con validaciones
- Métodos de negocio

**Uso:** Presentación (slide de Clean Architecture)

**Placeholder:** `[SCREENSHOT: domain-incident.png]`

---

### 1.3 Circuit Breaker Implementation

**Descripción:** Implementación del Circuit Breaker

**Ubicación:** `backend/internal/resilience/circuitbreaker.go`

**Qué capturar:**
- Struct `CircuitBreaker`
- Método `Execute()`
- Estados (Closed, Open, Half-Open)

**Uso:** Presentación (slide de resiliencia)

**Placeholder:** `[SCREENSHOT: circuit-breaker-code.png]`

---

## 2. Sistema en Ejecución

### 2.1 Docker Compose Up

**Descripción:** Terminal mostrando `docker-compose up` exitoso

**Comando:**
```bash
docker-compose up -d
docker-compose ps
```

**Qué capturar:**
- Servicios levantándose (api, postgres)
- Estado "Up" en ambos servicios
- Puertos expuestos (8080, 5432)

**Uso:** README (Quick Start), deployment guide

**Placeholder:** `[SCREENSHOT: docker-compose-up.png]`

---

### 2.2 Health Check Response

**Descripción:** Response del health check mostrando sistema saludable

**Comando:**
```bash
curl http://localhost:8080/health | jq
```

**Qué capturar:**
```json
{
  "status": "healthy",
  "timestamp": "2025-03-01T10:00:00Z",
  "checks": {
    "database": "healthy"
  }
}
```

**Uso:** README, demo, presentación

**Placeholder:** `[SCREENSHOT: health-check.png]`

---

### 2.3 Crear Incidente - Request/Response

**Descripción:** Creación exitosa de incidente con curl

**Comando:**
```bash
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database connection timeout",
    "description": "Production DB not responding",
    "severity": "high"
  }' | jq
```

**Qué capturar:**
- Request completo
- Response con ID generado, status "open", timestamp

**Uso:** README (Features), demo

**Placeholder:** `[SCREENSHOT: create-incident.png]`

---

### 2.4 Listar Incidentes

**Descripción:** Lista de incidentes con múltiples registros

**Comando:**
```bash
curl http://localhost:8080/api/v1/incidents | jq
```

**Qué capturar:**
- Array de incidentes
- Diferentes severidades (low, medium, high, critical)
- Diferentes estados (open, in_progress, resolved)

**Uso:** README, demo

**Placeholder:** `[SCREENSHOT: list-incidents.png]`

---

## 3. Observabilidad

### 3.1 Logs Estructurados

**Descripción:** Logs en JSON mostrando Request ID y contexto

**Comando:**
```bash
docker-compose logs api | tail -20
```

**Qué capturar:**
```json
{
  "time": "2025-03-01T10:00:00Z",
  "level": "INFO",
  "msg": "request completed",
  "request_id": "req-123abc",
  "method": "POST",
  "path": "/api/v1/incidents",
  "status": 201,
  "duration_ms": 45
}
```

**Uso:** Presentación (slide de observabilidad)

**Placeholder:** `[SCREENSHOT: structured-logs.png]`

---

### 3.2 Métricas Prometheus

**Descripción:** Endpoint /metrics mostrando métricas de negocio y técnicas

**Comando:**
```bash
curl http://localhost:8080/metrics | grep -E "(incidents_|http_requests_)"
```

**Qué capturar:**
```
incidents_total{severity="high"} 5
incidents_total{severity="critical"} 2
incidents_assigned_total 7
http_requests_total{method="POST",endpoint="/api/v1/incidents",status="201"} 10
http_request_duration_seconds_bucket{le="0.1"} 95
```

**Uso:** Presentación (slide de observabilidad)

**Placeholder:** `[SCREENSHOT: prometheus-metrics.png]`

---

### 3.3 Health Checks (3 tipos)

**Descripción:** Los 3 tipos de health checks funcionando

**Comandos:**
```bash
curl http://localhost:8080/health | jq
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

**Qué capturar:**
- /health con JSON completo
- /live con 200 OK
- /ready con 200 OK

**Uso:** Presentación, deployment guide

**Placeholder:** `[SCREENSHOT: health-checks-all.png]`

---

## 4. Seguridad

### 4.1 Security Headers

**Descripción:** Headers de seguridad OWASP en response

**Comando:**
```bash
curl -I http://localhost:8080/health
```

**Qué capturar:**
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

**Uso:** Presentación (slide de seguridad)

**Placeholder:** `[SCREENSHOT: security-headers.png]`

---

### 4.2 Rate Limiting

**Descripción:** Rate limiting en acción (429 Too Many Requests)

**Comando:**
```bash
# Hacer 15 requests rápidos
for i in {1..15}; do curl -w "\n" http://localhost:8080/health; done
```

**Qué capturar:**
- Primeros 10 requests: 200 OK
- Requests 11-15: 429 Too Many Requests
- Mensaje: "rate limit exceeded"

**Uso:** Presentación (slide de seguridad)

**Placeholder:** `[SCREENSHOT: rate-limiting.png]`

---

## 5. CI/CD

### 5.1 GitHub Actions Pipeline

**Descripción:** Pipeline completo pasando en GitHub

**Ubicación:** GitHub → Actions tab

**Qué capturar:**
- Workflow run completo
- 6 checks pasando (✓ lint, ✓ test, ✓ security, ✓ build, etc.)
- Tiempo de ejecución
- 100% success rate

**Uso:** Presentación (slide de CI/CD), README

**Placeholder:** `[SCREENSHOT: github-actions-pipeline.png]`

---

### 5.2 Pull Request con Checks

**Descripción:** PR con todos los checks pasando

**Ubicación:** GitHub → Pull Requests

**Qué capturar:**
- PR title y descripción
- 6 checks passed
- Merge button habilitado
- Conventional commit message

**Uso:** Presentación (slide de CI/CD)

**Placeholder:** `[SCREENSHOT: pr-with-checks.png]`

---

### 5.3 Security Scanning Results

**Descripción:** Resultados de security scanning (gosec, govulncheck)

**Ubicación:** GitHub Actions → Security job

**Qué capturar:**
- gosec: 0 issues found
- govulncheck: No vulnerabilities found
- gitleaks: No secrets detected

**Uso:** Presentación (slide de seguridad)

**Placeholder:** `[SCREENSHOT: security-scanning.png]`

---

## 6. Testing

### 6.1 Test Execution

**Descripción:** Tests corriendo exitosamente

**Comando:**
```bash
go test ./... -v
```

**Qué capturar:**
- Tests pasando (PASS)
- Coverage report
- Número de tests (51 test cases)

**Uso:** Presentación (slide de testing)

**Placeholder:** `[SCREENSHOT: tests-passing.png]`

---

### 6.2 Test Coverage

**Descripción:** Coverage report por paquete

**Comando:**
```bash
go test ./... -cover
```

**Qué capturar:**
```
ok      github.com/.../domain       0.123s  coverage: 85.7%
ok      github.com/.../usecase      0.234s  coverage: 82.3%
ok      github.com/.../infrastructure 0.456s coverage: 78.9%
```

**Uso:** Presentación (slide de testing)

**Placeholder:** `[SCREENSHOT: test-coverage.png]`

---

## 7. Documentación

### 7.1 README con Badges

**Descripción:** README principal con badges de estado

**Ubicación:** GitHub → README.md

**Qué capturar:**
- Badges (Build, Go Version, License, Coverage, Maturity)
- Descripción ejecutiva
- Diagrama de arquitectura
- Quick Start

**Uso:** Portafolio, LinkedIn

**Placeholder:** `[SCREENSHOT: readme-badges.png]`

---

### 7.2 Documentación Técnica

**Descripción:** Carpeta docs/ con toda la documentación

**Ubicación:** GitHub → docs/

**Qué capturar:**
- Lista de documentos (30+ archivos)
- Categorías (diseño, resumen, governance)
- Tamaños de archivos

**Uso:** Presentación (slide de documentación)

**Placeholder:** `[SCREENSHOT: docs-folder.png]`

---

### 7.3 Métricas de Madurez

**Descripción:** Documento de métricas mostrando 94% SENIOR

**Ubicación:** `docs/metricas-madurez.md`

**Qué capturar:**
- Tabla de dimensiones con puntuaciones
- Total: 33/35 (94%)
- Nivel: SENIOR
- Fortalezas y áreas de mejora

**Uso:** Presentación (slide de métricas), README

**Placeholder:** `[SCREENSHOT: maturity-metrics.png]`

---

## 8. Resiliencia

### 8.1 Circuit Breaker en Acción

**Descripción:** Circuit breaker abriéndose después de 5 fallos

**Setup:**
```bash
# Detener PostgreSQL
docker-compose stop postgres

# Intentar crear incidente 6 veces
for i in {1..6}; do
  curl -X POST http://localhost:8080/api/v1/incidents \
    -H "Content-Type: application/json" \
    -d '{"title":"Test","description":"Test","severity":"low"}'
  echo ""
done
```

**Qué capturar:**
- Primeros 5 requests: Error de DB
- Request 6: "circuit breaker is open"
- Logs mostrando transición de estado

**Uso:** Presentación (slide de resiliencia)

**Placeholder:** `[SCREENSHOT: circuit-breaker-open.png]`

---

### 8.2 Retry Policy Logs

**Descripción:** Logs mostrando retry con backoff exponencial

**Ubicación:** Logs de la aplicación

**Qué capturar:**
```
INFO  attempt 1 failed, retrying in 100ms
INFO  attempt 2 failed, retrying in 200ms
INFO  attempt 3 failed, giving up
```

**Uso:** Presentación (slide de resiliencia)

**Placeholder:** `[SCREENSHOT: retry-policy-logs.png]`

---

## Resumen de Capturas

**Total de screenshots necesarios:** 20

### Por Categoría:
- **Arquitectura y Código:** 3 screenshots
- **Sistema en Ejecución:** 4 screenshots
- **Observabilidad:** 3 screenshots
- **Seguridad:** 2 screenshots
- **CI/CD:** 3 screenshots
- **Testing:** 2 screenshots
- **Documentación:** 3 screenshots
- **Resiliencia:** 2 screenshots

### Prioridad:
- **Alta (README):** 1.1, 2.1, 2.2, 2.3, 5.1, 7.1
- **Media (Presentación):** 1.2, 1.3, 3.1, 3.2, 4.1, 5.2, 6.1, 7.3, 8.1
- **Baja (Documentación):** 2.4, 3.3, 4.2, 5.3, 6.2, 7.2, 8.2

---

## Instrucciones para Captura

### Herramientas Recomendadas:
- **macOS:** Cmd+Shift+4 (selección), Cmd+Shift+3 (pantalla completa)
- **Linux:** Flameshot, GNOME Screenshot
- **Windows:** Snipping Tool, Win+Shift+S

### Formato:
- **Formato:** PNG (mejor calidad)
- **Resolución:** 1920x1080 o superior
- **Compresión:** Usar TinyPNG para reducir tamaño

### Naming Convention:
```
screenshots/
├── 01-project-structure.png
├── 02-domain-incident.png
├── 03-circuit-breaker-code.png
├── 04-docker-compose-up.png
├── 05-health-check.png
...
```

---

**Última actualización:** Marzo 2025  
**Versión:** 1.0
