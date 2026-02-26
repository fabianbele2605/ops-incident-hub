# Fase 1 - Archivos de Configuración

## 1. .gitignore

### 1.1 Ubicación:
Crear archivo en raíz: `.gitignore`

### 1.2 Contenido que debes escribir:

```gitignore
# ============================================
# Sistema Operativo
# ============================================
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
*~

# ============================================
# IDEs y Editores
# ============================================
.vscode/
.idea/
*.swp
*.swo
*.swn
.project
.classpath
.settings/
*.sublime-project
*.sublime-workspace

# ============================================
# Backend (Go)
# ============================================
# Binarios
*.exe
*.exe~
*.dll
*.so
*.dylib
backend/bin/
backend/dist/

# Test coverage
*.out
coverage.txt
coverage.html

# Go workspace
go.work
go.work.sum

# Vendor (si usas vendoring)
backend/vendor/

# Archivos de configuración local
backend/.env
backend/.env.local
backend/config.local.yaml

# ============================================
# Frontend (Node.js/React)
# ============================================
# Dependencies
frontend/node_modules/
frontend/.pnp
frontend/.pnp.js

# Testing
frontend/coverage/

# Production
frontend/build/
frontend/dist/

# Misc
frontend/.DS_Store
frontend/.env
frontend/.env.local
frontend/.env.development.local
frontend/.env.test.local
frontend/.env.production.local

# Logs
frontend/npm-debug.log*
frontend/yarn-debug.log*
frontend/yarn-error.log*

# ============================================
# Terraform
# ============================================
# Local .terraform directories
**/.terraform/*

# .tfstate files
*.tfstate
*.tfstate.*

# Crash log files
crash.log
crash.*.log

# Exclude all .tfvars files
*.tfvars
*.tfvars.json

# Ignore override files
override.tf
override.tf.json
*_override.tf
*_override.tf.json

# Ignore CLI configuration files
.terraformrc
terraform.rc

# Lock file (comentar si quieres versionarlo)
# .terraform.lock.hcl

# ============================================
# Docker
# ============================================
*.log

# ============================================
# Secrets y Credenciales
# ============================================
*.pem
*.key
*.cert
*.crt
secrets/
.secrets
credentials.json
service-account.json

# ============================================
# Logs
# ============================================
logs/
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
lerna-debug.log*

# ============================================
# Temporales
# ============================================
tmp/
temp/
*.tmp
*.bak
*.cache

# ============================================
# Build artifacts
# ============================================
dist/
build/
out/

# ============================================
# Backups
# ============================================
*.backup
*.old
```

### 1.3 Explicación:

**Propósito:** Evitar que archivos innecesarios o sensibles se versionen en Git.

**Secciones:**
- **Sistema Operativo:** Archivos específicos de OS (macOS, Windows)
- **IDEs:** Configuraciones de editores (VSCode, IntelliJ)
- **Backend:** Binarios de Go, coverage, vendor
- **Frontend:** node_modules, builds, archivos .env
- **Terraform:** Estados, variables, archivos temporales
- **Secrets:** Certificados, claves, credenciales
- **Logs y temporales:** Archivos de log y temporales

**Importante:** Los archivos `.env` NUNCA deben versionarse (contienen secretos).

## 2. .editorconfig

### 2.1 Ubicación:
Crear archivo en raíz: `.editorconfig`

### 2.2 Contenido que debes escribir:

```editorconfig
# EditorConfig: https://EditorConfig.org

# Configuración raíz
root = true

# Configuración para todos los archivos
[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

# Go files
[*.go]
indent_style = tab
indent_size = 4

# Makefiles
[Makefile]
indent_style = tab

# Markdown
[*.md]
trim_trailing_whitespace = false
indent_size = 2

# YAML
[*.{yml,yaml}]
indent_size = 2

# JSON
[*.json]
indent_size = 2

# TypeScript/JavaScript
[*.{ts,tsx,js,jsx}]
indent_size = 2

# Shell scripts
[*.sh]
indent_size = 2

# SQL
[*.sql]
indent_size = 2

# Terraform
[*.tf]
indent_size = 2

# Python (si se usa en scripts)
[*.py]
indent_size = 4
```

### 2.3 Explicación:

**Propósito:** Mantener consistencia de formato entre diferentes editores y desarrolladores.

