# Demo Script - Ops Incident Hub

## Información General

**Duración:** 5-7 minutos  
**Audiencia:** Reclutadores técnicos, hiring managers, equipos de ingeniería  
**Objetivo:** Demostrar capacidades técnicas y valor de negocio del sistema

---

## Preparación (Antes de la Demo)

### Checklist Pre-Demo

- [ ] Sistema levantado con `docker-compose up -d`
- [ ] Health check funcionando: `curl http://localhost:8080/health`
- [ ] Base de datos limpia (opcional: `docker-compose down -v && docker-compose up -d`)
- [ ] Terminal preparada con comandos listos
- [ ] Navegador con tabs abiertos (Prometheus, logs, código)
- [ ] Postman/Insomnia con colección de requests (opcional)

### Comandos Pre-Cargados

```bash
# Terminal 1: Logs en vivo
docker-compose logs -f api

# Terminal 2: Comandos de demo
# (copiar comandos de las secciones siguientes)
```

---

## Script de Demo (5-7 minutos)

### 1. Introducción (30 segundos)

**Qué decir:**

> "Hola, soy Fabián Bele y les voy a mostrar **Ops Incident Hub**, una plataforma profesional para gestión de incidentes operativos que desarrollé siguiendo estándares de nivel senior.
>
> El sistema implementa **Clean Architecture**, observabilidad completa con **Prometheus**, resiliencia con **Circuit Breaker** y **Retry Policies**, y seguridad siguiendo recomendaciones **OWASP**.
>
> Vamos a ver el sistema en acción."

**Qué mostrar:**
- README.md con badges de madurez (94% SENIOR)
- Arquitectura en diagrama

---

### 2. Arquitectura (1 minuto)

**Qué decir:**

> "La arquitectura está basada en **Clean Architecture** con 4 capas bien definidas:
>
> 1. **Domain Layer** - Entidades y reglas de negocio puras
> 2. **Use Case Layer** - Orquestación de lógica de negocio
> 3. **Infrastructure Layer** - Implementaciones concretas (PostgreSQL, Circuit Breaker)
> 4. **API Layer** - Exposición HTTP con middlewares de seguridad
>
> Esta separación nos da **testabilidad**, **mantenibilidad** e **independencia de frameworks**."

**Qué mostrar:**
- Diagrama de arquitectura en README
- Estructura de carpetas en IDE:
  ```
  backend/internal/
  ├── domain/       # Entidades puras
  ├── usecase/      # Lógica de negocio
  ├── infrastructure/ # PostgreSQL, Circuit Breaker
  └── api/          # HTTP handlers, middlewares
  ```

---

### 3. Features en Acción (2 minutos)

#### 3.1 Health Check (15 segundos)

**Comando:**
```bash
curl http://localhost:8080/health | jq
```

**Qué decir:**

> "Primero verificamos que el sistema está saludable. El health check valida la conexión a la base de datos y retorna el estado de todas las dependencias."

**Output esperado:**
```json
{
  "status": "healthy",
  "timestamp": "2025-03-01T10:00:00Z",
  "checks": {
    "database": "healthy"
  }
}
```

#### 3.2 Crear Incidente (30 segundos)

**Comando:**
```bash
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database connection timeout",
    "description": "Production DB not responding in us-east-1",
    "severity": "high"
  }' | jq
```

**Qué decir:**

> "Ahora creamos un incidente de alta severidad. El sistema valida la entrada, genera un UUID único, registra métricas de negocio y persiste en PostgreSQL."

**Output esperado:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Database connection timeout",
  "description": "Production DB not responding in us-east-1",
  "severity": "high",
  "status": "open",
  "created_at": "2025-03-01T10:00:00Z"
}
```

**Qué mostrar en logs:**
```
INFO  incident created  incident_id=550e8400-e29b-41d4-a716-446655440000 severity=high
```

#### 3.3 Listar Incidentes (15 segundos)

**Comando:**
```bash
curl http://localhost:8080/api/v1/incidents | jq
```

**Qué decir:**

> "Listamos todos los incidentes. El sistema retorna un array con todos los incidentes activos."

#### 3.4 Asignar Incidente (30 segundos)

**Comando:**
```bash
# Primero crear un usuario (si no existe)
INCIDENT_ID="550e8400-e29b-41d4-a716-446655440000"
USER_ID="123e4567-e89b-12d3-a456-426614174000"

