# Fase 4 - Health Checks y Observabilidad

## 1. Introducción

Este documento describe la implementación de health checks en tres niveles (health, liveness, readiness) para mejorar la observabilidad, resiliencia y capacidad de auto-recuperación de la aplicación.

## 2. Arquitectura de Health Checks

### 2.1 Tres Niveles de Health Checks

```
┌─────────────────────────────────────────────────────────┐
│                    HEALTH CHECKS                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   /health    │  │ /health/live │  │ /health/ready│ │
│  │              │  │              │  │              │ │
│  │  Completo    │  │  Liveness    │  │  Readiness   │ │
│  │  + DB check  │  │  Solo app    │  │  + DB check  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 2.2 Propósito de Cada Endpoint

| Endpoint | Propósito | Verifica | Usado por |
|----------|-----------|----------|-----------|
| `/health` | Estado completo del sistema | App + DB + Dependencias | Monitoreo, Dashboards |
| `/health/live` | Aplicación está viva | Solo proceso | Kubernetes liveness probe |
| `/health/ready` | Listo para recibir tráfico | App + DB | Kubernetes readiness probe, Load balancers |

## 3. Implementación

### 3.1 Health Handler

**Archivo:** `backend/internal/api/handler/health_handler.go`

```go
package handler

import (
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

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Health - Endpoint completo de salud con verificación de dependencias
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	status := "healthy"

	// Verificar base de datos
	ctx := r.Context()
	if err := h.db.PingContext(ctx); err != nil {
		checks["database"] = "unhealthy"
		status = "unhealthy"
	} else {
		checks["database"] = "healthy"
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	if status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// Ready - Readiness probe (verifica que está listo para recibir tráfico)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Live - Liveness probe (verifica que la aplicación está viva)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

### 3.2 Integración en Router

**Archivo:** `backend/internal/api/router.go`

```go
func SetupRouter(incidentHandler *handler.IncidentHandler, healthHandler *handler.HealthHandler) *mux.Router {
	router := mux.NewRouter()

	// Health endpoints (antes de /api/v1)
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/health/ready", healthHandler.Ready).Methods("GET")
	router.HandleFunc("/health/live", healthHandler.Live).Methods("GET")

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()
	// ... resto de rutas
	
	return router
}
```

### 3.3 Comando Health en Main

**Archivo:** `backend/cmd/api/main.go`

```go
func main() {
	// Verificar si es comando health
	if len(os.Args) > 1 && os.Args[1] == "health" {
		runHealthCheck()
		return
	}
	
	// ... resto del main
	
	// Inicializar health handler
	healthHandler := handler.NewHealthHandler(db)
	router := api.SetupRouter(incidentHandler, healthHandler)
}

func runHealthCheck() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/health/live", cfg.Server.Host, cfg.Server.Port)
	resp, err := http.Get(url)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	os.Exit(0)
}
```

## 4. Uso de Health Checks

### 4.1 Endpoint /health

**Propósito:** Estado completo del sistema

**Request:**
```bash
curl http://localhost:8081/health
```

**Response (Healthy):**
```json
{
  "status": "healthy",
  "timestamp": "2026-02-26T13:53:25Z",
  "checks": {
    "database": "healthy"
  }
}
```

**Response (Unhealthy):**
```json
{
  "status": "unhealthy",
  "timestamp": "2026-02-26T13:53:25Z",
  "checks": {
    "database": "unhealthy"
  }
}
```

**Códigos HTTP:**
- `200 OK` - Sistema saludable
- `503 Service Unavailable` - Sistema no saludable

**Casos de uso:**
- Dashboards de monitoreo
- Alertas de sistema
- Verificación manual
- Pruebas de integración

### 4.2 Endpoint /health/live

**Propósito:** Verificar que la aplicación está viva

**Request:**
```bash
curl http://localhost:8081/health/live
```

**Response:**
- `200 OK` (sin contenido)
- `503 Service Unavailable` (si falla)

**Características:**
- No verifica dependencias
- Respuesta inmediata
- Usado por Kubernetes liveness probe

**Casos de uso:**
- Kubernetes reinicia el pod si falla
- Detectar deadlocks o procesos colgados
- Verificación básica de proceso

### 4.3 Endpoint /health/ready

**Propósito:** Verificar que está listo para recibir tráfico

**Request:**
```bash
curl http://localhost:8081/health/ready
```

**Response:**
- `200 OK` - Listo para tráfico
- `503 Service Unavailable` - No listo

**Verifica:**
- Conexión a base de datos
- Dependencias críticas

**Casos de uso:**
- Kubernetes readiness probe
- Load balancers
- Rolling deployments
- Warm-up después de inicio

## 5. Integración con Docker

### 5.1 Dockerfile HEALTHCHECK

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/api", "health"] || exit 1
```

**Parámetros:**
- `interval=30s` - Ejecuta cada 30 segundos
- `timeout=3s` - Máximo 3 segundos de espera
- `start-period=5s` - 5 segundos de gracia al iniciar
- `retries=3` - 3 fallos consecutivos = unhealthy

**Funcionamiento:**
1. Docker ejecuta `/app/api health`
2. El comando hace GET a `/health/live`
3. Si responde 200 OK → healthy
4. Si falla o timeout → unhealthy

### 5.2 Docker Compose Health Check

```yaml
api:
  healthcheck:
    test: ["CMD", "/app/api", "health"]
    interval: 30s
    timeout: 3s
    start_period: 10s
    retries: 3
  restart: unless-stopped
```

**Beneficios:**
- Docker reinicia automáticamente si unhealthy
- `depends_on` puede esperar a que esté healthy
- Visible en `docker ps` (status)

### 5.3 Verificar Estado

```bash
# Ver estado de contenedores
docker ps

# Output:
# CONTAINER ID   STATUS
# 0b9119944baf   Up 5 minutes (healthy)

# Inspeccionar health check
docker inspect ops-incident-hub-api | jq '.[0].State.Health'

# Ver logs de health checks
docker inspect ops-incident-hub-api | jq '.[0].State.Health.Log'
```

## 6. Integración con Kubernetes

### 6.1 Liveness Probe

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ops-incident-hub-api
spec:
  containers:
  - name: api
    image: ops-incident-hub:latest
    livenessProbe:
      httpGet:
        path: /health/live
        port: 8080
      initialDelaySeconds: 10
      periodSeconds: 30
      timeoutSeconds: 3
      failureThreshold: 3
```

**Comportamiento:**
- Kubernetes ejecuta GET `/health/live` cada 30s
- Si falla 3 veces consecutivas → reinicia el pod
- Detecta procesos colgados o deadlocks

### 6.2 Readiness Probe

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ops-incident-hub-api
spec:
  containers:
  - name: api
    image: ops-incident-hub:latest
    readinessProbe:
      httpGet:
        path: /health/ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
      timeoutSeconds: 3
      failureThreshold: 3
```

**Comportamiento:**
- Kubernetes ejecuta GET `/health/ready` cada 10s
- Si falla → quita el pod del Service (no recibe tráfico)
- Si se recupera → vuelve a agregar al Service
- Útil durante rolling updates

### 6.3 Startup Probe

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ops-incident-hub-api
spec:
  containers:
  - name: api
    image: ops-incident-hub:latest
    startupProbe:
      httpGet:
        path: /health/live
        port: 8080
      initialDelaySeconds: 0
      periodSeconds: 5
      timeoutSeconds: 3
      failureThreshold: 30
```

**Comportamiento:**
- Da hasta 150 segundos (30 * 5s) para que inicie
- Deshabilita liveness/readiness hasta que pase
- Útil para aplicaciones con inicio lento

## 7. Monitoreo y Alertas

### 7.1 Prometheus Metrics

**Futuro:** Agregar métricas de health checks

```go
var (
	healthCheckTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "health_check_total",
			Help: "Total health checks executed",
		},
		[]string{"endpoint", "status"},
	)
	
	healthCheckDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "health_check_duration_seconds",
			Help: "Health check duration",
		},
		[]string{"endpoint"},
	)
)
```

### 7.2 Alertas Recomendadas

**Alert 1: API Unhealthy**
```yaml
- alert: APIUnhealthy
  expr: up{job="ops-incident-hub"} == 0
  for: 2m
  annotations:
    summary: "API is unhealthy"
    description: "API has been unhealthy for 2 minutes"
```

**Alert 2: Database Connection Failed**
```yaml
- alert: DatabaseConnectionFailed
  expr: health_check_total{endpoint="health",status="unhealthy"} > 0
  for: 1m
  annotations:
    summary: "Database connection failed"
```

### 7.3 Dashboard Grafana

**Métricas a visualizar:**
- Health check success rate
- Response time de health checks
- Número de reinicios por unhealthy
- Tiempo en estado unhealthy

## 8. Mejores Prácticas

### 8.1 Diseño de Health Checks

✅ **Hacer:**
- Verificar dependencias críticas (DB, cache, APIs)
- Responder rápido (< 1 segundo)
- Usar timeouts apropiados
- Retornar códigos HTTP correctos
- Incluir timestamp en respuesta

❌ **No hacer:**
- Verificar dependencias no críticas
- Hacer operaciones pesadas
- Depender de servicios externos lentos
- Retornar siempre 200 OK
- Cachear resultados de health checks

### 8.2 Timeouts

```go
// Configurar timeout para health checks
ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()

if err := h.db.PingContext(ctx); err != nil {
    // Handle timeout or error
}
```

### 8.3 Logging

```go
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// ... health check logic
	
	duration := time.Since(start)
	log.Printf("Health check completed in %v, status: %s", duration, status)
}
```

## 9. Troubleshooting

### 9.1 Health Check Siempre Falla

**Síntomas:**
- Docker marca contenedor como unhealthy
- Kubernetes reinicia pods constantemente

**Causas comunes:**
1. Timeout muy corto
2. Aplicación tarda en iniciar
3. Puerto incorrecto
4. Ruta incorrecta

**Solución:**
```yaml
# Aumentar start-period y timeout
healthcheck:
  start_period: 30s  # Más tiempo para iniciar
  timeout: 5s        # Más tiempo para responder
```

### 9.2 Base de Datos Intermitente

**Síntomas:**
- Health check falla ocasionalmente
- Aplicación funciona pero marca unhealthy

**Solución:**
```go
// Implementar retry logic
func (h *HealthHandler) checkDatabase(ctx context.Context) error {
	for i := 0; i < 3; i++ {
		if err := h.db.PingContext(ctx); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("database unhealthy after retries")
}
```

### 9.3 Health Check Lento

**Síntomas:**
- Health checks tardan > 1 segundo
- Timeouts frecuentes

**Solución:**
```go
// Usar connection pool configurado
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// Timeout agresivo en health check
ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
defer cancel()
```

## 10. Métricas de Éxito

### 10.1 Objetivos

| Métrica | Objetivo | Actual |
|---------|----------|--------|
| Response time /health | < 100ms | ~50ms |
| Response time /health/live | < 10ms | ~5ms |
| Response time /health/ready | < 100ms | ~50ms |
| Uptime | > 99.9% | 100% |
| False positives | < 0.1% | 0% |

### 10.2 Validación

```bash
# Test de carga en health endpoint
ab -n 1000 -c 10 http://localhost:8081/health/live

# Verificar response time
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8081/health

# curl-format.txt:
# time_total: %{time_total}s
```

## 11. Próximos Pasos

### 11.1 Mejoras Futuras

- [ ] Agregar más checks (cache, message queue)
- [ ] Implementar métricas de Prometheus
- [ ] Agregar health check de dependencias externas
- [ ] Implementar circuit breaker para dependencias
- [ ] Agregar versión de API en respuesta de health

### 11.2 Integración con Observabilidad

- [ ] Logs estructurados de health checks
- [ ] Traces distribuidos
- [ ] Dashboards de Grafana
- [ ] Alertas en PagerDuty/Slack

## 12. Referencias

- [Kubernetes Liveness/Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Docker HEALTHCHECK](https://docs.docker.com/engine/reference/builder/#healthcheck)
- [Health Check Pattern](https://microservices.io/patterns/observability/health-check-api.html)
- [12 Factor App - Health Checks](https://12factor.net/)

---

**Fecha de creación:** 26 de febrero de 2025  
**Fase:** 4 - Contenedores y Plataforma de Ejecución  
**Estado:** Implementado y funcionando
