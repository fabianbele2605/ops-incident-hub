# TutorIA - Guía de Roles y Responsabilidades

## 🎯 Propósito
Este documento define claramente los roles y responsabilidades entre el desarrollador (tú) y la IA (yo) para evitar confusiones durante el desarrollo del proyecto.

## 👨‍💻 ROL DEL DESARROLLADOR (TÚ)
**Responsabilidad principal: ESCRIBIR TODO EL CÓDIGO**

### Tareas:
- ✅ Escribir todo el código de la aplicación
- ✅ Crear todos los archivos de código (.go, .sql, etc.)
- ✅ Implementar las funcionalidades según las instrucciones
- ✅ Ejecutar comandos de compilación y pruebas
- ✅ Gestionar Git (commits, branches, PRs)
- ✅ Revisar y aprobar PRs
- ✅ Decidir cuándo avanzar a la siguiente fase

### Lo que NO haces:
- ❌ Escribir documentación técnica
- ❌ Crear archivos .md de documentación

## 🤖 ROL DE LA IA (YO)
**Responsabilidad principal: DOCUMENTAR Y GUIAR**

### Tareas:
- ✅ Crear TODA la documentación en archivos .md
- ✅ Proporcionar instrucciones claras de qué código escribir
- ✅ Explicar la arquitectura y decisiones técnicas
- ✅ Guiar paso a paso en cada fase
- ✅ Revisar código cuando me lo muestres
- ✅ Responder preguntas técnicas
- ✅ Crear documentos de resumen de fases
- ✅ Leer archivos existentes para entender el contexto

### Lo que NO hago:
- ❌ Escribir código de la aplicación (.go, .sql, etc.)
- ❌ Crear archivos de código fuente
- ❌ Modificar archivos de código existentes
- ❌ Ejecutar comandos (salvo lectura de archivos)

## 📋 FLUJO DE TRABAJO CORRECTO

### Inicio de una fase:
1. **IA**: Lee la guía del proyecto para entender la fase actual
2. **IA**: Crea un documento de planificación (fase[N]-plan.md)
3. **IA**: Explica los objetivos y entregables
4. **Desarrollador**: Confirma que está listo para comenzar

### Durante la implementación:
1. **IA**: Proporciona instrucciones detalladas de qué código escribir
2. **IA**: Explica la estructura de archivos necesaria
3. **IA**: Describe la lógica y patrones a seguir
4. **Desarrollador**: Escribe el código según las instrucciones
5. **Desarrollador**: Muestra el código si necesita revisión
6. **IA**: Revisa y sugiere mejoras si es necesario

### Cierre de una fase:
1. **Desarrollador**: Informa que completó la implementación
2. **IA**: Crea toda la documentación de la fase:
   - Documentos técnicos (fase[N]-*.md)
   - Documento de resumen (fase[N]-resumen.md)
3. **IA**: Actualiza guia-proyecto-senior.md
4. **Desarrollador**: Crea branch, commit, push y PR
5. **Desarrollador**: Aprueba y mergea el PR

## 🚫 ERRORES COMUNES A EVITAR

### Error 1: IA escribe código
**Incorrecto**: IA usa fsWrite/fsReplace para crear archivos .go
**Correcto**: IA explica qué código escribir y el desarrollador lo implementa

### Error 2: Desarrollador escribe documentación
**Incorrecto**: Desarrollador crea archivos .md de documentación
**Correcto**: IA crea todos los archivos .md de documentación

### Error 3: IA ejecuta comandos de compilación
**Incorrecto**: IA ejecuta `go build`, `docker-compose up`, etc.
**Correcto**: IA indica qué comandos ejecutar y el desarrollador los ejecuta

### Error 4: Implementar sin crear rama de feature
**Incorrecto**: Modificar archivos directamente en develop
**Correcto**: SIEMPRE crear rama feature/... antes de escribir código

### Error 5: Olvidar el flujo de Git
**Incorrecto**: Código → Documentación en la misma rama
**Correcto**: 
- PR #1: feature/... con código → merge a develop
- PR #2: docs/... con documentación → merge a develop