curl -X POST http://localhost:8080/api/v1/incidents/$INCIDENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\"}" | jq
```

**Qué decir:**

> "Asignamos el incidente a un usuario. El sistema valida que el usuario existe, actualiza el estado a 'in_progress' y registra la métrica de asignación."

**Output esperado:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Database connection timeout",
  "status": "in_progress",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174000"
}
```

---

### 4. Observabilidad (1.5 minutos)

#### 4.1 Logs Estructurados (30 segundos)

**Qué decir:**

> "El sistema genera logs estructurados en JSON con **Request ID** único para tracing. Cada request tiene un UUID que se propaga por todas las capas."

**Qué mostrar:**
```bash
docker-compose logs api | tail -20
```

**Output esperado:**
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

#### 4.2 Métricas Prometheus (45 segundos)

**Comando:**
```bash
curl http://localhost:8080/metrics | grep incidents
```

**Qué decir:**

> "El sistema expone métricas de negocio y técnicas en formato Prometheus. Aquí vemos:
> - **incidents_total** - Total de incidentes por severidad
> - **incidents_assigned_total** - Total de incidentes asignados
> - **http_requests_total** - Requests HTTP por endpoint y status
> - **http_request_duration_seconds** - Latencia de requests"

**Output esperado:**
```
incidents_total{severity="high"} 1
incidents_assigned_total 1
http_requests_total{method="POST",endpoint="/api/v1/incidents",status="201"} 1
http_request_duration_seconds_bucket{le="0.1"} 1
```

#### 4.3 Health Checks Avanzados (15 segundos)

**Comandos:**
```bash
# Liveness probe
curl http://localhost:8080/health/live

# Readiness probe
curl http://localhost:8080/health/ready
```

**Qué decir:**

> "Tenemos 3 tipos de health checks:
> - **/health** - Estado completo con dependencias
> - **/health/live** - Liveness probe para Kubernetes
> - **/health/ready** - Readiness probe que valida DB"

---

### 5. Resiliencia (1 minuto)

#### 5.1 Circuit Breaker (30 segundos)

**Qué decir:**

> "El sistema implementa **Circuit Breaker** para proteger contra fallos en cascada. Si la base de datos falla 5 veces consecutivas, el circuit breaker se abre por 30 segundos y falla rápido sin intentar conectar."

**Qué mostrar:**
- Código de `internal/resilience/circuitbreaker.go`
- Configuración: `maxFailures: 5, resetTimeout: 30s`

**Demo (opcional):**
```bash
# Detener PostgreSQL para simular fallo
docker-compose stop postgres

# Intentar crear incidente (fallará rápido)
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","description":"Test","severity":"low"}'

# Ver logs: "circuit breaker is open"
docker-compose logs api | grep "circuit breaker"

# Reiniciar PostgreSQL
docker-compose start postgres
```

#### 5.2 Retry Policies (30 segundos)

**Qué decir:**

> "Implementamos **Retry Policies** con backoff exponencial. Si una operación falla temporalmente, el sistema reintenta hasta 3 veces con delays de 100ms, 200ms, 400ms antes de fallar definitivamente."

**Qué mostrar:**
- Código de `internal/resilience/retry.go`
- Configuración: `maxAttempts: 3, initialWait: 100ms, maxWait: 5s, multiplier: 2.0`

---

### 6. CI/CD y Seguridad (1 minuto)

#### 6.1 Pipeline CI/CD (30 segundos)

**Qué decir:**

> "El proyecto tiene un pipeline completo en **GitHub Actions** con 6 checks automáticos:
> 1. **Lint** - golangci-lint
> 2. **Test** - Unit e integration tests
> 3. **Security** - gosec (SAST)
> 4. **Vulnerability** - govulncheck
> 5. **Secrets** - gitleaks
> 6. **Build** - Docker multi-stage build"

**Qué mostrar:**
- Archivo `.github/workflows/ci.yml`
- PRs en GitHub con checks pasando (100% success rate)