**Configuraciones clave:**
- **charset = utf-8:** Codificación estándar
- **end_of_line = lf:** Unix line endings (evita problemas entre OS)
- **insert_final_newline = true:** Agrega línea vacía al final (buena práctica)
- **trim_trailing_whitespace = true:** Elimina espacios al final de líneas
- **indent_style:** Espacios para la mayoría, tabs para Go (convención de Go)

**Beneficio:** Todos los desarrolladores tendrán el mismo formato automáticamente.

## 3. README.md Principal

### 3.1 Ubicación:
Crear archivo en raíz: `README.md`

### 3.2 Contenido que debes escribir:

```markdown
# Ops Incident Hub

[![CI/CD](https://github.com/[tu-usuario]/azureSenior/workflows/CI/badge.svg)](https://github.com/[tu-usuario]/azureSenior/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?logo=typescript)](https://www.typescriptlang.org/)

> Plataforma cloud-native para gestión de incidentes operativos con trazabilidad completa, construida con estándares de nivel senior.

## 🎯 Descripción

Ops Incident Hub es una plataforma moderna para gestionar el ciclo de vida completo de incidentes operativos: registro, priorización, asignación, seguimiento y cierre, con auditoría completa y métricas en tiempo real.

### Características Principales

- ✅ **Gestión de Incidentes:** CRUD completo con estados y prioridades
- 👥 **Asignación Inteligente:** Asignación de responsables con SLAs
- 📊 **Dashboard Operativo:** Métricas y KPIs en tiempo real
- 🔍 **Auditoría Completa:** Trazabilidad de todos los cambios
- 🔔 **Alertas:** Notificaciones configurables
- 🔐 **Seguridad:** Autenticación con Azure AD B2C
- 📈 **Observabilidad:** Logs estructurados, métricas y trazas

## 🏗️ Arquitectura

```
Frontend (React + TS) → Backend API (Go) → PostgreSQL
                              ↓
                    Azure Monitor + App Insights
```

**Stack Tecnológico:**
- **Backend:** Go 1.21+ con Clean Architecture
- **Frontend:** TypeScript + React 18+
- **Base de Datos:** PostgreSQL 15+
- **Cloud:** Azure (Container Apps, Static Web Apps, PostgreSQL)
- **IaC:** Terraform
- **CI/CD:** GitHub Actions

📚 [Documentación de Arquitectura](./docs/fase0-arquitectura-v1.md)

## 🚀 Quick Start

### Requisitos Previos

- Go 1.21+
- Node.js 18+
- Docker y Docker Compose
- Azure CLI (para deployment)
- Terraform 1.5+ (para infraestructura)

### Instalación Local

```bash
# 1. Clonar repositorio
git clone https://github.com/[tu-usuario]/azureSenior.git
cd azureSenior

# 2. Configurar variables de entorno
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env.local

# 3. Levantar servicios con Docker Compose
docker-compose up -d

# 4. Ejecutar migraciones
cd backend
make migrate-up

# 5. Iniciar backend
make run

# 6. En otra terminal, iniciar frontend
cd frontend
npm install
npm start
```

La aplicación estará disponible en:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/swagger

## 📖 Documentación

### Para Desarrolladores

- [Guía de Contribución](./CONTRIBUTING.md)
- [Estructura del Proyecto](./docs/fase1-estructura-repositorio.md)
- [Git Workflow](./docs/fase1-git-workflow.md)
- [Decisiones Técnicas](./docs/fase0-decisiones-tecnicas.md)

### Para Operadores

- [Deployment Guide](./docs/deployment.md) *(próximamente)*
- [Runbooks](./docs/runbooks/) *(próximamente)*
- [Monitoreo y Alertas](./docs/monitoring.md) *(próximamente)*

### Arquitectura

- [Arquitectura v1](./docs/fase0-arquitectura-v1.md)
- [ADRs](./docs/adr/) *(próximamente)*

## 🛠️ Comandos Útiles

### Backend

```bash
# Ejecutar tests
make test

# Ejecutar tests con cobertura
make test-coverage

# Ejecutar linter
make lint

# Formatear código
make fmt

# Ejecutar migraciones
make migrate-up

# Rollback migraciones
make migrate-down

# Generar mocks
make mocks
```

### Frontend

```bash
# Instalar dependencias
npm install

# Iniciar en desarrollo
npm start

# Ejecutar tests
npm test

# Build para producción
npm run build

# Ejecutar linter
npm run lint

# Formatear código
npm run format
```

### Infraestructura

```bash
# Inicializar Terraform
cd infrastructure/terraform/environments/dev
terraform init