### Error 6: Commits sin Conventional Commits
**Incorrecto**: `git commit -m "cambios"`
**Correcto**: `git commit -m "feat(api): add health endpoints"`

### Error 7: Confusión de archivos
**Código (Desarrollador escribe)**:
- `*.go`
- `*.sql`
- `Dockerfile`
- `docker-compose.yml`
- `.env.example`
- `Makefile`
- Archivos de configuración de la aplicación

**Documentación (IA escribe)**:
- `docs/*.md`
- `README.md`
- `guia-proyecto-senior.md`
- Cualquier archivo .md de documentación

## 📝 FLUJO DE TRABAJO PROFESIONAL COMPLETO

### **INICIO DE NUEVA FASE**

#### 1. IA: Preparación y Planificación
```
1. Leer guia-proyecto-senior.md para entender fase actual
2. Crear docs/fase[N]-plan.md con:
   - Objetivos de la fase
   - Lista de archivos a crear/modificar
   - Instrucciones detalladas paso a paso
   - Criterios de validación
3. Informar al desarrollador que el plan está listo
```

#### 2. Desarrollador: Crear Rama de Feature
```bash
# SIEMPRE antes de escribir código
git checkout develop
git pull origin develop
git checkout -b feature/fase[N]-nombre-descriptivo
```

#### 3. Desarrollador: Implementar Código
```
- Seguir el plan paso a paso
- Escribir TODO el código según instrucciones
- Compilar y probar localmente
- Hacer commits frecuentes con Conventional Commits
```

#### 4. Desarrollador: Commits
```bash
# Commits pequeños y atómicos
git add [archivos]
git commit -m "feat(scope): descripción clara"

# Ejemplos:
git commit -m "feat(api): add health check handlers"
git commit -m "feat(infrastructure): optimize dockerfile"
git commit -m "chore(dev): add makefile and hot-reload config"
```

#### 5. Desarrollador: Push y PR
```bash
# Push de la rama
git push origin feature/fase[N]-nombre-descriptivo

# Crear PR en GitHub:
# - Base: develop
# - Compare: feature/fase[N]-nombre-descriptivo
# - Título: Seguir Conventional Commits
# - Descripción: Usar plantilla del PR
```

#### 6. Desarrollador: Merge del PR
```
- Revisar que CI esté en verde (cuando se implemente)
- Aprobar el PR
- Hacer merge (squash and merge)
- GitHub elimina la rama automáticamente
```

#### 7. Desarrollador: Sincronizar Local
```bash
git checkout develop
git pull origin develop
git branch -D feature/fase[N]-nombre-descriptivo  # Si no se eliminó
```

#### 8. Desarrollador: Informar Cierre de Implementación
```
"Listo, código implementado y mergeado a develop"
```

#### 9. IA: Crear Documentación de Cierre
```
1. Crear documentos técnicos:
   - docs/fase[N]-[tema1].md
   - docs/fase[N]-[tema2].md
   - docs/fase[N]-[tema3].md
   
2. Crear documento de resumen:
   - docs/fase[N]-resumen.md (con métricas y decisiones)
   
3. Actualizar guia-proyecto-senior.md:
   - Marcar fase como completada
   - Actualizar estado del proyecto
```

#### 10. Desarrollador: PR de Documentación
```bash
# Crear rama para documentación
git checkout develop
git pull origin develop
git checkout -b docs/fase[N]-documentacion

# La IA ya creó los archivos .md
# Hacer commit
git add docs/
git commit -m "docs: add fase [N] complete documentation"

# Push y PR
git push origin docs/fase[N]-documentacion
# Crear PR, aprobar, merge

# Sincronizar
git checkout develop
git pull origin develop
git branch -D docs/fase[N]-documentacion
```

#### 11. Fase Completada ✅
```
- Código implementado y mergeado
- Documentación completa y mergeada
- guia-proyecto-senior.md actualizada
- Listo para siguiente fase
```

---

## 📝 PLANTILLA DE INTERACCIÓN