#### 6.2 Seguridad (30 segundos)

**Comando:**
```bash
curl -I http://localhost:8080/health
```

**Qué decir:**

> "El sistema implementa seguridad en capas:
> - **Security Headers OWASP** - 7 headers de seguridad
> - **CORS** - Validación de origen con wildcard support
> - **Rate Limiting** - 10 req/s por IP con burst de 20
> - **Input Validation** - Validación en todos los endpoints"

**Output esperado (headers):**
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
```

---

### 7. Cierre (30 segundos)

**Qué decir:**

> "En resumen, **Ops Incident Hub** es un proyecto de nivel senior que demuestra:
>
> ✅ **Arquitectura limpia** - Clean Architecture con 4 capas  
> ✅ **Observabilidad completa** - Logs, métricas, tracing  
> ✅ **Resiliencia robusta** - Circuit breaker, retry, graceful shutdown  
> ✅ **Seguridad** - OWASP headers, CORS, rate limiting  
> ✅ **CI/CD maduro** - Pipeline completo con security scanning  
> ✅ **Documentación ejemplar** - 30+ documentos técnicos  
>
> **Madurez técnica: 94% (SENIOR)**
>
> El código está en GitHub y toda la documentación está disponible. ¿Tienen alguna pregunta?"

**Qué mostrar:**
- README con badges
- Métricas de madurez (94% SENIOR)
- Documentación completa en `docs/`

---

## Preguntas Frecuentes

### P: ¿Por qué Clean Architecture?

**R:** "Clean Architecture nos da separación de responsabilidades, testabilidad independiente de cada capa, y la capacidad de cambiar frameworks sin afectar la lógica de negocio. Por ejemplo, podríamos cambiar de PostgreSQL a MongoDB solo modificando la capa de infrastructure."

### P: ¿Cómo escala el sistema?

**R:** "El sistema está diseñado para escalar horizontalmente. Podemos levantar múltiples instancias de la API detrás de un load balancer. El Circuit Breaker y Connection Pooling protegen la base de datos. Para escalar más, podríamos agregar Redis para caching y RabbitMQ para procesamiento asíncrono."

### P: ¿Qué falta para producción?

**R:** "El sistema está listo para producción, pero hay mejoras identificadas en el backlog:
- Container hardening (non-root user)
- E2E tests con Testcontainers
- Distributed tracing con OpenTelemetry
- Dashboards en Grafana

Todas están priorizadas y estimadas en el roadmap técnico."

### P: ¿Cuánto tiempo tomó desarrollar?

**R:** "El proyecto se desarrolló en 10 fases durante 3 días, siguiendo un workflow profesional con:
- 34 PRs mergeados
- 100+ commits
- 30+ documentos técnicos
- 51 test cases
- 100% success rate en CI/CD"

---

## Tips para la Demo

### Antes de Empezar

1. **Practica el timing** - La demo debe durar 5-7 minutos máximo
2. **Prepara los comandos** - Copia/pega para evitar errores de tipeo
3. **Limpia la base de datos** - Empieza con datos frescos
4. **Verifica que todo funciona** - Corre todos los comandos antes

### Durante la Demo

1. **Habla con confianza** - Conoces el sistema mejor que nadie
2. **Muestra, no leas** - Explica lo que está pasando, no leas el código
3. **Maneja errores con gracia** - Si algo falla, explica por qué y cómo se recupera
4. **Conecta con el negocio** - No solo técnica, también valor de negocio

### Después de la Demo

1. **Comparte el repo** - GitHub link en el chat
2. **Ofrece deep dive** - "Puedo profundizar en cualquier área"
3. **Menciona la documentación** - "Toda la arquitectura está documentada"

---

## Recursos Adicionales

- **Repositorio:** https://github.com/fabianbele2605/ops-incident-hub
- **Documentación:** `docs/` folder
- **Arquitectura:** `docs/arquitectura-consolidada.md`
- **Métricas:** `docs/metricas-madurez.md`
- **Roadmap:** `docs/roadmap-tecnico.md`

---

**Última actualización:** Marzo 2025  
**Versión:** 1.0