# Planear cambios
terraform plan

# Aplicar cambios
terraform apply

# Destruir recursos
terraform destroy
```

## 🧪 Testing

### Backend

```bash
# Tests unitarios
go test ./internal/...

# Tests de integración
go test ./tests/integration/...

# Cobertura
go test -cover ./...
```

### Frontend

```bash
# Tests unitarios
npm test

# Tests con cobertura
npm test -- --coverage

# Tests en modo watch
npm test -- --watch
```

## 🚢 Deployment

### Entornos

- **dev:** Desarrollo y pruebas rápidas
- **staging:** Pre-producción, réplica de prod
- **prod:** Producción

### CI/CD

El proyecto usa GitHub Actions para CI/CD:

1. **PR a develop:** Ejecuta tests, linting, security scan
2. **Merge a develop:** Deploy automático a staging
3. **PR a main:** Revisión exhaustiva
4. **Merge a main:** Deploy automático a producción

📚 [Guía de Deployment](./docs/deployment.md) *(próximamente)*

## 📊 Estado del Proyecto

### Fases Completadas

- ✅ **Fase 0:** Definición y Diseño
- ✅ **Fase 1:** Fundación del Repositorio
- 🔄 **Fase 2:** Arquitectura de Aplicación (en progreso)

### Roadmap

- [ ] Fase 2: Arquitectura de Aplicación
- [ ] Fase 3: Infraestructura como Código
- [ ] Fase 4: Contenedores y Plataforma
- [ ] Fase 5: CI/CD Profesional
- [ ] Fase 6: Observabilidad y Operación
- [ ] Fase 7: Seguridad Integral
- [ ] Fase 8: Resiliencia y Continuidad
- [ ] Fase 9: Gobierno y Costos
- [ ] Fase 10: Cierre Profesional

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Por favor lee la [Guía de Contribución](./CONTRIBUTING.md) antes de enviar un PR.

### Proceso

1. Fork el proyecto
2. Crea tu rama de feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver [LICENSE](./LICENSE) para más detalles.

## 👥 Autores

- **[Tu Nombre]** - *Trabajo Inicial* - [@tu-usuario](https://github.com/tu-usuario)

## 🙏 Agradecimientos

- Proyecto desarrollado como parte de un portafolio profesional senior
- Inspirado en mejores prácticas de la industria
- Construido con estándares de producción

## 📞 Contacto

- GitHub: [@tu-usuario](https://github.com/tu-usuario)
- LinkedIn: [Tu Perfil](https://linkedin.com/in/tu-perfil)
- Email: tu-email@example.com

---

**Nota:** Este es un proyecto de portafolio que demuestra capacidades de nivel senior en arquitectura cloud, DevOps y desarrollo full-stack.
```

### 3.3 Explicación:

**Propósito:** Primera impresión del proyecto, punto de entrada para cualquier persona.

**Secciones clave:**
- **Badges:** Estado visual del proyecto (CI, licencia, versiones)
- **Descripción:** Qué es y qué hace
- **Quick Start:** Cómo empezar rápidamente
- **Documentación:** Enlaces a docs detalladas
- **Comandos útiles:** Referencia rápida
- **Contribuir:** Cómo participar
- **Roadmap:** Estado y próximos pasos

**Importante:** Mantener actualizado conforme avanza el proyecto.

## 4. LICENSE

### 4.1 Ubicación:
Crear archivo en raíz: `LICENSE`

### 4.2 Contenido que debes escribir (MIT License):

```
MIT License

Copyright (c) 2024 [Tu Nombre]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### 4.3 Explicación:

**Propósito:** Define los términos legales bajo los cuales se distribuye el código.

**MIT License:**
- Muy permisiva
- Permite uso comercial
- Permite modificación
- Permite distribución
- Solo requiere incluir el copyright notice

**Alternativas:**
- Apache 2.0 (más protección de patentes)
- GPL (requiere que derivados sean open source)
- Unlicense (dominio público)

## 5. docker-compose.yml (para desarrollo local)

### 5.1 Ubicación:
Crear archivo en raíz: `docker-compose.yml`

### 5.2 Contenido que debes escribir:

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: ops-incident-hub-db
    environment:
      POSTGRES_DB: ops_incident_hub
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres_dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis (para caché futuro)
  redis:
    image: redis:7-alpine
    container_name: ops-incident-hub-redis
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

### 5.3 Explicación:

**Propósito:** Levantar servicios de dependencias localmente para desarrollo.

**Servicios:**
- **postgres:** Base de datos principal
- **redis:** Caché (preparado para futuro)

**Características:**
- Usa imágenes Alpine (más ligeras)
- Health checks configurados
- Puertos expuestos para acceso local
- Volúmenes para persistencia de datos

**Uso:**
```bash
# Levantar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

