# Fase 7 - Seguridad Integral: Diseño e Implementación

## 📋 Información General

**Fase:** 7 - Seguridad Integral  
**Objetivo:** Aplicar seguridad de identidad, red, datos y cadena de suministro  
**Principio clave:** "La seguridad no es fase final; se trabaja desde el inicio"

## 🎯 Objetivos de la Fase

1. **Gestión de secretos:** No hardcodear credenciales, usar variables de entorno
2. **Hardening de contenedores:** Non-root user, read-only filesystem, minimal image
3. **Security headers HTTP:** Protección contra ataques comunes (XSS, CSRF, etc.)
4. **Rate limiting:** Protección contra abuse y DDoS
5. **Input validation:** Sanitización robusta de inputs

## 🔒 Alcance Adaptado

**Nota:** Esta fase se adapta para entorno local/Docker. Los componentes de Azure (Key Vault, Managed Identities) se implementarán cuando se despliegue en Azure.

### Implementación Actual (Local/Docker)
- ✅ Gestión de secretos con variables de entorno
- ✅ Hardening de contenedores Docker
- ✅ Security headers HTTP
- ✅ Rate limiting
- ✅ Input validation mejorada

### Implementación Futura (Azure)
- ⏳ Azure Key Vault para secretos
- ⏳ Managed Identities para autenticación
- ⏳ Azure Security Center
- ⏳ Network Security Groups

## 🛡️ Paso 1: Hardening de Contenedores

### Objetivo
Mejorar la seguridad del contenedor Docker siguiendo best practices.

### Mejoras a Implementar

#### 1. Non-Root User

**Problema actual:** Contenedor corre como root (riesgo de seguridad)

**Solución:** Crear usuario no privilegiado

```dockerfile
# Crear usuario no privilegiado
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Cambiar ownership
RUN chown -R appuser:appuser /app

# Cambiar a usuario no privilegiado
USER appuser
```

#### 2. Read-Only Root Filesystem

**Problema actual:** Filesystem escribible (riesgo de modificación)

**Solución:** Montar filesystem como read-only

```yaml
# docker-compose.yml
services:
  api:
    read_only: true
    tmpfs:
      - /tmp
```

#### 3. Security Options

**Mejoras:**
- Drop capabilities innecesarias
- No new privileges
- Seccomp profile

```yaml
# docker-compose.yml
services:
  api:
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

### Archivo a Modificar

**`backend/Dockerfile`:**
- Agregar usuario no privilegiado
- Cambiar ownership de archivos
- USER appuser al final

**`docker-compose.yml`:**
- Agregar security_opt
- Agregar read_only
- Agregar tmpfs para /tmp

## 🔐 Paso 2: Security Headers HTTP

### Objetivo
Agregar headers de seguridad HTTP para proteger contra ataques comunes.

### Headers a Implementar

```go
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### Estructura de Implementación

```
backend/internal/api/middleware/
├── security.go         # Security headers middleware
└── cors.go            # CORS middleware
```

### Qué Implementar

#### 1. `backend/internal/api/middleware/security.go`

**Responsabilidad:** Agregar security headers a todas las responses

**Funciones:**
```go
// SecurityHeadersMiddleware agrega headers de seguridad
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        next.ServeHTTP(w, r)
    })
}
```

#### 2. `backend/internal/api/middleware/cors.go`

**Responsabilidad:** Configurar CORS de forma segura

**Funciones:**
```go
// CORSMiddleware configura CORS
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Validar origin
            if isAllowedOrigin(origin, allowedOrigins) {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                w.Header().Set("Access-Control-Max-Age", "3600")
            }
            
            // Handle preflight
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Integración en Router

```go
// router.go
router.Use(middleware.SecurityHeadersMiddleware)
router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
```

## ⏱️ Paso 3: Rate Limiting

### Objetivo
Proteger la API contra abuse y ataques de fuerza bruta.

### Estrategia

**Librería:** `golang.org/x/time/rate`

**Límites propuestos:**
- Global: 100 requests/segundo
- Por IP: 10 requests/segundo
- Por endpoint crítico: 5 requests/minuto

### Estructura de Implementación

```
backend/internal/api/middleware/
└── ratelimit.go       # Rate limiting middleware
```

### Qué Implementar

#### 1. `backend/internal/api/middleware/ratelimit.go`

**Responsabilidad:** Limitar rate de requests

**Implementación:**
```go
package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "golang.org/x/time/rate"
)

// RateLimiter almacena limiters por IP
type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

// GetLimiter obtiene o crea un limiter para una IP
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.limiters[ip]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[ip] = limiter
    }
    
    return limiter
}

// RateLimitMiddleware middleware de rate limiting
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getIP(r)
        limiter := rl.GetLimiter(ip)
        
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// getIP extrae IP del request
func getIP(r *http.Request) string {
    // Intentar X-Forwarded-For primero
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return forwarded
    }
    
    // Usar RemoteAddr
    return r.RemoteAddr
}
```

### Integración

```go
// main.go
rateLimiter := middleware.NewRateLimiter(10, 20) // 10 req/s, burst 20

