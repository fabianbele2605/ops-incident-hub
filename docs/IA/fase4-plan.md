# Fase 4 - Contenedores y Plataforma de Ejecución - Plan de Implementación

## 🎯 Objetivo de la Fase
Ejecutar la aplicación con estándares de seguridad, recursos y disponibilidad mediante optimización de contenedores y configuración multi-entorno.

## 📋 Entregables Esperados
1. Dockerfile optimizado para producción
2. Health checks avanzados (liveness, readiness, startup)
3. Configuración de recursos y límites
4. Docker Compose para múltiples entornos (dev, staging, prod)
5. Documentación completa de la fase

## 🔧 Implementación Paso a Paso

### **Paso 1: Optimizar Dockerfile para Producción**

#### Archivo: `backend/Dockerfile`

**Objetivo**: Crear imagen optimizada, segura y mínima.

**Características a implementar**:
- Multi-stage build con etapa de compilación y runtime
- Imagen base `scratch` para runtime (imagen mínima)
- Compilación con flags de optimización: `-ldflags='-w -s'`
- Usuario no-root (`nobody:nobody`)
- Copiar solo archivos necesarios (binario + migraciones)
- Labels de metadata
- Health check integrado

**Estructura sugerida**:
```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /build
# Copiar go.mod y go.sum primero (cache layer)
# go mod download
# Copiar código fuente
# Compilar con: CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o api ./cmd/api

# Stage 2: Runtime
FROM scratch
# Copiar ca-certificates, timezone, passwd/group
# Copiar binario y migraciones
# USER nobody:nobody
# EXPOSE 8080
# HEALTHCHECK
# ENTRYPOINT
```

---

### **Paso 2: Crear Health Check Handler**

#### Archivo: `backend/internal/api/handler/health_handler.go`

**Objetivo**: Endpoints de salud para Kubernetes/Docker.

**Endpoints a implementar**:

1. **GET /health** - Health check completo
   - Verifica conexión a base de datos
   - Retorna JSON con status y checks
   - Status 200 si healthy, 503 si unhealthy

2. **GET /health/live** - Liveness probe
   - Verifica que la aplicación está viva
   - Solo retorna 200 OK
   - No verifica dependencias

3. **GET /health/ready** - Readiness probe
   - Verifica que la app está lista para recibir tráfico
   - Verifica conexión a base de datos
   - Status 200 si ready, 503 si not ready

**Estructura del handler**:
```go
type HealthHandler struct {
    db *sql.DB
}

type HealthResponse struct {
    Status    string            `json:"status"`
    Timestamp string            `json:"timestamp"`
    Checks    map[string]string `json:"checks"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request)
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request)
```

---

### **Paso 3: Actualizar Router**

#### Archivo: `backend/internal/api/router.go`

**Cambios necesarios**:
1. Agregar parámetro `healthHandler *handler.HealthHandler` a `SetupRouter`
2. Registrar rutas de health antes del prefijo `/api/v1`:
   - `GET /health` → `healthHandler.Health`
   - `GET /health/live` → `healthHandler.Live`
   - `GET /health/ready` → `healthHandler.Ready`

---

### **Paso 4: Actualizar Main.go**

#### Archivo: `backend/cmd/api/main.go`

**Cambios necesarios**:

1. **Agregar comando health** al inicio de `main()`:
```go
if len(os.Args) > 1 && os.Args[1] == "health" {
    runHealthCheck()
    return
}
```

2. **Inicializar HealthHandler**:
```go
healthHandler := handler.NewHealthHandler(db)
```

3. **Pasar healthHandler al router**:
```go
router := api.SetupRouter(incidentHandler, healthHandler)
```

4. **Agregar función runHealthCheck** al final del archivo:
```go
func runHealthCheck() {
    cfg, err := config.Load()
    if err != nil {
        os.Exit(1)
    }
    
    url := fmt.Sprintf("http://%s:%s/health/live", cfg.Server.Host, cfg.Server.Port)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        os.Exit(1)
    }
    os.Exit(0)
}
```

---

### **Paso 5: Actualizar Docker Compose**

#### Archivo: `docker-compose.yml`

**Cambios en el servicio `api`**:

1. Agregar health check:
```yaml
healthcheck:
  test: ["CMD", "/app/api", "health"]
  interval: 30s
  timeout: 3s
  start_period: 10s
  retries: 3
