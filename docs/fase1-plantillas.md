# Fase 1 - Plantillas de GitHub

## 1. Plantilla de Pull Request

### 1.1 Ubicación
Crear archivo: `.github/PULL_REQUEST_TEMPLATE.md`

### 1.2 Contenido que debes escribir:

```markdown
## 📋 Descripción

<!-- Describe claramente qué cambios introduce este PR y por qué son necesarios -->

## 🔧 Tipo de cambio

<!-- Marca con 'x' las opciones que apliquen -->

- [ ] 🐛 Bug fix (cambio que corrige un issue)
- [ ] ✨ Nueva feature (cambio que agrega funcionalidad)
- [ ] 💥 Breaking change (fix o feature que causa que funcionalidad existente no funcione como antes)
- [ ] 📝 Documentación (cambios solo en documentación)
- [ ] ♻️ Refactorización (cambio que no corrige bug ni agrega feature)
- [ ] ⚡ Performance (cambio que mejora rendimiento)
- [ ] ✅ Tests (agregar o corregir tests)
- [ ] 🔧 Chore (cambios en build, CI, dependencias)

## 🔗 Issues relacionados

<!-- Referencia a issues que este PR resuelve o está relacionado -->

Closes #
Related to #

## 🧪 Testing

<!-- Describe cómo probaste estos cambios -->

- [ ] Tests unitarios agregados/actualizados
- [ ] Tests de integración agregados/actualizados
- [ ] Probado manualmente en local
- [ ] Probado en entorno de desarrollo

### Pasos para probar:

1. 
2. 
3. 

## 📸 Screenshots (si aplica)

<!-- Agrega capturas de pantalla si hay cambios visuales -->

## ✅ Checklist

<!-- Verifica que todo esté completo antes de solicitar review -->

- [ ] Mi código sigue las convenciones de estilo del proyecto
- [ ] He realizado self-review de mi código
- [ ] He comentado mi código en áreas complejas
- [ ] He actualizado la documentación correspondiente
- [ ] Mis cambios no generan nuevos warnings
- [ ] He agregado tests que prueban que mi fix funciona o que mi feature funciona
- [ ] Tests unitarios nuevos y existentes pasan localmente
- [ ] He verificado que no hay conflictos con la rama base
- [ ] He seguido Conventional Commits en mis mensajes de commit

## 📚 Documentación adicional

<!-- Enlaces a documentación relevante, ADRs, diseños, etc -->

## 🚀 Deployment notes

<!-- Notas especiales para deployment, migraciones, variables de entorno, etc -->

## 👀 Reviewers

<!-- Menciona a personas específicas si necesitas su review -->

@username

---

**Nota para reviewers:** Por favor revisen especialmente [área específica que necesita atención]
```

### 1.3 Explicación:

**Propósito:** Estandarizar la información que debe incluir cada PR para facilitar el review.

**Secciones clave:**
- **Descripción:** Contexto del cambio
- **Tipo de cambio:** Clasificación rápida
- **Issues relacionados:** Trazabilidad
- **Testing:** Cómo se validó
- **Checklist:** Verificación de calidad
- **Screenshots:** Evidencia visual si aplica

## 2. Plantillas de Issues

### 2.1 Bug Report

#### Ubicación:
Crear archivo: `.github/ISSUE_TEMPLATE/bug_report.md`

#### Contenido que debes escribir:

```markdown
---
name: 🐛 Bug Report
about: Reportar un bug para ayudarnos a mejorar
title: '[BUG] '
labels: bug
assignees: ''
---

## 🐛 Descripción del Bug

<!-- Descripción clara y concisa del bug -->

## 🔄 Pasos para Reproducir

1. 
2. 
3. 
4. 

## ✅ Comportamiento Esperado

<!-- Describe qué debería suceder -->

## ❌ Comportamiento Actual

<!-- Describe qué está sucediendo actualmente -->

## 📸 Screenshots

<!-- Si aplica, agrega screenshots para explicar el problema -->

## 🌍 Entorno

- **Entorno:** [dev/staging/prod]
- **Navegador:** [Chrome, Firefox, Safari, etc]
- **Versión:** [v1.2.3]
- **OS:** [Windows, macOS, Linux]

## 📋 Logs/Errores

<!-- Pega logs relevantes o mensajes de error -->

```
[Pega logs aquí]
```

## 🔍 Contexto Adicional

<!-- Cualquier otra información relevante sobre el problema -->

## 💡 Posible Solución

<!-- Si tienes idea de cómo solucionarlo, compártela -->

## 🎯 Prioridad Sugerida

- [ ] 🔴 Crítica (bloquea funcionalidad principal)
- [ ] 🟠 Alta (afecta funcionalidad importante)
- [ ] 🟡 Media (afecta funcionalidad secundaria)
- [ ] 🟢 Baja (mejora menor)
```