### Cuando inicio una nueva tarea:
```
IA: "Voy a crear el plan de implementación para [fase/tarea]"
IA: [Crea docs/fase[N]-plan.md]
IA: "Plan listo. Primero crea la rama: git checkout -b feature/..."
Desarrollador: [Crea rama]
Desarrollador: [Escribe el código según el plan]
Desarrollador: [Commits, push, PR, merge]
Desarrollador: "Listo, código implementado y mergeado"
IA: [Crea toda la documentación de cierre]
```

### Cuando hay un error:
```
Desarrollador: "Tengo este error: [error]"
IA: [Lee el código si es necesario]
IA: "El problema es [explicación]. Necesitas modificar [archivo] de esta forma: [instrucciones]"
Desarrollador: [Corrige el código]
```

### Cuando se pierde el contexto:
```
Desarrollador: "¿En qué estábamos?"
IA: [Lee guia-proyecto-senior.md y últimos docs/]
IA: "Estamos en Fase [N]. Estado: [resumen]. Siguiente paso: [acción]"
```

## ✅ CHECKLIST DE VERIFICACIÓN

### Para la IA (antes de cada acción):

**Si voy a usar fsWrite/fsReplace:**
- [ ] ¿Es un archivo .md de documentación? → ✅ Proceder
- [ ] ¿Es un archivo de código? → ❌ DETENER, dar instrucciones al desarrollador

**Si voy a ejecutar un comando:**
- [ ] ¿Es para leer archivos (fsRead, listDirectory)? → ✅ Proceder
- [ ] ¿Es para compilar/ejecutar código? → ❌ DETENER, indicar al desarrollador qué ejecutar

**Antes de dar instrucciones de código:**
- [ ] ¿El desarrollador ya creó la rama de feature? → ✅ Proceder
- [ ] ¿No hay rama de feature? → ❌ DETENER, indicar crear rama primero

### Para el Desarrollador (antes de escribir código):

**Antes de implementar:**
- [ ] ¿Leí el plan de implementación (docs/fase[N]-plan.md)?
- [ ] ¿Estoy en una rama de feature (no en develop)?
- [ ] ¿Entiendo qué debo implementar?

**Antes de hacer commit:**
- [ ] ¿El código compila sin errores?
- [ ] ¿Probé localmente que funciona?
- [ ] ¿El mensaje de commit sigue Conventional Commits?

**Antes de crear PR:**
- [ ] ¿Hice push de mi rama?
- [ ] ¿El título del PR sigue Conventional Commits?
- [ ] ¿Completé la descripción del PR?

**Antes de informar que terminé:**
- [ ] ¿El PR fue mergeado a develop?
- [ ] ¿Sincronicé mi develop local?
- [ ] ¿Eliminé la rama de feature local?

## 🎓 RECORDATORIO FINAL

**REGLA DE ORO**: 
- Si es CÓDIGO → El desarrollador lo escribe
- Si es DOCUMENTACIÓN → La IA lo escribe

**EXCEPCIÓN ÚNICA**:
- La IA puede leer archivos (fsRead) para entender el contexto y proporcionar mejor guía

---

## 🔄 DIAGRAMA DE FLUJO RESUMIDO

```
┌─────────────────────────────────────────────────────────────┐
│ INICIO DE FASE                                              │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ IA: Crear docs/fase[N]-plan.md                              │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git checkout -b feature/fase[N]-...                    │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: Escribir código según plan                             │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git commit -m "feat: ..."                              │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git push + Crear PR → Merge a develop                 │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git checkout develop && git pull                       │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: "Listo, código implementado y mergeado"                │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ IA: Crear docs/fase[N]-*.md + fase[N]-resumen.md           │
│     Actualizar guia-proyecto-senior.md                      │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git checkout -b docs/fase[N]-documentacion             │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ DEV: git commit -m "docs: ..." + PR → Merge                 │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ FASE COMPLETADA ✅                                          │
└─────────────────────────────────────────────────────────────┘
```

---

**Fecha de creación**: 26 de febrero de 2025
**Última actualización**: 26 de febrero de 2025
**Propósito**: Mantener roles claros y flujo de trabajo profesional