```

2. Agregar restart policy:
```yaml
restart: unless-stopped
```

3. Agregar resource limits (opcional para dev):
```yaml
deploy:
  resources:
    limits:
      cpus: '1'
      memory: 512M
    reservations:
      cpus: '0.5'
      memory: 256M
```

---

### **Paso 6: Crear Docker Compose para Desarrollo**

#### Archivo: `docker-compose.dev.yml`

**Objetivo**: Configuración optimizada para desarrollo con hot-reload.

**Características**:
- Usar `Dockerfile.dev` con hot-reload (air)
- Montar volúmenes para código fuente
- Variables de entorno para desarrollo
- Logs más verbosos

---

### **Paso 7: Crear Dockerfile para Desarrollo**

#### Archivo: `backend/Dockerfile.dev`

**Objetivo**: Imagen de desarrollo con hot-reload.

**Características**:
- Basado en `golang:1.22-alpine`
- Instalar `air` para hot-reload: `go install github.com/cosmtrek/air@latest`
- Montar código como volumen
- Ejecutar con `air`

---

### **Paso 8: Configuración de Air**

#### Archivo: `backend/.air.toml`

**Objetivo**: Configurar hot-reload para desarrollo.

**Configuración básica**:
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor"]
  delay = 1000
```

---

### **Paso 9: Actualizar .dockerignore**

#### Archivo: `backend/.dockerignore`

**Agregar**:
```
tmp/
.air.toml
Dockerfile*
.dockerignore
```

---

### **Paso 10: Crear Makefile**

#### Archivo: `Makefile` (en raíz del proyecto)

**Objetivo**: Comandos útiles para desarrollo y operación.

**Comandos a incluir**:
- `make build` - Compilar aplicación
- `make run` - Ejecutar localmente
- `make docker-build` - Construir imagen
- `make docker-up` - Levantar servicios
- `make docker-down` - Detener servicios
- `make docker-logs` - Ver logs
- `make health` - Verificar salud
- `make clean` - Limpiar archivos generados

---

## 📊 Orden de Implementación Recomendado

1. ✅ Health Handler (`health_handler.go`)
2. ✅ Actualizar Router (`router.go`)
3. ✅ Actualizar Main (`main.go`)
4. ✅ Compilar y probar: `go build ./cmd/api`
5. ✅ Optimizar Dockerfile (`Dockerfile`)
6. ✅ Actualizar Docker Compose (`docker-compose.yml`)
7. ✅ Crear Dockerfile.dev (`Dockerfile.dev`)
8. ✅ Crear configuración Air (`.air.toml`)
9. ✅ Actualizar .dockerignore
10. ✅ Crear Makefile
11. ✅ Probar con Docker: `docker-compose up --build`

---

## ✅ Criterios de Validación

Antes de considerar el paso completado, verificar:

- [ ] La aplicación compila sin errores
- [ ] Los endpoints de health responden correctamente:
  - `curl http://localhost:8081/health` → JSON con status
  - `curl http://localhost:8081/health/live` → 200 OK
  - `curl http://localhost:8081/health/ready` → 200 OK
- [ ] Docker build exitoso con imagen optimizada
- [ ] Docker Compose levanta servicios correctamente
- [ ] Health checks de Docker funcionan: `docker ps` muestra "healthy"
- [ ] La aplicación se reinicia automáticamente si falla

---

## 🎯 Resultado Esperado

Al finalizar este paso, tendrás:
- ✅ Dockerfile optimizado (~10MB vs ~300MB)
- ✅ Health checks funcionando en 3 niveles
- ✅ Configuración lista para múltiples entornos
- ✅ Herramientas de desarrollo (hot-reload)
- ✅ Comandos automatizados (Makefile)

---

## 📝 Próximos Pasos

Una vez completada la implementación:
1. Informarme para crear la documentación de la fase
2. Crear PR con los cambios
3. Mergear a develop
4. Continuar con configuraciones multi-entorno (staging, prod)
