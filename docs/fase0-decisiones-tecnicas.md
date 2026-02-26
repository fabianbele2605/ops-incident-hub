# Decisiones Técnicas Iniciales - Ops Incident Hub

## 1. Stack Tecnológico

| Decisión | Tecnología Elegida | Justificación | Alternativas Consideradas |
|----------|-------------------|---------------|---------------------------|
| Backend | Go 1.21+ | - Rendimiento superior<br>- Concurrencia nativa (goroutines)<br>- Binarios estáticos (fácil despliegue)<br>- Tipado fuerte<br>- Excelente para APIs REST | Node.js (menos rendimiento), Java (más pesado) |
| Frontend | TypeScript + React | - Tipado estático reduce errores<br>- Ecosistema maduro<br>- Componentes reutilizables<br>- Experiencia del equipo | Vue.js, Angular |
| Base de Datos | PostgreSQL 15+ | - ACID completo<br>- Soporte JSON/JSONB (flexibilidad)<br>- Madurez y estabilidad<br>- Excelente para auditoría | MySQL (menos features), MongoDB (no relacional) |
| IaC | Terraform | - Multi-cloud (portabilidad)<br>- Declarativo<br>- Estado compartido<br>- Módulos reutilizables | Bicep (solo Azure), Pulumi |
| CI/CD | GitHub Actions | - Integración nativa con repo<br>- Gratuito para proyectos públicos<br>- Sintaxis simple YAML | Azure DevOps, GitLab CI |

## 2. Servicios de Azure

| Servicio | Propósito | Justificación |
|----------|-----------|---------------|
| Azure Container Apps | Hosting del backend | - Serverless containers<br>- Escalado automático<br>- Menor complejidad que AKS<br>- Ideal para APIs stateless |
| Azure Database for PostgreSQL | Base de datos principal | - Managed service<br>- Backups automáticos<br>- Alta disponibilidad<br>- Parches automáticos |
| Azure Static Web Apps | Hosting del frontend | - CDN global incluido<br>- CI/CD integrado<br>- SSL automático<br>- Bajo costo |
| Azure Key Vault | Gestión de secretos | - Rotación automática<br>- Auditoría completa<br>- Integración con identidades administradas |
| Azure Monitor + App Insights | Observabilidad | - Logs centralizados<br>- Métricas y trazas<br>- Alertas configurables |
| Azure Container Registry | Registro de imágenes | - Privado y seguro<br>- Escaneo de vulnerabilidades<br>- Integración con Container Apps |

## 3. Arquitectura de Aplicación

| Decisión | Elección | Justificación |
|----------|----------|---------------|
| Patrón arquitectónico | Clean Architecture (Hexagonal) | - Separación de responsabilidades<br>- Testeable<br>- Independiente de frameworks<br>- Facilita cambios futuros |
| Estructura del backend | Capas: domain, usecase, infrastructure, api | - Domain: entidades y lógica de negocio<br>- Usecase: casos de uso<br>- Infrastructure: DB, externos<br>- API: handlers HTTP |
| API Design | RESTful + OpenAPI 3.0 | - Estándar de industria<br>- Documentación automática<br>- Fácil integración |
| Autenticación | Azure AD B2C + JWT | - Identidades administradas<br>- OAuth 2.0 / OIDC<br>- Sin gestión de passwords |

## 4. Estrategia de Entornos

| Entorno | Propósito | Características |
|---------|-----------|-----------------|
| dev | Desarrollo local y pruebas rápidas | - Recursos mínimos<br>- Sin HA<br>- Datos sintéticos |
| staging | Pre-producción y validación | - Réplica de prod<br>- Datos anonimizados<br>- Testing de integración |
| prod | Producción | - Alta disponibilidad<br>- Backups automáticos<br>- Monitoreo 24/7 |

## 5. Estrategia de Datos

| Aspecto | Decisión | Justificación |
|---------|----------|---------------|
| Migraciones | golang-migrate | - Versionado de esquema<br>- Rollback seguro<br>- Integrable en CI/CD |
| Auditoría | Tabla de audit_log + triggers | - Trazabilidad completa<br>- Inmutable<br>- Consultas eficientes |
| Backups | Automáticos diarios + retención 30 días | - Recuperación ante desastres<br>- Cumplimiento normativo |

## 6. Seguridad

| Aspecto | Decisión | Justificación |
|---------|----------|---------------|
| Secretos | Azure Key Vault + Managed Identity | - Sin credenciales en código<br>- Rotación automática |
| Red | Private endpoints + NSG | - Tráfico privado<br>- Superficie de ataque reducida |
| HTTPS | Obligatorio en todos los entornos | - Cifrado en tránsito<br>- Estándar de industria |
| RBAC | Azure AD roles + permisos granulares | - Principio de mínimo privilegio |

## 7. Observabilidad

| Componente | Herramienta | Justificación |
|------------|-------------|---------------|
| Logs | Structured logging (JSON) + Azure Monitor | - Consultas eficientes<br>- Correlación de eventos |
| Métricas | Prometheus format + App Insights | - Estándar de industria<br>- Dashboards automáticos |
| Trazas | OpenTelemetry | - Trazabilidad distribuida<br>- Vendor-neutral |
| Alertas | Azure Monitor Alerts | - Notificaciones proactivas<br>- Integración con runbooks |

## 8. Versionado y Branching

| Aspecto | Estrategia | Justificación |
|---------|-----------|---------------|
| Branching | Git Flow simplificado | - main: producción<br>- develop: integración<br>- feature/*: desarrollo |
| Versionado | Semantic Versioning (SemVer) | - MAJOR.MINOR.PATCH<br>- Claro y predecible |
| Commits | Conventional Commits | - Changelog automático<br>- Trazabilidad |

## 9. Criterios de Calidad

| Aspecto | Umbral Mínimo | Herramienta |
|---------|---------------|-------------|
| Cobertura de tests | 80% | go test -cover |
| Linting | 0 errores críticos | golangci-lint |
| Vulnerabilidades | 0 críticas/altas | Trivy, Snyk |
| Performance | API < 200ms p95 | Load testing |

## 10. Próximos Pasos

- [ ] Validar estas decisiones con el equipo
- [ ] Crear diagrama de arquitectura v1
- [ ] Definir modelo de datos inicial
- [ ] Iniciar Fase 1: Fundación del repositorio
