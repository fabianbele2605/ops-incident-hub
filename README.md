# Ops Incident Hub

[![CI/CD](https://github.com/fabianbele2605/ops-incident-hub/workflows/CI/badge.svg)](https://github.com/fabianbele2605/ops-incident-hub/actions)
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
git clone https://github.com/fabianbele2605/ops-incident-hub.git
cd ops-incident-hub

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

- **Fabian Bele** - *Trabajo Inicial* - [@fabianbele2605](https://github.com/fabianbele2605)

## 🙏 Agradecimientos

- Proyecto desarrollado como parte de un portafolio profesional senior
- Inspirado en mejores prácticas de la industria
- Construido con estándares de producción

## 📞 Contacto

- GitHub: [@fabianbele2605](https://github.com/fabianbele2605)
- LinkedIn: [Tu Perfil](https://www.linkedin.com/in/fabian-enrique-bele%C3%B1o-robles-696960261/)
- Email: fabianrobles321@outlook.com | fabianbele19@gmail.com

---

**Nota:** Este es un proyecto de portafolio que demuestra capacidades de nivel senior en arquitectura cloud, DevOps y desarrollo full-stack.