# Fase 3 - Migraciones de Base de Datos

## 1. Objetivo

Implementar un sistema de migraciones automáticas que gestione el schema de la base de datos de forma versionada y reproducible.

## 2. Estructura de Archivos

```
backend/migrations/
├── 001_create_tables.up.sql    # Migración: crear tablas
└── 001_create_tables.down.sql  # Rollback: eliminar tablas
```

## 3. Sistema de Migraciones

### 3.1 Tabla de Control

Se crea automáticamente una tabla `schema_migrations`:

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
)
```

### 3.2 Flujo de Ejecución

1. Verificar tabla `schema_migrations`
2. Obtener migraciones ya aplicadas
3. Buscar archivos `.up.sql` en carpeta `migrations/`
4. Aplicar migraciones pendientes en orden
5. Registrar migración aplicada

### 3.3 Características

- ✅ **Automáticas** - Se ejecutan al iniciar la aplicación
- ✅ **Idempotentes** - Se pueden ejecutar múltiples veces
- ✅ **Versionadas** - Orden numérico (001, 002, ...)
- ✅ **Transaccionales** - Rollback automático si falla
- ✅ **Auditables** - Registro de cuándo se aplicaron

## 4. Migración 001: Create Tables

### 4.1 Tabla Users

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_role CHECK (role IN ('admin', 'operator', 'viewer'))
);
```

**Características:**
- UUID como primary key
- Email único (índice automático)
- Role con constraint CHECK
- Timestamp de creación

### 4.2 Tabla Incidents

```sql
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    CONSTRAINT check_severity CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT check_status CHECK (status IN ('open', 'assigned', 'in_progress', 'resolved', 'closed'))
);
```

**Características:**
- UUID como primary key
- Foreign keys a users
- ON DELETE SET NULL para assigned_to (opcional)
- ON DELETE RESTRICT para created_by (obligatorio)
- Constraints CHECK para severity y status
- Campo JSONB para metadata flexible
- Timestamps de auditoría

### 4.3 Índices

```sql
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_assigned_to ON incidents(assigned_to);
CREATE INDEX idx_incidents_created_by ON incidents(created_by);
CREATE INDEX idx_incidents_created_at ON incidents(created_at DESC);
CREATE INDEX idx_users_email ON users(email);
```

**Justificación:**
- `status` y `severity`: Filtros frecuentes
- `assigned_to` y `created_by`: Joins y filtros
- `created_at DESC`: Ordenamiento por fecha
- `email`: Búsquedas por email

### 4.4 Datos Iniciales

```sql
INSERT INTO users (id, email, name, role, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin@ops-incident-hub.com',
    'Admin User',
    'admin',
    NOW()
) ON CONFLICT (id) DO NOTHING;
```

**Usuario admin por defecto** para desarrollo y testing.

## 5. Rollback (Down Migration)

```sql
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_incidents_created_at;
DROP INDEX IF EXISTS idx_incidents_created_by;
DROP INDEX IF EXISTS idx_incidents_assigned_to;
DROP INDEX IF EXISTS idx_incidents_severity;
DROP INDEX IF EXISTS idx_incidents_status;

DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS users;
```

**Orden importante:** Primero índices, luego tablas (por foreign keys).

## 6. Ejecución de Migraciones

### 6.1 En Desarrollo (Docker)

Las migraciones se ejecutan automáticamente al iniciar el contenedor:

```go
if err := postgres.RunMigrations(db, "migrations"); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

### 6.2 Logs de Ejecución

```
Applied migration: 001_create_tables
Migrations applied successfully
```

### 6.3 Verificación Manual

```bash
docker exec -it ops-incident-hub-db psql -U postgres -d ops_incident_hub

# Ver migraciones aplicadas
SELECT * FROM schema_migrations;

# Ver tablas creadas
\dt

# Ver estructura de tabla
\d incidents
```

## 7. Buenas Prácticas

### 7.1 Naming Convention

- Formato: `NNN_descripcion.up.sql` y `NNN_descripcion.down.sql`
- Números secuenciales con padding (001, 002, ...)
- Nombres descriptivos en snake_case

### 7.2 Contenido de Migraciones

- ✅ Una migración = un cambio lógico
- ✅ Siempre incluir `IF NOT EXISTS` / `IF EXISTS`
- ✅ Incluir rollback (down migration)
- ✅ Probar rollback antes de aplicar
- ❌ No modificar migraciones ya aplicadas

### 7.3 Constraints y Validaciones

- Usar CHECK constraints para valores enum
- Usar NOT NULL para campos obligatorios
- Usar UNIQUE para campos únicos
- Usar foreign keys para integridad referencial

## 8. Decisiones Técnicas

### 8.1 Migraciones Automáticas vs Manual

**Decisión:** Automáticas al inicio de la aplicación  
**Justificación:**
- Simplifica despliegues
- Garantiza schema actualizado
- Reduce errores humanos
- Ideal para desarrollo

**Alternativa en producción:** Usar herramientas como `golang-migrate` o `flyway` con control manual.

### 8.2 UUID vs Auto-increment

**Decisión:** UUID como primary key  
**Justificación:**
- Generación distribuida sin colisiones
- No expone información de volumen
- Compatible con sistemas distribuidos
- Mejor para APIs públicas

### 8.3 JSONB para Metadata

**Decisión:** Campo `metadata` como JSONB  
**Justificación:**
- Flexibilidad para datos no estructurados
- Queries eficientes con índices GIN
- Validación de JSON por PostgreSQL
- Extensibilidad sin cambios de schema

## 9. Validaciones Realizadas

```bash
✅ Migraciones se ejecutan automáticamente
✅ Tablas creadas correctamente
✅ Índices aplicados
✅ Constraints funcionando
✅ Usuario admin creado
✅ Foreign keys validadas
```

## 10. Próximos Pasos

Para agregar nuevas migraciones:

1. Crear archivos `002_nombre.up.sql` y `002_nombre.down.sql`
2. Escribir SQL de migración
3. Escribir SQL de rollback
4. Probar localmente
5. Commitear y desplegar

---

**Migraciones de Base de Datos:** ✅ **IMPLEMENTADAS**