#### Explicación:

**Propósito:** Capturar toda la información necesaria para reproducir y solucionar un bug.

**Secciones clave:**
- **Pasos para reproducir:** Fundamental para debugging
- **Comportamiento esperado vs actual:** Clarifica el problema
- **Entorno:** Ayuda a identificar si es específico de plataforma
- **Logs:** Evidencia técnica del error

### 2.2 Feature Request

#### Ubicación:
Crear archivo: `.github/ISSUE_TEMPLATE/feature_request.md`

#### Contenido que debes escribir:

```markdown
---
name: ✨ Feature Request
about: Sugerir una nueva funcionalidad
title: '[FEATURE] '
labels: enhancement
assignees: ''
---

## 🎯 Problema/Necesidad

<!-- Describe el problema que esta feature resolvería -->
<!-- Ejemplo: "Como usuario, no puedo filtrar incidentes por fecha, lo cual dificulta..." -->

## 💡 Solución Propuesta

<!-- Describe cómo te gustaría que funcionara la feature -->

## 🔄 Alternativas Consideradas

<!-- Describe alternativas que hayas considerado -->

## 📊 Beneficios

<!-- ¿Qué valor aporta esta feature? -->

- 
- 
- 

## 📸 Mockups/Diseños (si aplica)

<!-- Agrega mockups, wireframes o diseños si los tienes -->

## 🎯 Criterios de Aceptación

<!-- Define cuándo esta feature se considerará completa -->

- [ ] 
- [ ] 
- [ ] 

## 🔗 Referencias

<!-- Enlaces a documentación, ejemplos en otros sistemas, etc -->

## 🚀 Prioridad Sugerida

- [ ] 🔴 Crítica (necesaria para MVP)
- [ ] 🟠 Alta (importante para usuarios)
- [ ] 🟡 Media (mejora significativa)
- [ ] 🟢 Baja (nice to have)

## 📝 Notas Adicionales

<!-- Cualquier otra información relevante -->
```

#### Explicación:

**Propósito:** Capturar ideas de nuevas funcionalidades con suficiente detalle para evaluación.

**Secciones clave:**
- **Problema/Necesidad:** Justifica por qué es necesaria
- **Solución propuesta:** Describe la implementación deseada
- **Criterios de aceptación:** Define cuándo está completa
- **Prioridad:** Ayuda a priorizar el backlog

### 2.3 Documentation

#### Ubicación:
Crear archivo: `.github/ISSUE_TEMPLATE/documentation.md`

#### Contenido que debes escribir:

```markdown
---
name: 📝 Documentation
about: Mejora o corrección en documentación
title: '[DOCS] '
labels: documentation
assignees: ''
---

## 📚 Tipo de Documentación

- [ ] 📖 Documentación técnica
- [ ] 🎓 Tutorial/Guía
- [ ] 📋 README
- [ ] 🏗️ Arquitectura
- [ ] 🔧 Runbook
- [ ] 📊 ADR (Architecture Decision Record)
- [ ] 🐛 Corrección de error en docs existente

## 📝 Descripción

<!-- Describe qué documentación necesita ser creada o actualizada -->

## 📍 Ubicación

<!-- Dónde debería estar esta documentación -->

Archivo: `docs/...`

## 🎯 Audiencia

<!-- ¿Para quién es esta documentación? -->

- [ ] Desarrolladores
- [ ] Operadores/DevOps
- [ ] Usuarios finales
- [ ] Contribuidores
- [ ] Arquitectos

## ✅ Contenido Esperado

<!-- Lista de tópicos que debe cubrir -->

- [ ] 
- [ ] 
- [ ] 

## 🔗 Referencias

<!-- Enlaces a código, issues, o documentación relacionada -->

## 📋 Checklist

- [ ] Definir estructura del documento
- [ ] Escribir contenido
- [ ] Agregar diagramas (si aplica)
- [ ] Agregar ejemplos de código (si aplica)
- [ ] Revisar ortografía y gramática
- [ ] Validar con equipo
```

