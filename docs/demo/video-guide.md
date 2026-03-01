# Video Guide - Ops Incident Hub

## Guía para Grabar Video Demo

Esta guía te ayudará a grabar un video profesional de 5-7 minutos demostrando el proyecto.

---

## 1. Preparación

### 1.1 Herramientas Necesarias

**Grabación de Pantalla:**
- **macOS:** QuickTime Player (gratis), ScreenFlow (pago)
- **Linux:** OBS Studio (gratis), SimpleScreenRecorder
- **Windows:** OBS Studio (gratis), Camtasia (pago)
- **Online:** Loom (gratis para videos cortos)

**Edición (opcional):**
- **Básica:** iMovie (macOS), OpenShot (Linux/Windows)
- **Avanzada:** DaVinci Resolve (gratis), Adobe Premiere

**Audio:**
- Micrófono USB o headset con buena calidad
- Ambiente silencioso
- Prueba de audio antes de grabar

### 1.2 Configuración de Pantalla

**Resolución:**
- **Recomendada:** 1920x1080 (Full HD)
- **Mínima:** 1280x720 (HD)

**Preparación:**
- Cerrar aplicaciones innecesarias
- Desactivar notificaciones (Do Not Disturb)
- Limpiar escritorio
- Aumentar tamaño de fuente en terminal (16-18pt)
- Usar tema de alto contraste

**Terminal:**
```bash
# Aumentar tamaño de fuente
# Usar tema claro o oscuro con buen contraste
# Configurar prompt simple (sin muchos colores)
```

**IDE:**
- Aumentar tamaño de fuente (14-16pt)
- Usar tema con buen contraste
- Ocultar paneles innecesarios

### 1.3 Checklist Pre-Grabación

- [ ] Sistema levantado (`docker-compose up -d`)
- [ ] Health check funcionando
- [ ] Base de datos limpia
- [ ] Comandos preparados en archivo de texto
- [ ] Navegador con tabs necesarios abiertos
- [ ] Micrófono funcionando
- [ ] Notificaciones desactivadas
- [ ] Ambiente silencioso
- [ ] Script de demo revisado

---

## 2. Estructura del Video

### Duración Total: 5-7 minutos

**Segmentos:**
1. **Intro** (30s) - Presentación y contexto
2. **Arquitectura** (1m) - Explicación de capas
3. **Demo en vivo** (2m) - Features funcionando
4. **Observabilidad** (1.5m) - Logs, métricas, health checks
5. **Resiliencia** (1m) - Circuit breaker, retry
6. **CI/CD** (30s) - Pipeline y seguridad
7. **Cierre** (30s) - Resumen y métricas

---

## 3. Script Detallado

### Segmento 1: Intro (30 segundos)

**Qué grabar:**
- Pantalla mostrando README con badges
- Transición a diagrama de arquitectura

**Qué decir:**
> "Hola, soy Fabián Bele y en este video voy a demostrar **Ops Incident Hub**, una plataforma profesional para gestión de incidentes operativos.
>
> El sistema implementa Clean Architecture, observabilidad completa con Prometheus, resiliencia con Circuit Breaker, y seguridad siguiendo recomendaciones OWASP.
>
> El proyecto alcanzó un 94% de madurez técnica, nivel SENIOR."

**Timing:** 0:00 - 0:30

---

### Segmento 2: Arquitectura (1 minuto)

**Qué grabar:**
- Diagrama de arquitectura en README
- Estructura de carpetas en IDE
- Código de domain/incident.go (brevemente)

**Qué decir:**
> "La arquitectura está basada en Clean Architecture con 4 capas:
>
> 1. **Domain** - Entidades y reglas de negocio puras, sin dependencias externas
> 2. **Use Case** - Orquestación de lógica de negocio
> 3. **Infrastructure** - Implementaciones concretas como PostgreSQL y Circuit Breaker
> 4. **API** - Exposición HTTP con middlewares de seguridad
>
> Esta separación nos da testabilidad, mantenibilidad e independencia de frameworks."

**Timing:** 0:30 - 1:30

---

### Segmento 3: Demo en Vivo (2 minutos)

**Qué grabar:**
- Terminal con comandos curl
- Responses en JSON (usar `jq` para formato)
- Logs en otra terminal (split screen)

**Comandos a ejecutar:**

```bash
# 1. Health check (15s)
curl http://localhost:8080/health | jq

# 2. Crear incidente (30s)
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database connection timeout",
    "description": "Production DB not responding",
    "severity": "high"
  }' | jq

# 3. Listar incidentes (15s)
curl http://localhost:8080/api/v1/incidents | jq

# 4. Asignar incidente (30s)
INCIDENT_ID="<id-del-paso-2>"
USER_ID="123e4567-e89b-12d3-a456-426614174000"
curl -X POST http://localhost:8080/api/v1/incidents/$INCIDENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\"}" | jq
```