## 6. .env.example (Backend)

### 6.1 Ubicación:
Crear archivo: `backend/.env.example`

### 6.2 Contenido que debes escribir:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ops_incident_hub
DB_USER=postgres
DB_PASSWORD=postgres_dev_password
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Azure AD B2C Configuration
AZURE_AD_TENANT_ID=your-tenant-id
AZURE_AD_CLIENT_ID=your-client-id
AZURE_AD_CLIENT_SECRET=your-client-secret

# JWT Configuration
JWT_SECRET=your-jwt-secret-change-in-production
JWT_EXPIRATION_HOURS=1

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# Observability
OTEL_ENABLED=false
OTEL_ENDPOINT=http://localhost:4318

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 6.3 Explicación:

**Propósito:** Plantilla de variables de entorno para el backend.

**Secciones:**
- **Server:** Configuración del servidor HTTP
- **Database:** Conexión a PostgreSQL
- **Redis:** Conexión a caché
- **Azure AD:** Autenticación
- **JWT:** Tokens de sesión
- **Logging:** Nivel y formato de logs
- **Observability:** OpenTelemetry
- **CORS:** Orígenes permitidos

**Importante:** 
- Este archivo SÍ se versiona (es ejemplo)
- El archivo `.env` real NO se versiona (contiene secretos)
- Cada desarrollador copia `.env.example` a `.env` y ajusta valores

## 7. .env.example (Frontend)

### 7.1 Ubicación:
Crear archivo: `frontend/.env.example`

### 7.2 Contenido que debes escribir:

```env
# API Configuration
REACT_APP_API_URL=http://localhost:8080/api/v1
REACT_APP_API_TIMEOUT=30000

# Azure AD B2C Configuration
REACT_APP_AZURE_AD_CLIENT_ID=your-client-id
REACT_APP_AZURE_AD_AUTHORITY=https://your-tenant.b2clogin.com/your-tenant.onmicrosoft.com/B2C_1_signupsignin
REACT_APP_AZURE_AD_REDIRECT_URI=http://localhost:3000

# Feature Flags
REACT_APP_ENABLE_ANALYTICS=false
REACT_APP_ENABLE_DEBUG=true

# Environment
REACT_APP_ENVIRONMENT=development
```

### 7.3 Explicación:

**Propósito:** Variables de entorno para el frontend.

**Nota importante:** En React, las variables DEBEN empezar con `REACT_APP_` para ser accesibles.

**Secciones:**
- **API:** URL del backend
- **Azure AD:** Configuración de autenticación
- **Feature Flags:** Activar/desactivar funcionalidades
- **Environment:** Identificar entorno

## 8. Resumen de Archivos a Crear

| Archivo | Ubicación | Propósito |
|---------|-----------|-----------|
| .gitignore | Raíz | Archivos a ignorar por Git |
| .editorconfig | Raíz | Configuración de editor |
| README.md | Raíz | Documentación principal |
| LICENSE | Raíz | Licencia del proyecto |
| docker-compose.yml | Raíz | Servicios para desarrollo local |
| .env.example | backend/ | Plantilla de variables backend |
| .env.example | frontend/ | Plantilla de variables frontend |

## 9. Orden de Creación

1. Crear `.gitignore` (primero, para no versionar archivos innecesarios)
2. Crear `.editorconfig` (para formato consistente)
3. Crear `LICENSE`
4. Crear `README.md`
5. Crear `docker-compose.yml`
6. Crear `.env.example` en backend y frontend
7. Hacer commit inicial

## 10. Validación

Después de crear todos los archivos, verifica:

```bash
# Ver estructura
tree -L 2 -a

# Verificar que .gitignore funciona
git status

# Probar docker-compose
docker-compose up -d
docker-compose ps
docker-compose down
```

## 11. Próximos Pasos

1. Crear todos los archivos listados
2. Personalizar con tu información (nombre, usuario GitHub, etc)
3. Hacer commit inicial: `git commit -m "chore: initialize project configuration"`
4. Crear repositorio en GitHub
5. Push inicial: `git push -u origin main`