#### Explicación:

**Propósito:** Trackear necesidades de documentación de forma organizada.

**Secciones clave:**
- **Tipo:** Clasifica el tipo de documentación
- **Audiencia:** Define para quién es
- **Contenido esperado:** Guía de qué debe incluir

### 2.4 Config File

#### Ubicación:
Crear archivo: `.github/ISSUE_TEMPLATE/config.yml`

#### Contenido que debes escribir:

```yaml
blank_issues_enabled: false
contact_links:
  - name: 💬 Discusiones
    url: https://github.com/[tu-usuario]/azureSenior/discussions
    about: Para preguntas generales, usa Discussions en lugar de Issues
  - name: 📚 Documentación
    url: https://github.com/[tu-usuario]/azureSenior/tree/main/docs
    about: Revisa la documentación antes de crear un issue
```

#### Explicación:

**Propósito:** Configurar el comportamiento de issues y agregar enlaces útiles.

**Configuración:**
- `blank_issues_enabled: false`: Obliga a usar plantillas
- `contact_links`: Redirige a otros recursos antes de crear issue

## 3. Archivo CONTRIBUTING.md

### 3.1 Ubicación:
Crear archivo en raíz: `CONTRIBUTING.md`

### 3.2 Contenido que debes escribir:

```markdown
# Guía de Contribución - Ops Incident Hub

¡Gracias por tu interés en contribuir! Este documento te guiará en el proceso.

## 📋 Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [Cómo Contribuir](#cómo-contribuir)
- [Configuración del Entorno](#configuración-del-entorno)
- [Estándares de Código](#estándares-de-código)
- [Proceso de Pull Request](#proceso-de-pull-request)
- [Convenciones de Commits](#convenciones-de-commits)

## 🤝 Código de Conducta

Este proyecto sigue un código de conducta profesional. Se espera que todos los contribuidores:

- Sean respetuosos y constructivos
- Acepten críticas constructivas
- Se enfoquen en lo mejor para el proyecto
- Muestren empatía hacia otros miembros

## 🚀 Cómo Contribuir

### Reportar Bugs

1. Verifica que el bug no haya sido reportado antes
2. Usa la plantilla de Bug Report
3. Incluye pasos detallados para reproducir
4. Agrega logs y screenshots si es posible

### Sugerir Features

1. Verifica que no exista un issue similar
2. Usa la plantilla de Feature Request
3. Explica claramente el problema que resuelve
4. Proporciona ejemplos de uso

### Contribuir Código

1. Haz fork del repositorio
2. Crea una rama desde `develop`
3. Implementa tus cambios
4. Agrega tests
5. Actualiza documentación
6. Crea un Pull Request

## 🛠️ Configuración del Entorno

### Requisitos Previos

- Go 1.21+
- Node.js 18+
- Docker y Docker Compose
- Git
- Make (opcional)

### Setup Local

```bash
# 1. Clonar repositorio
git clone https://github.com/[tu-usuario]/azureSenior.git
cd azureSenior

# 2. Configurar backend
cd backend
go mod download
cp .env.example .env

# 3. Configurar frontend
cd ../frontend
npm install
cp .env.example .env.local

# 4. Levantar servicios con Docker
docker-compose up -d

# 5. Ejecutar migraciones
cd backend
make migrate-up

# 6. Ejecutar tests
make test
```

## 📏 Estándares de Código

### Backend (Go)

- Seguir [Effective Go](https://golang.org/doc/effective_go)
- Usar `gofmt` para formateo
- Ejecutar `golangci-lint` antes de commit
- Cobertura de tests mínima: 80%
- Documentar funciones públicas con godoc

```bash
# Formatear código
gofmt -w .

# Ejecutar linter
golangci-lint run

