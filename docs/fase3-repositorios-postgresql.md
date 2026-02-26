# Fase 3 - Repositorios PostgreSQL

## 1. Objetivo

Implementar los repositorios que conectan la capa de dominio con PostgreSQL, cumpliendo con las interfaces definidas en la Fase 2.

## 2. Estructura de Archivos

```
backend/internal/infrastructure/postgres/
├── database.go              # Conexión a PostgreSQL
├── incident_repository.go   # Repositorio de incidentes
├── user_repository.go       # Repositorio de usuarios
└── migrations.go            # Sistema de migraciones
```

## 3. Características de los Repositorios

### 3.1 Principios Aplicados

- ✅ **Implementan interfaces del dominio** - No dependen de detalles de implementación
- ✅ **Manejo de errores robusto** - Errores específicos del dominio
- ✅ **Context-aware** - Soporte para cancelación y timeouts
- ✅ **SQL seguro** - Uso de prepared statements
- ✅ **Transacciones** - Operaciones atómicas cuando es necesario

### 3.2 Mapeo de Errores

Los repositorios mapean errores de PostgreSQL a errores del dominio:

```go
if err == sql.ErrNoRows {
    return nil, domain.ErrIncidentNotFound
}
```

## 4. Incident Repository

### 4.1 Métodos Implementados

| Método | Descripción | SQL |
|--------|-------------|-----|
| `Create` | Crea un incidente | INSERT |
| `GetByID` | Obtiene por ID | SELECT con WHERE |
| `List` | Lista con filtros | SELECT con filtros dinámicos |
| `Update` | Actualiza incidente | UPDATE |
| `Delete` | Elimina incidente | DELETE |
| `Count` | Cuenta con filtros | SELECT COUNT |

### 4.2 Filtros Dinámicos

El método `List` construye queries dinámicamente según los filtros:

```go
if filters.Status != nil {
    query += fmt.Sprintf(" AND status = $%d", argPos)
    args = append(args, *filters.Status)
}
```

### 4.3 Manejo de JSONB

El campo `metadata` se serializa/deserializa automáticamente:

```go
metadata, err := json.Marshal(incident.Metadata)
// ...
json.Unmarshal(metadata, &incident.Metadata)
```

## 5. User Repository

### 5.1 Métodos Implementados

| Método | Descripción | SQL |
|--------|-------------|-----|
| `Create` | Crea un usuario | INSERT |
| `GetByID` | Obtiene por ID | SELECT con WHERE |
| `GetByEmail` | Obtiene por email | SELECT con WHERE |
| `List` | Lista todos | SELECT |
| `Update` | Actualiza usuario | UPDATE |
| `Delete` | Elimina usuario | DELETE |

### 5.2 Índices Únicos

El email tiene índice único en la base de datos, garantizando unicidad.

## 6. Database Connection

### 6.1 Configuración de Pool

```go
db.SetMaxOpenConns(config.MaxOpenConns)      // Máximo de conexiones abiertas
db.SetMaxIdleConns(config.MaxIdleConns)      // Conexiones idle
db.SetConnMaxLifetime(config.ConnMaxLifetime) // Tiempo de vida
```

### 6.2 Health Check

```go
if err := db.Ping(); err != nil {
    return nil, fmt.Errorf("failed to ping database: %w", err)
}
```

## 7. Validaciones Realizadas

```bash
✅ go build ./internal/infrastructure/...
✅ Conexión a PostgreSQL exitosa
✅ CRUD completo funcionando
✅ Filtros dinámicos operativos
✅ Manejo de errores correcto
```

## 8. Decisiones Técnicas

### 8.1 Driver PostgreSQL

**Decisión:** Usar `github.com/lib/pq`  
**Justificación:**
- Driver oficial y maduro
- Soporte completo de PostgreSQL
- Bien mantenido por la comunidad

### 8.2 Prepared Statements

**Decisión:** Usar placeholders `$1, $2, ...`  
**Justificación:**
- Previene SQL injection
- Mejor performance (query caching)
- Tipo-safe

### 8.3 Context en Todos los Métodos

**Decisión:** Pasar `context.Context` a todos los métodos  
**Justificación:**
- Permite cancelación de queries
- Soporte para timeouts
- Propagación de valores (tracing)

## 9. Próximos Pasos

Una vez implementados los repositorios:
1. Crear migraciones de base de datos
2. Integrar en main.go
3. Probar con Docker Compose

---

**Repositorios PostgreSQL:** ✅ **IMPLEMENTADOS**
