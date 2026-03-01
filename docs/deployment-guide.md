# Guía de Deployment - Ops Incident Hub

## Índice

1. [Deployment Local (Docker Compose)](#1-deployment-local-docker-compose)
2. [Deployment AWS (ECS + RDS)](#2-deployment-aws-ecs--rds)
3. [Deployment Azure (Container Apps)](#3-deployment-azure-container-apps)
4. [Variables de Entorno](#4-variables-de-entorno)
5. [Troubleshooting](#5-troubleshooting)
6. [Rollback](#6-rollback)

---

## 1. Deployment Local (Docker Compose)

### 1.1 Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Git

### 1.2 Configuration

**Paso 1:** Clonar el repositorio

```bash
git clone https://github.com/fabianbele2605/ops-incident-hub.git
cd ops-incident-hub
```

**Paso 2:** Crear archivo `.env` (opcional)

```bash
cat > .env << EOF
# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ops_incident_hub

# Security
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Logging
LOG_LEVEL=debug
LOG_FORMAT=text
EOF
```

### 1.3 Start Services

```bash
# Levantar servicios en background
docker-compose up -d

# Ver logs
docker-compose logs -f

# Ver solo logs de la API
docker-compose logs -f api
```

### 1.4 Verify Deployment

```bash
# Health check
curl http://localhost:8080/health

# Respuesta esperada:
# {
#   "status": "healthy",
#   "timestamp": "2025-03-01T10:00:00Z",
#   "checks": {
#     "database": "healthy"
#   }
# }

# Verificar métricas
curl http://localhost:8080/metrics

# Verificar base de datos
docker-compose exec postgres psql -U postgres -d ops_incident_hub -c "\dt"
```

### 1.5 Access Endpoints

- **API:** http://localhost:8080
- **Health:** http://localhost:8080/health
- **Metrics:** http://localhost:8080/metrics
- **PostgreSQL:** localhost:5432

### 1.6 Stop Services

```bash
# Detener servicios
docker-compose down

# Detener y eliminar volúmenes (CUIDADO: elimina datos)
docker-compose down -v
```

---

## 2. Deployment AWS (ECS + RDS)

### 2.1 Prerequisites

- AWS CLI 2.0+
- Terraform 1.5+
- AWS Account con permisos de administrador
- Docker (para build de imagen)

### 2.2 Infrastructure Setup

**Paso 1:** Configurar AWS CLI

```bash
aws configure
# AWS Access Key ID: <tu-access-key>
# AWS Secret Access Key: <tu-secret-key>
# Default region name: us-east-1
# Default output format: json
```

**Paso 2:** Navegar a directorio de Terraform

```bash
cd infrastructure/terraform/aws
```

**Paso 3:** Inicializar Terraform

```bash
terraform init
```

**Paso 4:** Revisar plan de infraestructura

```bash
terraform plan
```

**Paso 5:** Aplicar infraestructura

```bash
terraform apply

# Confirmar con: yes
```

**Recursos creados:**
- VPC con subnets públicas y privadas
- RDS PostgreSQL (db.t3.micro)
- ECS Cluster
- ECS Service con Fargate
- Application Load Balancer
- CloudWatch Log Groups
- Security Groups

### 2.3 Deploy Application

**Paso 1:** Build y push de imagen Docker

```bash
# Autenticar con ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Build imagen
docker build -t ops-incident-hub:latest .

# Tag imagen
docker tag ops-incident-hub:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/ops-incident-hub:latest

# Push imagen
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/ops-incident-hub:latest
```

**Paso 2:** Actualizar ECS Task Definition

```bash
# Forzar nuevo deployment
aws ecs update-service \
  --cluster ops-incident-hub-cluster \
  --service ops-incident-hub-service \
  --force-new-deployment
```

### 2.4 Configure Monitoring

**CloudWatch Logs:**

```bash
# Ver logs de la aplicación
aws logs tail /ecs/ops-incident-hub --follow
```

**CloudWatch Alarms:**

```bash
# Crear alarma de CPU alta
aws cloudwatch put-metric-alarm \
  --alarm-name ops-incident-hub-high-cpu \
  --alarm-description "CPU utilization > 80%" \
  --metric-name CPUUtilization \
  --namespace AWS/ECS \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 2
```

### 2.5 Verify Deployment

```bash
# Obtener URL del Load Balancer
terraform output alb_dns_name

# Health check
curl http://<alb-dns-name>/health

# Verificar métricas
curl http://<alb-dns-name>/metrics
```

### 2.6 Access Endpoints

- **API:** http://<alb-dns-name>
- **Health:** http://<alb-dns-name>/health
- **Metrics:** http://<alb-dns-name>/metrics

---

## 3. Deployment Azure (Container Apps)

### 3.1 Prerequisites

- Azure CLI 2.50+
- Terraform 1.5+
- Azure Subscription
- Docker (para build de imagen)

### 3.2 Infrastructure Setup

**Paso 1:** Login en Azure

```bash
az login
```

**Paso 2:** Crear Resource Group

```bash
az group create \
  --name ops-incident-hub-rg \
  --location eastus
```

**Paso 3:** Navegar a directorio de Terraform

```bash
cd infrastructure/terraform/azure
```

**Paso 4:** Inicializar Terraform

```bash
terraform init
```

**Paso 5:** Aplicar infraestructura

```bash
terraform apply

# Confirmar con: yes
```

**Recursos creados:**
- Resource Group
- Azure Database for PostgreSQL
- Container Registry
- Container Apps Environment
- Container App
- Application Insights
- Log Analytics Workspace

### 3.3 Deploy Application

**Paso 1:** Build y push de imagen Docker

```bash
# Login en Azure Container Registry
az acr login --name opsincidenthubacr

# Build imagen
docker build -t ops-incident-hub:latest .

# Tag imagen
docker tag ops-incident-hub:latest opsincidenthubacr.azurecr.io/ops-incident-hub:latest

# Push imagen
docker push opsincidenthubacr.azurecr.io/ops-incident-hub:latest
```

**Paso 2:** Actualizar Container App

```bash
# Actualizar con nueva imagen
az containerapp update \
  --name ops-incident-hub-app \
  --resource-group ops-incident-hub-rg \
  --image opsincidenthubacr.azurecr.io/ops-incident-hub:latest
```

### 3.4 Configure Monitoring

**Application Insights:**

```bash
# Obtener Instrumentation Key
az monitor app-insights component show \
  --app ops-incident-hub-insights \
  --resource-group ops-incident-hub-rg \
  --query instrumentationKey
```

**Log Analytics:**

```bash
# Ver logs
az monitor log-analytics query \
  --workspace <workspace-id> \
  --analytics-query "ContainerAppConsoleLogs_CL | where ContainerAppName_s == 'ops-incident-hub-app' | order by TimeGenerated desc | limit 100"
```

### 3.5 Verify Deployment

```bash
# Obtener URL de la aplicación
az containerapp show \
  --name ops-incident-hub-app \
  --resource-group ops-incident-hub-rg \
  --query properties.configuration.ingress.fqdn

# Health check
curl https://<app-url>/health

# Verificar métricas
curl https://<app-url>/metrics
```

### 3.6 Access Endpoints

- **API:** https://<app-url>
- **Health:** https://<app-url>/health
- **Metrics:** https://<app-url>/metrics
- **Application Insights:** Azure Portal

---

## 4. Variables de Entorno

### 4.1 Variables Requeridas

| Variable | Descripción | Ejemplo | Requerida |
|----------|-------------|---------|-----------|
| `SERVER_PORT` | Puerto del servidor | `8080` | ✅ |
| `SERVER_HOST` | Host del servidor | `0.0.0.0` | ✅ |
| `DB_HOST` | Host de PostgreSQL | `localhost` | ✅ |
| `DB_PORT` | Puerto de PostgreSQL | `5432` | ✅ |
| `DB_USER` | Usuario de PostgreSQL | `postgres` | ✅ |
| `DB_PASSWORD` | Password de PostgreSQL | `secret` | ✅ |
| `DB_NAME` | Nombre de la base de datos | `ops_incident_hub` | ✅ |

### 4.2 Variables Opcionales

| Variable | Descripción | Default | Opcional |
|----------|-------------|---------|----------|
| `SERVER_READ_TIMEOUT` | Timeout de lectura | `10s` | ✅ |
| `SERVER_WRITE_TIMEOUT` | Timeout de escritura | `10s` | ✅ |
| `SERVER_SHUTDOWN_TIMEOUT` | Timeout de shutdown | `30s` | ✅ |
| `ENVIRONMENT` | Ambiente (dev/staging/prod) | `development` | ✅ |
| `DB_MAX_OPEN_CONNS` | Conexiones máximas abiertas | `25` | ✅ |
| `DB_MAX_IDLE_CONNS` | Conexiones máximas idle | `5` | ✅ |
| `DB_CONN_MAX_LIFETIME` | Lifetime de conexión | `5m` | ✅ |
| `DB_CONN_MAX_IDLE_TIME` | Tiempo máximo idle | `1m` | ✅ |
| `ALLOWED_ORIGINS` | Orígenes CORS permitidos | `*` | ✅ |
| `LOG_LEVEL` | Nivel de log (debug/info/warn/error) | `info` | ✅ |
| `LOG_FORMAT` | Formato de log (text/json) | `json` | ✅ |

### 4.3 Configuración por Ambiente

**Development:**
```env
ENVIRONMENT=development
LOG_LEVEL=debug
LOG_FORMAT=text
DB_MAX_OPEN_CONNS=10
ALLOWED_ORIGINS=http://localhost:3000
```

**Staging:**
```env
ENVIRONMENT=staging
LOG_LEVEL=info
LOG_FORMAT=json
DB_MAX_OPEN_CONNS=25
ALLOWED_ORIGINS=https://staging.example.com
```

**Production:**
```env
ENVIRONMENT=production
LOG_LEVEL=warn
LOG_FORMAT=json
DB_MAX_OPEN_CONNS=50
ALLOWED_ORIGINS=https://example.com,https://*.example.com
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
```

### 4.4 Secrets Management

**Local (Docker Compose):**
- Usar archivo `.env`
- NO commitear `.env` al repositorio

**AWS:**
- Usar AWS Secrets Manager
- Referenciar secrets en ECS Task Definition

```bash
# Crear secret
aws secretsmanager create-secret \
  --name ops-incident-hub/db-password \
  --secret-string "super-secret-password"
```

**Azure:**
- Usar Azure Key Vault
- Referenciar secrets en Container App

```bash
# Crear secret
az keyvault secret set \
  --vault-name ops-incident-hub-kv \
  --name db-password \
  --value "super-secret-password"
```

---

## 5. Troubleshooting

### 5.1 Problema: API no responde

**Síntomas:**
- `curl http://localhost:8080/health` no responde
- Container está corriendo pero no acepta conexiones

**Diagnóstico:**
```bash
# Verificar logs
docker-compose logs api

# Verificar que el container está corriendo
docker-compose ps

# Verificar puertos
netstat -tulpn | grep 8080
```

**Solución:**
```bash
# Reiniciar servicio
docker-compose restart api

# Si persiste, rebuild
docker-compose up -d --build api
```

### 5.2 Problema: Error de conexión a base de datos

**Síntomas:**
- Logs muestran: `failed to connect to database`
- Health check retorna `unhealthy`

**Diagnóstico:**
```bash
# Verificar que PostgreSQL está corriendo
docker-compose ps postgres

# Verificar logs de PostgreSQL
docker-compose logs postgres

# Intentar conexión manual
docker-compose exec postgres psql -U postgres -d ops_incident_hub
```

**Solución:**
```bash
# Verificar variables de entorno
docker-compose exec api env | grep DB_

# Reiniciar PostgreSQL
docker-compose restart postgres

# Verificar network
docker network inspect ops-incident-hub_default
```

### 5.3 Problema: Migraciones no se aplican

**Síntomas:**
- Logs muestran: `migration failed`
- Tablas no existen en la base de datos

**Diagnóstico:**
```bash
# Verificar tablas existentes
docker-compose exec postgres psql -U postgres -d ops_incident_hub -c "\dt"

# Verificar logs de migraciones
docker-compose logs api | grep migration
```

**Solución:**
```bash
# Aplicar migraciones manualmente
docker-compose exec api /app/api migrate up

# Si falla, verificar archivos de migración
ls -la backend/migrations/
```

### 5.4 Problema: Rate limiting muy agresivo

**Síntomas:**
- Requests retornan `429 Too Many Requests`
- Logs muestran: `rate limit exceeded`

**Diagnóstico:**
```bash
# Verificar configuración de rate limiting
docker-compose exec api env | grep RATE_

# Ver logs de rate limiting
docker-compose logs api | grep "rate limit"
```

**Solución:**
```bash
# Ajustar rate limit en código (middleware/rate_limiter.go)
# O esperar 1 minuto para que se limpie el cache
```

### 5.5 Problema: Circuit breaker abierto

**Síntomas:**
- Requests fallan con `503 Service Unavailable`
- Logs muestran: `circuit breaker is open`

**Diagnóstico:**
```bash
# Verificar estado de la base de datos
docker-compose exec postgres pg_isready

# Ver logs de circuit breaker
docker-compose logs api | grep "circuit breaker"
```

**Solución:**
```bash
# Esperar 30 segundos para que el circuit breaker intente recuperarse
# O reiniciar la base de datos si está caída
docker-compose restart postgres
```

---

## 6. Rollback

### 6.1 Rollback Local (Docker Compose)

```bash
# Detener servicios
docker-compose down

# Checkout a versión anterior
git checkout <commit-hash>

# Rebuild y levantar
docker-compose up -d --build
```

### 6.2 Rollback AWS (ECS)

**Opción 1: Rollback a Task Definition anterior**

```bash
# Listar Task Definitions
aws ecs list-task-definitions --family-prefix ops-incident-hub

# Actualizar servicio con Task Definition anterior
aws ecs update-service \
  --cluster ops-incident-hub-cluster \
  --service ops-incident-hub-service \
  --task-definition ops-incident-hub:5  # versión anterior
```

**Opción 2: Rollback con Terraform**

```bash
# Checkout a versión anterior del código
git checkout <commit-hash>

# Aplicar infraestructura anterior
cd infrastructure/terraform/aws
terraform apply
```

### 6.3 Rollback Azure (Container Apps)

**Opción 1: Rollback a revisión anterior**

```bash
# Listar revisiones
az containerapp revision list \
  --name ops-incident-hub-app \
  --resource-group ops-incident-hub-rg

# Activar revisión anterior
az containerapp revision activate \
  --name ops-incident-hub-app \
  --resource-group ops-incident-hub-rg \
  --revision <revision-name>
```

**Opción 2: Rollback con imagen anterior**

```bash
# Actualizar con imagen anterior
az containerapp update \
  --name ops-incident-hub-app \
  --resource-group ops-incident-hub-rg \
  --image opsincidenthubacr.azurecr.io/ops-incident-hub:<previous-tag>
```

### 6.4 Rollback de Base de Datos

**CUIDADO:** Rollback de base de datos puede causar pérdida de datos.

```bash
# Listar migraciones aplicadas
docker-compose exec api /app/api migrate version

# Rollback a versión específica
docker-compose exec api /app/api migrate down <version>

# Rollback 1 versión
docker-compose exec api /app/api migrate down 1
```

---

## 7. Verificación Post-Deployment

### 7.1 Checklist de Verificación

- [ ] Health check retorna `200 OK`
- [ ] Base de datos está accesible
- [ ] Migraciones aplicadas correctamente
- [ ] Logs se están generando
- [ ] Métricas están disponibles en `/metrics`
- [ ] Endpoints de API responden correctamente
- [ ] Rate limiting funciona
- [ ] CORS configurado correctamente
- [ ] Security headers presentes

### 7.2 Tests de Humo

```bash
# Health check
curl -f http://<url>/health || echo "FAIL: Health check"

# Crear incidente
curl -X POST http://<url>/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","description":"Test","severity":"low"}' \
  || echo "FAIL: Create incident"

# Listar incidentes
curl -f http://<url>/api/v1/incidents || echo "FAIL: List incidents"

# Métricas
curl -f http://<url>/metrics || echo "FAIL: Metrics"
```

---

## 8. Contacto y Soporte

**Documentación adicional:**
- [Guía de Operación](guia-operacion.md)
- [Arquitectura Consolidada](arquitectura-consolidada.md)
- [Troubleshooting Avanzado](guia-operacion.md#troubleshooting)

**Soporte:**
- GitHub Issues: https://github.com/fabianbele2605/ops-incident-hub/issues
- Email: fabian.bele@example.com

---

**Última actualización:** Marzo 2025  
**Versión:** 1.0
