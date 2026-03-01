# Guía de Operación - Ops Incident Hub

## 🚀 Inicio Rápido

### Requisitos Previos
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (solo para desarrollo)
- Git

### Instalación

```bash
# 1. Clonar repositorio
git clone https://github.com/fabianbele2605/ops-incident-hub.git
cd ops-incident-hub

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus valores

# 3. Levantar servicios
docker-compose up -d

# 4. Verificar salud
curl http://localhost:8081/health
```

### Verificación

```bash
# Health check completo
curl http://localhost:8081/health | jq

# Liveness probe
curl http://localhost:8081/health/live

# Readiness probe
curl http://localhost:8081/health/ready

# Métricas
curl http://localhost:8081/metrics
```

## 📊 Monitoreo

### Health Checks

**Endpoint principal:** `/health`
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "dependencies": {
    "database": "up"
  },
  "timestamp": "2025-03-01T10:30:00Z"
}
```

**Liveness probe:** `/health/live`
- 200 OK: Proceso vivo
- Uso: Kubernetes liveness probe

**Readiness probe:** `/health/ready`
- 200 OK: Listo para tráfico
- 503 Service Unavailable: No listo
- Uso: Kubernetes readiness probe

### Métricas Prometheus

**Endpoint:** `/metrics`

**Métricas de negocio:**
- `incidents_total{severity="critical"}` - Total de incidentes por severidad
- `incidents_assigned_total` - Total de asignaciones

**Métricas técnicas:**
- `http_requests_total{method,endpoint,status}` - Requests HTTP
- `http_request_duration_seconds{method,endpoint}` - Latencia

### Logs

**Formato:** JSON (producción), texto (desarrollo)

**Campos clave:**
- `time` - Timestamp
- `level` - Nivel (INFO, WARN, ERROR)
- `msg` - Mensaje
- `request_id` - ID único de request
- `method`, `path`, `status` - Datos HTTP

**Ejemplo:**
```json
{
  "time": "2025-03-01T10:30:00Z",
  "level": "INFO",
  "msg": "http request completed",
  "request_id": "req-abc123",
  "method": "POST",
  "path": "/api/v1/incidents",
  "status": 201,
  "duration_ms": 45
}
```

## 🔧 Troubleshooting

### Problema 1: Base de datos no conecta

**Síntomas:**
- Health check retorna `"database": "down"`
- Logs: "failed to ping database"
- HTTP 503 en `/health/ready`

**Diagnóstico:**
```bash
# Verificar PostgreSQL
docker-compose ps postgres

# Ver logs de PostgreSQL
docker-compose logs postgres

# Ver logs de API
docker-compose logs api
```

**Soluciones:**
1. Verificar que PostgreSQL esté corriendo
2. Verificar credenciales en `.env`
3. Verificar conectividad de red
4. Reiniciar servicios: `docker-compose restart`

---

### Problema 2: Rate limiting activo

**Síntomas:**
- HTTP 429 Too Many Requests
- Logs: "rate limit exceeded"

**Diagnóstico:**
```bash
# Ver logs de rate limiting
docker-compose logs api | grep "rate limit"

# Verificar IP del cliente
curl -v http://localhost:8081/health
```

**Soluciones:**
1. Esperar 1 minuto (cleanup automático)
2. Verificar que no sea un ataque
3. Ajustar límites si es legítimo:
   - Editar `cmd/api/main.go`
   - Cambiar `NewRateLimiter(10, 20)` a valores mayores
   - Rebuild y redeploy

---

### Problema 3: Circuit breaker abierto

**Síntomas:**
- Errores "circuit breaker is open"
- Requests fallan inmediatamente
- Logs: "circuit breaker state: open"

**Diagnóstico:**
```bash
# Verificar estado de DB
docker-compose exec postgres pg_isready

