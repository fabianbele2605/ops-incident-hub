# Fase 3 - Main.go y Punto de Entrada

## 1. Objetivo

Crear el punto de entrada de la aplicación con inyección de dependencias completa, inicialización ordenada y manejo de ciclo de vida.

## 2. Ubicación

```
backend/cmd/api/main.go
```

## 3. Responsabilidades del Main

### 3.1 Inicialización

1. ✅ Cargar configuración
2. ✅ Conectar a base de datos
3. ✅ Ejecutar migraciones
4. ✅ Inicializar repositorios
5. ✅ Inicializar casos de uso
6. ✅ Inicializar handlers
7. ✅ Configurar router
8. ✅ Iniciar servidor HTTP

### 3.2 Ciclo de Vida

1. ✅ Startup logging
2. ✅ Graceful shutdown
3. ✅ Manejo de señales (SIGINT, SIGTERM)
4. ✅ Cleanup de recursos

## 4. Flujo de Ejecución

```
main()
  ├─> Cargar Config
  ├─> Conectar DB
  ├─> Ejecutar Migraciones
  ├─> Inyección de Dependencias
  │    ├─> Repositorios
  │    ├─> Casos de Uso
  │    └─> Handlers
  ├─> Setup Router
  ├─> Iniciar Servidor (goroutine)
  └─> Esperar Señal de Shutdown
       └─> Graceful Shutdown
```

## 5. Inyección de Dependencias

### 5.1 Patrón Constructor

Cada componente recibe sus dependencias en el constructor:

```go
// Repositorios (no tienen dependencias)
incidentRepo := postgres.NewIncidentRepository(db)
userRepo := postgres.NewUserRepository(db)

// Casos de uso (dependen de repositorios)
createIncidentUC := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
assignIncidentUC := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
listIncidentsUC := incident.NewListIncidentUseCase(incidentRepo)

// Handlers (dependen de casos de uso)
incidentHandler := handler.NewIncidentHandler(
    createIncidentUC,
    assignIncidentUC,
    listIncidentsUC,
)
```

### 5.2 Ventajas

- ✅ **Explícito** - Dependencias claras y visibles
- ✅ **Testeable** - Fácil crear mocks
- ✅ **Sin magia** - No usa reflection ni frameworks DI
- ✅ **Type-safe** - Errores en compile-time

## 6. Configuración del Servidor

### 6.1 HTTP Server

```go
server := &http.Server{
    Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
    Handler:      router,
    ReadTimeout:  cfg.Server.ReadTimeout,
    WriteTimeout: cfg.Server.WriteTimeout,
}
```

### 6.2 Timeouts

- **ReadTimeout:** Tiempo máximo para leer request
- **WriteTimeout:** Tiempo máximo para escribir response
- **ShutdownTimeout:** Tiempo máximo para graceful shutdown

## 7. Graceful Shutdown

### 7.1 Manejo de Señales

```go
shutdown := make(chan os.Signal, 1)
signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
```

Captura:
- `SIGINT` (Ctrl+C)
- `SIGTERM` (Docker stop, Kubernetes)

### 7.2 Proceso de Shutdown

```go
ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    // Forzar cierre si timeout
    server.Close()
}
```

**Pasos:**
1. Dejar de aceptar nuevas conexiones
2. Esperar que requests activos terminen
3. Si timeout, forzar cierre
4. Cerrar conexión a DB

## 8. Logging

### 8.1 Startup Logs

```
Starting Ops Incident Hub API
Environment: development
Server: 0.0.0.0:8080
Database connected successfully
Applied migration: 001_create_tables
Migrations applied successfully
Server listening on 0.0.0.0:8080
```

### 8.2 Shutdown Logs

```
Received signal: interrupt. Starting graceful shutdown...
Server stopped gracefully
```

## 9. Manejo de Errores

### 9.1 Errores Fatales

Errores que impiden el inicio de la aplicación:

```go
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

- Configuración inválida
- No puede conectar a DB
- Migraciones fallan

### 9.2 Errores de Runtime

Errores durante la ejecución:

```go
case err := <-serverErrors:
    log.Fatalf("Server error: %v", err)
```

- Puerto ya en uso
- Error al iniciar servidor

## 10. Compilación

### 10.1 Build Local

```bash
cd backend
go build -o bin/api ./cmd/api
```

### 10.2 Build en Docker

```dockerfile
RUN go build -o bin/api ./cmd/api
```

### 10.3 Ejecución

```bash
# Local
./bin/api

# Docker
docker-compose up
```

## 11. Variables de Entorno

El main.go lee configuración de variables de entorno:

```bash
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
# ... etc
```

Ver `.env.example` para lista completa.

## 12. Decisiones Técnicas

### 12.1 Inyección Manual vs Framework

**Decisión:** Inyección manual de dependencias  
**Justificación:**
- Más simple y explícito
- Sin dependencias externas
- Fácil de entender y debuggear
- Suficiente para el tamaño del proyecto

**Alternativas:** `wire`, `dig`, `fx` (para proyectos más grandes)

### 12.2 Graceful Shutdown

**Decisión:** Implementar graceful shutdown  
**Justificación:**
- Evita pérdida de requests en progreso
- Importante para producción
- Requerido por Kubernetes
- Buena práctica general

### 12.3 Logging Simple

**Decisión:** Usar `log` estándar por ahora  
**Justificación:**
- Suficiente para desarrollo
- Sin dependencias externas
- Fácil de reemplazar después

**Mejora futura (Fase 6):** Usar `zerolog` o `zap` para logging estructurado

## 13. Validaciones Realizadas

```bash
✅ go build ./cmd/api/...
✅ Aplicación inicia correctamente
✅ Configuración cargada
✅ DB conectada
✅ Migraciones aplicadas
✅ Servidor escuchando
✅ Graceful shutdown funciona
```

## 14. Testing del Main

### 14.1 Prueba Manual

```bash
# Terminal 1: Iniciar aplicación
./bin/api

# Terminal 2: Probar endpoints
curl http://localhost:8080/health

# Terminal 1: Detener con Ctrl+C
# Verificar: "Server stopped gracefully"
```

### 14.2 Prueba con Docker

```bash
docker-compose up
# Verificar logs
docker-compose down
# Verificar graceful shutdown
```

## 15. Próximos Pasos

Mejoras futuras para el main.go:

- [ ] Logging estructurado (Fase 6)
- [ ] Métricas de Prometheus (Fase 6)
- [ ] Health checks avanzados (Fase 6)
- [ ] Tracing distribuido (Fase 6)
- [ ] Feature flags (Fase 9)

---

**Main.go:** ✅ **IMPLEMENTADO**
