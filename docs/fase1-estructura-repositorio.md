# Fase 1 - Estructura del Repositorio

## 1. Estructura de Carpetas Definitiva

```
azureSenior/
├── .github/                      # Configuración de GitHub
│   ├── workflows/                # GitHub Actions (CI/CD)
│   ├── ISSUE_TEMPLATE/           # Plantillas de issues
│   └── PULL_REQUEST_TEMPLATE.md  # Plantilla de PR
│
├── backend/                      # Backend en Go
│   ├── cmd/
│   │   └── api/                  # Entry point de la API
│   │       └── main.go
│   ├── internal/                 # Código privado del backend
│   │   ├── domain/               # Entidades y lógica de negocio
│   │   ├── usecase/              # Casos de uso
│   │   ├── infrastructure/       # Implementaciones (DB, auth, etc)
│   │   └── api/                  # HTTP handlers
│   ├── pkg/                      # Código reutilizable público
│   ├── migrations/               # Migraciones de base de datos
│   ├── tests/                    # Tests de integración
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── README.md
│
├── frontend/                     # Frontend en TypeScript
│   ├── src/
│   │   ├── components/           # Componentes React
│   │   ├── pages/                # Páginas/vistas
│   │   ├── services/             # Llamadas a API
│   │   ├── hooks/                # Custom hooks
│   │   ├── types/                # TypeScript types
│   │   ├── utils/                # Utilidades
│   │   └── App.tsx
│   ├── public/
│   ├── package.json
│   ├── tsconfig.json
│   ├── Dockerfile
│   └── README.md
│
├── infrastructure/               # Infraestructura como código
│   ├── terraform/
│   │   ├── modules/              # Módulos reutilizables
│   │   │   ├── container-apps/
│   │   │   ├── database/
│   │   │   ├── networking/
│   │   │   └── monitoring/
│   │   ├── environments/         # Configuración por entorno
│   │   │   ├── dev/
│   │   │   ├── staging/
│   │   │   └── prod/
│   │   └── README.md
│   └── scripts/                  # Scripts de automatización
│
├── docs/                         # Documentación del proyecto
│   ├── README.md                 # Índice de documentación
│   ├── fase0-*.md                # Documentos de Fase 0
│   ├── fase1-*.md                # Documentos de Fase 1
│   ├── architecture/             # Diagramas y arquitectura
│   ├── runbooks/                 # Runbooks operativos
│   └── adr/                      # Architecture Decision Records
│
├── .gitignore                    # Archivos ignorados por Git
├── .editorconfig                 # Configuración de editor
├── README.md                     # README principal del proyecto
├── LICENSE                       # Licencia del proyecto
├── CONTRIBUTING.md               # Guía de contribución
└── guia-proyecto-senior.md       # Guía maestra del proyecto
```

## 2. Explicación de Carpetas Principales

### 2.1 Backend (Go)
**Propósito:** Contiene toda la lógica del servidor API

**Subcarpetas clave:**
- `cmd/api/`: Punto de entrada de la aplicación (main.go)
- `internal/`: Código privado que no puede ser importado por otros proyectos
  - `domain/`: Entidades de negocio (Incident, User) - sin dependencias externas
  - `usecase/`: Lógica de casos de uso (CreateIncident, AssignIncident)
  - `infrastructure/`: Implementaciones concretas (PostgreSQL, Azure AD)
  - `api/`: Handlers HTTP, middlewares, routing
- `pkg/`: Código reutilizable que puede ser importado externamente
- `migrations/`: Archivos SQL de migraciones de base de datos
- `tests/`: Tests de integración y e2e

**Razón:** Sigue Clean Architecture, separando responsabilidades y facilitando testing

### 2.2 Frontend (TypeScript + React)
**Propósito:** Interfaz de usuario de la aplicación

**Subcarpetas clave:**
- `src/components/`: Componentes React reutilizables
- `src/pages/`: Páginas completas de la aplicación
- `src/services/`: Lógica de comunicación con API
- `src/hooks/`: Custom hooks de React
- `src/types/`: Definiciones de tipos TypeScript
- `src/utils/`: Funciones auxiliares

**Razón:** Organización estándar de React, facilita escalabilidad

### 2.3 Infrastructure (Terraform)
**Propósito:** Infraestructura como código para Azure

**Subcarpetas clave:**
- `terraform/modules/`: Módulos reutilizables por recurso
- `terraform/environments/`: Configuración específica por entorno (dev/staging/prod)
- `scripts/`: Scripts de automatización (deploy, destroy, etc)

**Razón:** Separación de módulos y entornos permite reutilización y gestión independiente

### 2.4 Docs
**Propósito:** Toda la documentación técnica y operativa

**Subcarpetas clave:**
- `architecture/`: Diagramas y documentos de arquitectura
- `runbooks/`: Procedimientos operativos
- `adr/`: Architecture Decision Records (decisiones importantes)

**Razón:** Documentación centralizada y organizada por tipo