**Qué decir:**
> "Vamos a ver el sistema en acción.
>
> Primero verificamos que está saludable con el health check. El sistema valida la conexión a PostgreSQL y retorna el estado.
>
> Ahora creamos un incidente de alta severidad. El sistema genera un UUID único, valida la entrada, y persiste en la base de datos.
>
> Podemos listar todos los incidentes activos.
>
> Y finalmente asignamos el incidente a un usuario. El estado cambia automáticamente a 'in_progress'."

**Timing:** 1:30 - 3:30

---

### Segmento 4: Observabilidad (1.5 minutos)

**Qué grabar:**
- Logs estructurados en terminal
- Endpoint /metrics con métricas
- Health checks (3 tipos)

**Comandos a ejecutar:**

```bash
# 1. Logs estructurados (30s)
docker-compose logs api | tail -20

# 2. Métricas Prometheus (45s)
curl http://localhost:8080/metrics | grep incidents

# 3. Health checks (15s)
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

**Qué decir:**
> "El sistema tiene observabilidad completa.
>
> Los logs son estructurados en JSON con Request ID único para tracing. Cada request tiene un UUID que se propaga por todas las capas.
>
> Exponemos métricas de negocio como total de incidentes por severidad, y métricas técnicas como latencia de requests.
>
> Tenemos 3 tipos de health checks: uno completo con dependencias, liveness probe para Kubernetes, y readiness probe que valida la base de datos."

**Timing:** 3:30 - 5:00

---

### Segmento 5: Resiliencia (1 minuto)

**Qué grabar:**
- Código de circuit breaker (brevemente)
- Demo de circuit breaker en acción (opcional)
- Código de retry policy

**Qué decir:**
> "El sistema implementa patrones de resiliencia.
>
> El Circuit Breaker protege contra fallos en cascada. Si la base de datos falla 5 veces consecutivas, el circuit breaker se abre por 30 segundos y falla rápido sin intentar conectar.
>
> Las Retry Policies manejan fallos temporales con backoff exponencial. El sistema reintenta hasta 3 veces con delays crecientes antes de fallar definitivamente.
>
> También tenemos Graceful Shutdown con 30 segundos de timeout y Connection Pooling optimizado."

**Timing:** 5:00 - 6:00

---

### Segmento 6: CI/CD (30 segundos)

**Qué grabar:**
- GitHub Actions pipeline pasando
- PR con checks verdes
- Archivo .github/workflows/ci.yml (brevemente)

**Qué decir:**
> "El proyecto tiene un pipeline completo en GitHub Actions con 6 checks automáticos: lint, tests, security scanning con gosec, vulnerability detection con govulncheck, secrets detection con gitleaks, y build.
>
> Todos los PRs pasan por estos checks antes de mergear. Tenemos 100% success rate en 34 PRs mergeados."

**Timing:** 6:00 - 6:30

---

### Segmento 7: Cierre (30 segundos)

**Qué grabar:**
- README con badges
- Métricas de madurez (94% SENIOR)
- Documentación en carpeta docs/

**Qué decir:**
> "En resumen, Ops Incident Hub demuestra:
>
> - Arquitectura limpia con Clean Architecture
> - Observabilidad completa con logs, métricas y tracing
> - Resiliencia robusta con Circuit Breaker y Retry Policies
> - Seguridad con headers OWASP, CORS y rate limiting
> - CI/CD maduro con security scanning
> - Documentación ejemplar con 30+ documentos técnicos
>
> El proyecto alcanzó 94% de madurez técnica, nivel SENIOR.
>
> Todo el código está en GitHub y la documentación completa está disponible. Gracias por ver el video."

**Timing:** 6:30 - 7:00

---

## 4. Tips de Grabación

### Audio

**Antes de grabar:**
- Prueba el micrófono
- Habla a 15-20cm del micrófono
- Evita ruidos de fondo (ventilador, AC, tráfico)
- Usa un script pero no lo leas textualmente

**Durante la grabación:**
- Habla claro y pausado
- Usa un tono profesional pero amigable
- Haz pausas entre secciones
- Si te equivocas, pausa 3 segundos y repite (fácil de editar)

### Video

**Movimiento de cursor:**
- Mueve el cursor suavemente
- Resalta áreas importantes
- No muevas el cursor nerviosamente

**Transiciones:**
- Usa transiciones suaves entre ventanas
- Pausa 1-2 segundos antes de cambiar de ventana
- Usa Cmd+Tab (macOS) o Alt+Tab (Windows/Linux) para cambiar

**Zoom:**
- Usa zoom para mostrar detalles importantes
- Vuelve a zoom normal después

### Errores Comunes

**Evitar:**
- ❌ Hablar muy rápido
- ❌ Leer el código línea por línea
- ❌ Mostrar errores sin explicar
- ❌ Usar jerga sin explicar
- ❌ Video muy largo (>10 minutos)

**Hacer:**
- ✅ Hablar con confianza
- ✅ Explicar el "por qué", no solo el "qué"
- ✅ Mostrar valor de negocio, no solo técnica
- ✅ Mantener el ritmo dinámico
- ✅ Terminar con call-to-action

---

## 5. Post-Producción

### Edición Básica

**Cortes necesarios:**
- Eliminar pausas largas (>3 segundos)
- Eliminar errores y repeticiones
- Eliminar tiempos de carga largos

**Mejoras opcionales:**
- Agregar intro con título (5 segundos)
- Agregar outro con links (5 segundos)
- Agregar música de fondo suave (opcional)
- Agregar subtítulos (muy recomendado)

### Exportación

**Configuración recomendada:**
- **Formato:** MP4 (H.264)
- **Resolución:** 1920x1080 (Full HD)
- **Frame rate:** 30 fps
- **Bitrate:** 5-8 Mbps
- **Audio:** AAC, 128-192 kbps

**Tamaño esperado:**
- 5 minutos: ~150-250 MB
- 7 minutos: ~200-350 MB

---

## 6. Publicación

### Plataformas

**YouTube:**
- Mejor para videos largos (>5 min)
- SEO friendly
- Embeddable en portafolio

**Loom:**
- Mejor para videos cortos (<10 min)
- Fácil de compartir
- Transcripción automática

**LinkedIn:**
- Mejor para alcance profesional
- Máximo 10 minutos
- Nativo en la plataforma

### Metadata

**Título:**
```
Ops Incident Hub - Sistema de Gestión de Incidentes con Clean Architecture | Demo Técnica
```

**Descripción:**
```
Demo técnica de Ops Incident Hub, una plataforma profesional para gestión de incidentes operativos.