# Ver logs de circuit breaker
docker-compose logs api | grep "circuit"
```

**Soluciones:**
1. Verificar estado de base de datos
2. Esperar 30 segundos (reset timeout)
3. Si persiste, revisar logs de DB
4. Reiniciar servicios si es necesario

---

### Problema 4: Memoria alta

**Síntomas:**
- Contenedor usa mucha memoria
- Performance degradada

**Diagnóstico:**
```bash
# Ver uso de recursos
docker stats

# Ver métricas de Go
curl http://localhost:8081/metrics | grep go_
```

**Soluciones:**
1. Verificar connection pool (no debe crecer indefinidamente)
2. Revisar goroutines activas
3. Ajustar límites de memoria en docker-compose.yml
4. Considerar restart periódico

---

### Problema 5: Logs no aparecen

**Síntomas:**
- No hay logs en stdout
- Logs vacíos

**Diagnóstico:**
```bash
# Ver logs del contenedor
docker-compose logs -f api

# Verificar configuración
docker-compose exec api env | grep LOG
```

**Soluciones:**
1. Verificar `LOG_LEVEL` en `.env`
2. Verificar `LOG_FORMAT` (json/text)
3. Reiniciar contenedor
4. Verificar que no haya buffering

## 🚀 Deployment

### Build Local

```bash
# Build de la aplicación
cd backend
go build -o bin/api cmd/api/main.go

# Run local
./bin/api
```

### Build Docker

```bash
# Build de imagen
docker build -t ops-incident-hub:latest .

# Tag para registry
docker tag ops-incident-hub:latest registry.example.com/ops-incident-hub:1.0.0

# Push a registry
docker push registry.example.com/ops-incident-hub:1.0.0
```

### Run en Producción

```bash
# Run con variables de entorno
docker run -d \
  --name ops-incident-hub \
  -p 8080:8080 \
  -e DB_HOST=postgres.example.com \
  -e DB_PASSWORD=secret \
  -e ENVIRONMENT=production \
  ops-incident-hub:latest

# Verificar salud
curl http://localhost:8080/health
```

### Graceful Shutdown

```bash
# Enviar SIGTERM (graceful)
docker stop ops-incident-hub

# Logs mostrarán:
# - "shutdown signal received"
# - "server stopped gracefully"

# Force stop (si graceful falla)
docker kill ops-incident-hub
```

## 🔄 Mantenimiento

### Backups de Base de Datos

```bash
# Backup manual
docker-compose exec postgres pg_dump -U postgres ops_incident_hub > backup.sql

# Restore
docker-compose exec -T postgres psql -U postgres ops_incident_hub < backup.sql
```

**Frecuencia recomendada:** Diaria

### Updates de Aplicación

```bash
# 1. Pull latest code
git pull origin main

# 2. Rebuild imagen
docker-compose build api

# 3. Rolling update
docker-compose up -d api

# 4. Verificar salud
curl http://localhost:8081/health
```

### Limpieza de Recursos

```bash
# Limpiar contenedores parados
docker-compose down

# Limpiar volúmenes (CUIDADO: borra datos)
docker-compose down -v

# Limpiar imágenes antiguas
docker image prune -a
```

### Monitoring Continuo

**Recomendaciones:**
1. Configurar Prometheus para scraping de `/metrics`
2. Crear dashboards en Grafana
3. Configurar alertas para:
   - Health check failures
   - High error rate (>5%)
   - High latency (>1s p95)
   - Circuit breaker open
4. Revisar logs diariamente

## 📞 Soporte

### Logs Importantes

```bash
# Últimos 100 logs
docker-compose logs --tail=100 api

# Logs en tiempo real
docker-compose logs -f api

# Logs de errores
docker-compose logs api | grep ERROR

# Logs por request ID
docker-compose logs api | grep "req-abc123"
```

### Información de Debug

```bash
# Versión de la aplicación
curl http://localhost:8081/health | jq .version

# Estado de dependencias
curl http://localhost:8081/health | jq .dependencies

# Métricas actuales
curl http://localhost:8081/metrics
```

---

**Documento creado:** 1 de marzo de 2025  
**Versión:** 1.0  
**Próxima actualización:** Cuando se agreguen nuevas features