// router.go
router.Use(rateLimiter.RateLimitMiddleware)
```

## ✅ Paso 4: Input Validation Mejorada

### Objetivo
Validar y sanitizar todos los inputs para prevenir inyecciones.

### Mejoras a Implementar

#### 1. Validación de UUIDs

**Problema:** UUIDs no validados pueden causar errores

**Solución:**
```go
// En handlers
func (h *IncidentHandler) Assign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    
    // Validar UUID
    id, err := uuid.Parse(idStr)
    if err != nil {
        http.Error(w, "Invalid incident ID format", http.StatusBadRequest)
        return
    }
    
    // ... resto del código
}
```

#### 2. Validación de Strings

**Problema:** Strings sin validar pueden contener caracteres peligrosos

**Solución:**
```go
// Validar longitud
if len(input.Title) < 3 || len(input.Title) > 200 {
    return errors.New("title must be between 3 and 200 characters")
}

// Sanitizar HTML
import "html"
input.Title = html.EscapeString(input.Title)
```

#### 3. Validación de Enums

**Problema:** Valores de enum no validados

**Solución:**
```go
// En domain/incident.go
func isValidSeverity(s Severity) bool {
    switch s {
    case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow:
        return true
    default:
        return false
    }
}
```

### Archivo a Crear

**`backend/internal/api/validation/validator.go`:**
- Funciones de validación reutilizables
- Sanitización de inputs
- Validación de formatos

## 🔧 Paso 5: Configuración Segura

### Objetivo
Asegurar que la configuración no exponga información sensible.

### Mejoras

#### 1. Validar Variables de Entorno

```go
// config/config.go
func Load() (*Config, error) {
    cfg := &Config{}
    
    // Validar que variables críticas existan
    if os.Getenv("DB_PASSWORD") == "" {
        return nil, errors.New("DB_PASSWORD is required")
    }
    
    // ... resto de la carga
    
    return cfg, nil
}
```

#### 2. No Loggear Secretos

```go
// Nunca loggear passwords, tokens, etc.
logger.Info("database connected",
    "host", cfg.Database.Host,
    "port", cfg.Database.Port,
    // NO loggear: "password", cfg.Database.Password
)
```

#### 3. Timeouts Configurados

```go
// Ya implementado en main.go
server := &http.Server{
    ReadTimeout:  cfg.Server.ReadTimeout,
    WriteTimeout: cfg.Server.WriteTimeout,
    IdleTimeout:  cfg.Server.IdleTimeout,
}
```

## 📦 Resumen de Archivos a Crear/Modificar

### Archivos a Crear

```
backend/internal/api/middleware/
├── security.go         # Security headers
├── cors.go            # CORS configuration
└── ratelimit.go       # Rate limiting

backend/internal/api/validation/
└── validator.go       # Input validation helpers
```

### Archivos a Modificar

```
backend/
├── Dockerfile                              # Non-root user, security
├── docker-compose.yml                      # Security options
├── internal/api/router.go                  # Agregar middlewares
├── internal/api/handler/incident_handler.go # Validación mejorada
└── internal/config/config.go               # Validación de config
```

## 🔄 Orden de Implementación Recomendado

### Paso 1: Security Headers (PR #24)
1. Crear `middleware/security.go`
2. Crear `middleware/cors.go`
3. Integrar en router
4. Probar headers con curl

**Validación:** Headers presentes en responses

### Paso 2: Rate Limiting (PR #25)
1. Crear `middleware/ratelimit.go`
2. Integrar en router
3. Probar límites con script

**Validación:** 429 Too Many Requests después de límite

### Paso 3: Hardening de Contenedores (PR #26)
1. Modificar Dockerfile (non-root user)
2. Modificar docker-compose.yml (security options)
3. Rebuild y probar

**Validación:** Contenedor corre como non-root

### Paso 4: Input Validation (PR #27)
1. Crear `validation/validator.go`
2. Mejorar validación en handlers
3. Agregar tests de validación

**Validación:** Inputs inválidos rechazados

## ✅ Criterios de Éxito de la Fase

Al finalizar Fase 7, el proyecto debe tener:

### Security Headers
- [x] X-Content-Type-Options
- [x] X-Frame-Options
- [x] X-XSS-Protection
- [x] Strict-Transport-Security
- [x] Content-Security-Policy
- [x] CORS configurado

### Rate Limiting
- [x] Rate limiting por IP
- [x] Límites configurables
- [x] Response 429 cuando se excede

### Container Security
- [x] Non-root user
- [x] Read-only filesystem
- [x] Security options (no-new-privileges)
- [x] Minimal capabilities

### Input Validation
- [x] UUIDs validados
- [x] Strings sanitizados
- [x] Enums validados
- [x] Longitudes verificadas

### Configuration
- [x] Variables de entorno validadas
- [x] No secretos en logs
- [x] Timeouts configurados

## 🎯 Entregables de la Fase

1. **Código:**
   - Middlewares de seguridad
   - Rate limiting
   - Container hardening
   - Input validation

2. **Documentación:**
   - Este documento de diseño
   - Documentación de implementación
   - Guía de seguridad
   - Documento de resumen

3. **Evidencia:**
   - Headers de seguridad en responses
   - Rate limiting funcionando
   - Contenedor corriendo como non-root
   - Validación de inputs

## 📚 Referencias Técnicas

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [HTTP Security Headers](https://owasp.org/www-project-secure-headers/)

## 🚀 Próximos Pasos

Una vez completada esta fase, estarás listo para:
- **Fase 8:** Resiliencia y continuidad (backups, disaster recovery)
- **Fase 9:** Gobierno y costos
- Despliegue en Azure con Key Vault y Managed Identities

---

**Documento creado:** 27 de febrero de 2025  
**Fase:** 7 - Seguridad Integral  
**Estado:** Diseño completo - Listo para implementación  
**Siguiente acción:** Developer implementa Paso 1 (Security Headers)
