# Fase 1 - Git Workflow y Versionado

## 1. Estrategia de Branching: Git Flow Simplificado

### 1.1 Ramas Principales

#### main
- **Propósito:** Código en producción
- **Protección:** Requiere PR aprobado + CI pasando
- **Deploy:** Automático a producción
- **Nunca:** Hacer commits directos

#### develop
- **Propósito:** Integración de features
- **Protección:** Requiere PR + CI pasando
- **Deploy:** Automático a staging
- **Merge desde:** feature/*, fix/*, docs/*

### 1.2 Ramas Temporales

#### feature/nombre-descriptivo
- **Propósito:** Nuevas funcionalidades
- **Origen:** develop
- **Destino:** develop
- **Ejemplo:** feature/create-incident-endpoint
- **Duración:** Corta (1-5 días idealmente)

#### fix/nombre-descriptivo
- **Propósito:** Corrección de bugs
- **Origen:** develop (o main si es hotfix)
- **Destino:** develop (o main si es hotfix)
- **Ejemplo:** fix/incident-validation-error

#### docs/nombre-descriptivo
- **Propósito:** Solo cambios en documentación
- **Origen:** develop
- **Destino:** develop
- **Ejemplo:** docs/add-api-documentation

#### refactor/nombre-descriptivo
- **Propósito:** Refactorización sin cambio funcional
- **Origen:** develop
- **Destino:** develop
- **Ejemplo:** refactor/incident-repository-interface

#### hotfix/nombre-descriptivo
- **Propósito:** Corrección urgente en producción
- **Origen:** main
- **Destino:** main Y develop
- **Ejemplo:** hotfix/critical-auth-bug

## 2. Flujo de Trabajo

### 2.1 Desarrollo de Nueva Feature

```
1. Actualizar develop local:
   git checkout develop
   git pull origin develop

2. Crear rama de feature:
   git checkout -b feature/nombre-descriptivo

3. Desarrollar y hacer commits:
   git add .
   git commit -m "feat: descripción del cambio"

4. Subir rama al remoto:
   git push origin feature/nombre-descriptivo

5. Crear Pull Request en GitHub:
   - Base: develop
   - Compare: feature/nombre-descriptivo

6. Esperar aprobación y CI verde

7. Hacer merge (squash and merge recomendado)

8. Eliminar rama:
   git branch -d feature/nombre-descriptivo
   git push origin --delete feature/nombre-descriptivo
```

### 2.2 Corrección de Bug

```
1. Actualizar develop:
   git checkout develop
   git pull origin develop

2. Crear rama de fix:
   git checkout -b fix/nombre-bug

3. Corregir y hacer commits:
   git commit -m "fix: descripción de la corrección"

4. Subir y crear PR igual que feature

5. Merge a develop
```

### 2.3 Hotfix en Producción

```
1. Desde main:
   git checkout main
   git pull origin main

2. Crear rama de hotfix:
   git checkout -b hotfix/nombre-critico

3. Corregir el problema:
   git commit -m "fix: corrección crítica en producción"

4. Crear PR a main

5. Después de merge a main, también merge a develop:
   git checkout develop
   git merge main
   git push origin develop
```

### 2.4 Release a Producción

```
1. Validar que staging (develop) está estable

2. Crear PR de develop a main

3. Revisión exhaustiva del PR

4. Merge a main (crea tag automático)

5. Deploy automático a producción

6. Monitorear métricas y logs
```

## 3. Conventional Commits

### 3.1 Formato

```
<tipo>(<scope>): <descripción corta>

[cuerpo opcional]

[footer opcional]
```

### 3.2 Tipos de Commits

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| feat | Nueva funcionalidad | `feat(api): add create incident endpoint` |
| fix | Corrección de bug | `fix(auth): resolve token expiration issue` |
| docs | Solo documentación | `docs(readme): update installation steps` |
| style | Formato, sin cambio de lógica | `style(backend): format code with gofmt` |
| refactor | Refactorización | `refactor(domain): simplify incident validation` |
| test | Agregar o modificar tests | `test(usecase): add unit tests for create incident` |
| chore | Tareas de mantenimiento | `chore(deps): update go dependencies` |
| perf | Mejora de rendimiento | `perf(api): optimize database queries` |
| ci | Cambios en CI/CD | `ci(github): add security scanning to pipeline` |
| build | Cambios en build | `build(docker): optimize backend image size` |

### 3.3 Scope (Opcional)

Indica el área afectada:
- `api`: Capa de API/handlers
- `domain`: Lógica de negocio
- `infrastructure`: Implementaciones (DB, auth)
- `frontend`: Código del frontend
- `terraform`: Infraestructura
- `ci`: CI/CD
- `docs`: Documentación

### 3.4 Ejemplos Completos

**Feature simple:**
```
feat(api): add list incidents endpoint

Implements GET /api/v1/incidents with pagination support.
Includes filtering by status and severity.
```

**Fix con breaking change:**
```
fix(auth)!: change JWT token structure

BREAKING CHANGE: Token payload structure changed.
Clients must update to new token format.

Migration guide in docs/migration-v2.md
```

**Chore:**
```
chore(deps): update dependencies

- go 1.21 -> 1.22
- react 18.2 -> 18.3
```

### 3.5 Reglas de Commits

1. **Descripción en imperativo:** "add" no "added" ni "adds"
2. **Primera letra minúscula:** "add feature" no "Add feature"
3. **Sin punto final:** "add feature" no "add feature."
4. **Máximo 72 caracteres** en la primera línea
5. **Cuerpo opcional** separado por línea en blanco
6. **Referencias a issues:** "Closes #123" en el footer

## 4. Pull Request Guidelines

### 4.1 Título del PR

Debe seguir Conventional Commits:
```
feat(api): add incident assignment endpoint
fix(frontend): resolve dashboard loading issue
docs: update architecture documentation
```

### 4.2 Descripción del PR

Usar la plantilla (se creará en siguiente documento):

```markdown
## Descripción
[Descripción clara del cambio]

## Tipo de cambio
- [ ] Bug fix
- [ ] Nueva feature
- [ ] Breaking change
- [ ] Documentación

## Checklist
- [ ] Tests agregados/actualizados
- [ ] Documentación actualizada
- [ ] CI pasando
- [ ] Code review solicitado

## Testing
[Cómo se probó el cambio]

## Screenshots (si aplica)
[Capturas de pantalla]
```

### 4.3 Proceso de Review

1. **Autor crea PR** con descripción completa
2. **CI ejecuta automáticamente** (tests, linting, security)
3. **Reviewer asignado** revisa código
4. **Comentarios y cambios** si es necesario
5. **Aprobación** cuando todo está correcto
6. **Merge** (squash and merge preferido)
7. **Eliminar rama** automáticamente

### 4.4 Criterios de Aprobación

Un PR puede ser aprobado si:
- ✅ CI está en verde (todos los checks pasan)
- ✅ Al menos 1 aprobación de reviewer
- ✅ No hay conflictos con la rama base
- ✅ Descripción completa y clara
- ✅ Tests incluidos (si aplica)
- ✅ Documentación actualizada (si aplica)

## 5. Semantic Versioning (SemVer)

### 5.1 Formato

```
MAJOR.MINOR.PATCH

Ejemplo: 1.2.3
```

### 5.2 Incremento de Versión

| Tipo | Cuándo | Ejemplo |
|------|--------|---------|
| MAJOR | Breaking changes | 1.0.0 → 2.0.0 |
| MINOR | Nueva funcionalidad (compatible) | 1.0.0 → 1.1.0 |
| PATCH | Bug fixes (compatible) | 1.0.0 → 1.0.1 |

### 5.3 Pre-releases

```
1.0.0-alpha.1    # Primera versión alpha
1.0.0-beta.1     # Primera versión beta
1.0.0-rc.1       # Release candidate
1.0.0            # Release final
```

### 5.4 Tags en Git

Cada release debe tener un tag:

```bash
# Crear tag anotado
git tag -a v1.0.0 -m "Release version 1.0.0"

# Subir tag al remoto
git push origin v1.0.0

# Listar tags
git tag -l
```

### 5.5 Changelog Automático

Los commits siguiendo Conventional Commits permiten generar changelog automático:

```
# v1.1.0 (2024-01-15)

## Features
- feat(api): add incident assignment endpoint (#45)
- feat(frontend): add dashboard metrics (#47)

## Bug Fixes
- fix(auth): resolve token refresh issue (#46)

## Documentation
- docs: update API documentation (#48)
```

## 6. Protección de Ramas

### 6.1 Configuración para main

En GitHub Settings → Branches → Branch protection rules:

- ✅ Require pull request before merging
- ✅ Require approvals (mínimo 1)
- ✅ Dismiss stale pull request approvals when new commits are pushed
- ✅ Require status checks to pass before merging
  - CI/CD pipeline
  - Tests
  - Linting
  - Security scan
- ✅ Require branches to be up to date before merging
- ✅ Require conversation resolution before merging
- ✅ Do not allow bypassing the above settings
- ✅ Restrict who can push to matching branches (solo admins)

### 6.2 Configuración para develop

- ✅ Require pull request before merging
- ✅ Require status checks to pass before merging
- ✅ Require branches to be up to date before merging

## 7. Comandos Git Útiles

### 7.1 Actualizar rama con cambios de develop

```bash
# Opción 1: Rebase (historial limpio)
git checkout feature/mi-feature
git rebase develop

# Opción 2: Merge (preserva historial)
git checkout feature/mi-feature
git merge develop
```

### 7.2 Deshacer último commit (sin perder cambios)

```bash
git reset --soft HEAD~1
```

### 7.3 Ver historial de commits

```bash
# Formato compacto
git log --oneline --graph --all

# Con detalles
git log --graph --pretty=format:'%h - %s (%cr) <%an>'
```

### 7.4 Limpiar ramas locales eliminadas en remoto

```bash
git fetch --prune
git branch -vv | grep ': gone]' | awk '{print $1}' | xargs git branch -d
```

## 8. Buenas Prácticas

### 8.1 Commits

- ✅ Commits pequeños y atómicos
- ✅ Un commit = un cambio lógico
- ✅ Mensaje descriptivo y claro
- ✅ Commits frecuentes (no esperar días)
- ❌ No hacer commits de código que no compila
- ❌ No hacer commits con TODOs sin resolver

### 8.2 Branches

- ✅ Nombres descriptivos y cortos
- ✅ Ramas de vida corta (días, no semanas)
- ✅ Eliminar ramas después de merge
- ✅ Mantener ramas actualizadas con develop
- ❌ No acumular muchos commits sin merge
- ❌ No crear ramas desde ramas (solo desde develop/main)

### 8.3 Pull Requests

- ✅ PRs pequeños y enfocados
- ✅ Descripción completa y clara
- ✅ Screenshots si hay cambios visuales
- ✅ Responder a comentarios rápidamente
- ✅ Resolver conversaciones antes de merge
- ❌ No hacer PRs gigantes (>500 líneas)
- ❌ No hacer merge sin aprobación

## 9. Troubleshooting

### 9.1 Conflictos en Merge

```bash
# 1. Actualizar develop
git checkout develop
git pull origin develop

# 2. Volver a tu rama
git checkout feature/mi-feature

# 3. Hacer rebase
git rebase develop

# 4. Resolver conflictos manualmente en archivos

# 5. Continuar rebase
git add .
git rebase --continue

# 6. Forzar push (solo en ramas de feature)
git push origin feature/mi-feature --force-with-lease
```

### 9.2 Commit en rama equivocada

```bash
# 1. Guardar cambios
git stash

# 2. Cambiar a rama correcta
git checkout rama-correcta

# 3. Aplicar cambios
git stash pop

# 4. Hacer commit
git commit -m "mensaje"
```

## 10. Próximos Pasos

1. Configurar protección de ramas en GitHub
2. Crear plantillas de PR e issues
3. Configurar GitHub Actions para CI
4. Hacer primer commit siguiendo estas convenciones
5. Practicar flujo completo con una feature de prueba