# Ejecutar tests con cobertura
go test -cover ./...
```

### Frontend (TypeScript)

- Seguir guía de estilo de Airbnb
- Usar ESLint y Prettier
- Componentes funcionales con hooks
- Props tipadas con TypeScript
- Tests con React Testing Library

```bash
# Formatear código
npm run format

# Ejecutar linter
npm run lint

# Ejecutar tests
npm test
```

### Documentación

- Markdown para toda la documentación
- Diagramas con Mermaid o Draw.io
- Ejemplos de código funcionales
- Mantener índice actualizado

## 🔄 Proceso de Pull Request

### Antes de Crear el PR

- [ ] Código formateado correctamente
- [ ] Tests agregados y pasando
- [ ] Linter sin errores
- [ ] Documentación actualizada
- [ ] Commits siguiendo Conventional Commits
- [ ] Rama actualizada con develop

### Crear el PR

1. Usa la plantilla de PR
2. Título siguiendo Conventional Commits
3. Descripción clara y completa
4. Referencia issues relacionados
5. Agrega screenshots si hay cambios visuales

### Durante el Review

- Responde a comentarios rápidamente
- Haz cambios solicitados en nuevos commits
- Resuelve conversaciones cuando estén completas
- Mantén el PR actualizado con develop

### Después del Merge

- Elimina tu rama
- Verifica que el deploy fue exitoso
- Cierra issues relacionados si no se cerraron automáticamente

## 📝 Convenciones de Commits

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```
<tipo>(<scope>): <descripción>

[cuerpo opcional]

[footer opcional]
```

### Tipos

- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Solo documentación
- `style`: Formato, sin cambio de lógica
- `refactor`: Refactorización
- `test`: Agregar o modificar tests
- `chore`: Tareas de mantenimiento
- `perf`: Mejora de rendimiento
- `ci`: Cambios en CI/CD

### Ejemplos

```bash
feat(api): add incident assignment endpoint
fix(auth): resolve token expiration issue
docs(readme): update installation steps
refactor(domain): simplify incident validation
test(usecase): add unit tests for create incident
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

# Tests e2e (cuando estén implementados)
npm run test:e2e
```

## 📚 Recursos Adicionales

- [Documentación del Proyecto](./docs/README.md)
- [Arquitectura](./docs/fase0-arquitectura-v1.md)
- [Decisiones Técnicas](./docs/fase0-decisiones-tecnicas.md)
- [Git Workflow](./docs/fase1-git-workflow.md)

## ❓ Preguntas

Si tienes preguntas:

1. Revisa la [documentación](./docs/)
2. Busca en [issues existentes](https://github.com/[tu-usuario]/azureSenior/issues)
3. Crea un [Discussion](https://github.com/[tu-usuario]/azureSenior/discussions)

## 🙏 Agradecimientos

¡Gracias por contribuir a Ops Incident Hub!
```

### 3.3 Explicación:

**Propósito:** Guía completa para nuevos contribuidores sobre cómo participar en el proyecto.

**Secciones clave:**
- **Setup del entorno:** Instrucciones paso a paso
- **Estándares de código:** Herramientas y convenciones
- **Proceso de PR:** Flujo completo de contribución
- **Testing:** Cómo ejecutar y escribir tests

## 4. Resumen de Archivos a Crear

| Archivo | Ubicación | Propósito |
|---------|-----------|-----------|
| PULL_REQUEST_TEMPLATE.md | `.github/` | Plantilla de PRs |
| bug_report.md | `.github/ISSUE_TEMPLATE/` | Reportar bugs |
| feature_request.md | `.github/ISSUE_TEMPLATE/` | Solicitar features |
| documentation.md | `.github/ISSUE_TEMPLATE/` | Issues de documentación |
| config.yml | `.github/ISSUE_TEMPLATE/` | Configuración de issues |
| CONTRIBUTING.md | Raíz del proyecto | Guía de contribución |

## 5. Próximos Pasos

1. Crear todos los archivos listados arriba
2. Personalizar con tu usuario de GitHub
3. Probar creando un issue de prueba
4. Probar creando un PR de prueba
5. Ajustar plantillas según necesidades del equipo