### 2.5 .github
**Propósito:** Configuración específica de GitHub

**Contenido:**
- `workflows/`: Pipelines de CI/CD con GitHub Actions
- `ISSUE_TEMPLATE/`: Plantillas para crear issues
- `PULL_REQUEST_TEMPLATE.md`: Plantilla para PRs

**Razón:** Automatización y estandarización del flujo de trabajo

## 3. Archivos de Configuración Raíz

### 3.1 .gitignore
**Propósito:** Definir qué archivos NO deben versionarse

**Debe incluir:**
- Archivos de build (binarios, node_modules, dist/)
- Archivos de configuración local (.env, *.local)
- Archivos de IDE (.vscode/, .idea/)
- Archivos de Terraform (*.tfstate, .terraform/)
- Archivos temporales (*.log, *.tmp)

### 3.2 .editorconfig
**Propósito:** Configuración consistente de editor entre desarrolladores

**Define:**
- Tipo de indentación (espacios vs tabs)
- Tamaño de indentación
- Charset (UTF-8)
- Fin de línea (LF)

### 3.3 README.md
**Propósito:** Punto de entrada del proyecto, primera impresión

**Debe contener:**
- Descripción del proyecto
- Tecnologías utilizadas
- Requisitos previos
- Instrucciones de instalación
- Comandos principales
- Enlaces a documentación
- Badges de estado (build, coverage, etc)

### 3.4 CONTRIBUTING.md
**Propósito:** Guía para contribuidores

**Debe contener:**
- Cómo configurar el entorno de desarrollo
- Estándares de código
- Proceso de PR
- Convenciones de commits

### 3.5 LICENSE
**Propósito:** Licencia del proyecto

**Recomendación:** MIT License (permisiva y común en proyectos open source)

## 4. Convenciones de Nombres

### 4.1 Archivos
- **Go:** snake_case para archivos (incident_repository.go)
- **TypeScript:** kebab-case para archivos (incident-list.tsx)
- **Documentación:** kebab-case (fase1-estructura.md)
- **Configuración:** lowercase con puntos (.gitignore, .editorconfig)

### 4.2 Carpetas
- **Todas:** lowercase con guiones (container-apps, pull-request-template)
- **Excepción:** Carpetas de Go siguen convención Go (internal, pkg)

### 4.3 Branches
- **main:** Rama principal (producción)
- **develop:** Rama de integración
- **feature/nombre-feature:** Nuevas funcionalidades
- **fix/nombre-fix:** Correcciones de bugs
- **docs/nombre-doc:** Cambios en documentación
- **refactor/nombre-refactor:** Refactorizaciones

## 5. Orden de Creación

Para crear esta estructura, sigue este orden:

1. **Crear carpetas principales:**
   ```
   backend/
   frontend/
   infrastructure/
   docs/
   .github/
   ```

2. **Crear subcarpetas de backend:**
   ```
   backend/cmd/api/
   backend/internal/domain/
   backend/internal/usecase/
   backend/internal/infrastructure/
   backend/internal/api/
   backend/pkg/
   backend/migrations/
   backend/tests/
   ```

3. **Crear subcarpetas de frontend:**
   ```
   frontend/src/components/
   frontend/src/pages/
   frontend/src/services/
   frontend/src/hooks/
   frontend/src/types/
   frontend/src/utils/
   frontend/public/
   ```

4. **Crear subcarpetas de infrastructure:**
   ```
   infrastructure/terraform/modules/
   infrastructure/terraform/environments/dev/
   infrastructure/terraform/environments/staging/
   infrastructure/terraform/environments/prod/
   infrastructure/scripts/
   ```

5. **Crear subcarpetas de docs:**
   ```
   docs/architecture/
   docs/runbooks/
   docs/adr/
   ```

6. **Crear subcarpetas de .github:**
   ```
   .github/workflows/
   .github/ISSUE_TEMPLATE/
   ```

7. **Crear archivos de configuración raíz:**
   ```
   .gitignore
   .editorconfig
   README.md
   CONTRIBUTING.md
   LICENSE
   ```

## 6. Archivos Placeholder

En cada carpeta vacía, crear un archivo `.gitkeep` para que Git la versione:

**Razón:** Git no versiona carpetas vacías, solo archivos. El `.gitkeep` es una convención para mantener la estructura.

**Carpetas que necesitan .gitkeep:**
- backend/pkg/
- backend/tests/
- frontend/public/
- docs/architecture/
- docs/runbooks/
- docs/adr/
- infrastructure/scripts/

## 7. Próximos Pasos

Una vez creada la estructura:

1. Crear archivos de configuración (.gitignore, .editorconfig)
2. Crear README.md principal
3. Crear plantillas de GitHub (PR, issues)
4. Configurar Git Flow
5. Hacer commit inicial con mensaje: `chore: initialize project structure`

## 8. Validación

Para validar que la estructura está correcta, ejecuta:

```bash
tree -L 3 -a
```

Deberías ver la estructura completa con todas las carpetas y archivos de configuración.
