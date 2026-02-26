# Fase 3 - Sistema de Configuración

## 1. Objetivo

Implementar un sistema de configuración robusto que permita gestionar diferentes entornos (development, staging, production) de forma segura y mantenible.

## 2. Principios de Configuración

- ✅ **12-Factor App** - Configuración en variables de entorno
- ✅ **Seguridad** - Secretos nunca en código
- ✅ **Validación** - Configuración validada al inicio
- ✅ **Defaults** - Valores por defecto sensatos
- ✅ **Documentación** - Cada variable documentada

## 3. Estructura de Archivos

```
backend/
├── internal/
│   └── config/
│       └── config.go          # Sistema de configuración
├── .env.example               # Plantilla de variables
└── .env                       # Variables locales (git ignored)
```

## 4. Archivo: `backend/internal/config/config.go`

**Contenido que debes escribir:**

```go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LogConfig configuración de logging
type LogConfig struct {
	Level  string
	Format string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "ops_incident_hub"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate valida la configuración
func (c *Config) Validate() error {
	// Validar servidor
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	// Validar base de datos
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validar entorno
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[c.Server.Environment] {
		return fmt.Errorf("invalid ENVIRONMENT: %s", c.Server.Environment)
	}

	return nil
}

// IsDevelopment verifica si está en desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction verifica si está en producción
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// DatabaseDSN retorna el DSN de PostgreSQL
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// getEnv obtiene una variable de entorno o retorna el valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv obtiene una variable de entorno como int
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv obtiene una variable de entorno como duration
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
```

**Explicación:**

- **Config struct:** Agrupa toda la configuración
- **Load():** Carga desde variables de entorno con defaults
- **Validate():** Valida que la configuración sea correcta
- **Helper methods:** IsDevelopment(), IsProduction(), DatabaseDSN()
- **Type-safe getters:** getEnv, getIntEnv, getDurationEnv

---

## 5. Archivo: `backend/.env.example`

**Contenido que debes escribir:**

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_SHUTDOWN_TIMEOUT=30s
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ops_incident_hub
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

**Explicación:**

- Plantilla de todas las variables necesarias
- Valores por defecto para desarrollo local
- Documentación implícita de qué configurar

---

## 6. Archivo: `backend/.env` (local)

**Crear tu archivo local:**

```bash
cd backend
cp .env.example .env
```

**Editar `.env` con tus valores locales:**

```bash
# Ejemplo para desarrollo local
SERVER_PORT=8080
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password_local
DB_NAME=ops_incident_hub_dev
DB_SSLMODE=disable

LOG_LEVEL=debug
LOG_FORMAT=text
```

**IMPORTANTE:** Este archivo NO se commitea (ya está en `.gitignore`)

---

## 7. Actualizar `.gitignore`

Verificar que `.env` esté ignorado:

```bash
# En backend/.gitignore o raíz
.env
*.env
!.env.example
```

---

## 8. Validar Compilación

```bash
cd backend
go build ./internal/config/...
```

---

## 9. Características del Sistema de Configuración

✅ **12-Factor App compliant** - Configuración en variables de entorno  
✅ **Type-safe** - Tipos específicos (int, duration, string)  
✅ **Validación temprana** - Falla rápido si hay errores  
✅ **Defaults sensatos** - Funciona out-of-the-box en desarrollo  
✅ **Documentado** - .env.example como documentación  
✅ **Seguro** - Secretos nunca en código  

---

## 10. Uso en la Aplicación

```go
// En main.go (próximo paso)
cfg, err := config.Load()
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Usar configuración
fmt.Printf("Starting server on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
fmt.Printf("Environment: %s\n", cfg.Server.Environment)
fmt.Printf("Database: %s\n", cfg.DatabaseDSN())
```

---

## 11. Próximos Pasos

Una vez creado el sistema de configuración:

1. Crear repositorios PostgreSQL
2. Crear migraciones de base de datos
3. Crear archivo main.go
4. Crear docker-compose.yml

---

## 12. Orden de Creación

1. Crear carpeta: `mkdir -p backend/internal/config`
2. Crear `backend/internal/config/config.go`
3. Crear `backend/.env.example`
4. Copiar a `backend/.env` y personalizar
5. Validar compilación

---

**Crea estos archivos y avísame cuando termines para continuar con el siguiente paso.**