🏗️ Arquitectura: Clean Architecture con 4 capas
📊 Observabilidad: Prometheus + slog + Request ID tracing
🛡️ Resiliencia: Circuit Breaker + Retry Policies
🔒 Seguridad: OWASP headers + CORS + Rate Limiting
🚀 CI/CD: GitHub Actions con security scanning
📈 Madurez: 94% (SENIOR)

🔗 Repositorio: https://github.com/fabianbele2605/ops-incident-hub
📚 Documentación: https://github.com/fabianbele2605/ops-incident-hub/tree/main/docs

Tech Stack: Go, PostgreSQL, Docker, Prometheus, GitHub Actions

#golang #cleanarchitecture #devops #backend #softwareengineering
```

**Tags:**
```
golang, clean architecture, devops, backend, software engineering, 
incident management, observability, prometheus, docker, postgresql,
circuit breaker, microservices, rest api, ci/cd, github actions
```

**Thumbnail:**
- Captura del README con badges
- Título: "Ops Incident Hub"
- Subtítulo: "94% SENIOR"
- Logo de Go + PostgreSQL + Docker

---

## 7. Checklist Final

### Antes de Publicar

- [ ] Video dura 5-7 minutos
- [ ] Audio es claro y sin ruidos
- [ ] Todos los comandos funcionan
- [ ] No hay información sensible (passwords, tokens)
- [ ] Transiciones son suaves
- [ ] Título y descripción son claros
- [ ] Tags son relevantes
- [ ] Thumbnail es atractivo

### Después de Publicar

- [ ] Compartir en LinkedIn
- [ ] Agregar link en README
- [ ] Agregar link en portafolio
- [ ] Compartir en comunidades relevantes (Reddit, Dev.to)

---

## 8. Alternativas al Video

Si no puedes grabar video, considera:

**GIF Animado:**
- Herramienta: LICEcap, Kap (macOS)
- Duración: 30-60 segundos
- Mostrar feature principal

**Slides con Audio:**
- Herramienta: Google Slides + grabación
- Más fácil de editar
- Menos técnico

**Demo Interactivo:**
- Herramienta: Asciinema (terminal recording)
- Reproducible en navegador
- Ideal para demos de CLI

---

## Recursos Adicionales

**Tutoriales de OBS Studio:**
- https://obsproject.com/wiki/
- YouTube: "OBS Studio Tutorial for Beginners"

**Edición de Video:**
- DaVinci Resolve: https://www.blackmagicdesign.com/products/davinciresolve
- OpenShot: https://www.openshot.org/

**Música de Fondo (libre de derechos):**
- YouTube Audio Library
- Incompetech
- Bensound

---

**Última actualización:** Marzo 2025  
**Versión:** 1.0
