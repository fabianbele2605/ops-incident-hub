# Enunciado Maestro - Proyecto Azure Senior

## 1. Vision del proyecto
Construir una plataforma cloud en Azure, lista para produccion, orientada a estandares senior: arquitectura limpia, seguridad por defecto, observabilidad completa, automatizacion de extremo a extremo y operacion confiable.

Este proyecto sera la evolucion natural de tu repositorio anterior, pero ahora con enfoque de producto real, escala y mantenibilidad.

## 2. Objetivo general
Diseñar, implementar y operar una aplicacion moderna sobre Azure con practicas de nivel senior, demostrando capacidad de:
- Diseñar arquitectura y decisiones tecnicas.
- Implementar CI/CD robusto con controles de calidad.
- Gestionar infraestructura como codigo con entornos separados.
- Aplicar seguridad en cada capa.
- Operar con monitoreo, alertas y runbooks.
- Escalar de forma controlada.

## 3. Alcance
El proyecto cubre:
- Backend principal y servicios de soporte.
- Infraestructura cloud en Azure.
- Plataforma de despliegue y ejecucion en contenedores.
- Pipeline de integracion y despliegue continuo.
- Observabilidad, seguridad, resiliencia y operacion.
- Documentacion tecnica y operativa.

El proyecto no cubre:
- Funcionalidades de negocio demasiado amplias sin valor tecnico.
- Integraciones no esenciales que desvien el foco.

## 4. Resultado esperado
Al finalizar, debes tener un repositorio senior que sirva para portafolio, entrevistas y referencia profesional:
- Arquitectura explicada y justificada.
- Automatizacion repetible en dev, staging y prod.
- Evidencia de calidad, seguridad y operacion.
- Documentacion clara para desarrolladores y operadores.

## 5. Principios de trabajo
- Todo cambio relevante queda versionado y documentado.
- Todo despliegue debe ser trazable y reversible.
- La seguridad no es fase final; se trabaja desde el inicio.
- Lo que no se observa, no se puede operar.
- Si no puede automatizarse, no esta listo para produccion.

## 6. Fases del proyecto (paso a paso)

## Fase 0 - Definicion y diseño
Objetivo:
Definir claramente que problema resuelve la aplicacion y como se vera la arquitectura objetivo.

Actividades:
- Definir alcance funcional minimo viable.
- Definir diagrama de arquitectura inicial.
- Definir criterios de aceptacion por fase.
- Definir convenciones tecnicas y estandares.

Entregables:
- Documento de vision.
- Diagrama de arquitectura v1.
- Lista de decisiones tecnicas iniciales.

Criterio de exito:
- Alcance claro, sin ambiguedad, con decisiones justificadas.

## Fase 1 - Fundacion del repositorio senior
Objetivo:
Establecer una base profesional de trabajo para todo el ciclo de vida.

Actividades:
- Estructurar monorepo o multi-repo con criterio claro.
- Definir estrategia de ramas y versionado.
- Definir estandar de commits y pull requests.
- Crear plantilla de documentacion tecnica.

Entregables:
- Estructura de carpetas definitiva.
- Politicas de colaboracion y calidad.
- Plantillas de PR, incidencias y cambios.

Criterio de exito:
- Flujo de trabajo claro para cualquier colaborador tecnico.

## Fase 2 - Arquitectura de aplicacion
Objetivo:
Diseñar la aplicacion con separacion de responsabilidades y mantenibilidad a largo plazo.

Actividades:
- Definir capas y fronteras de arquitectura.
- Definir contratos entre componentes.
- Definir estrategia de errores, validaciones y trazabilidad.
- Definir modelo de configuracion por entorno.

Entregables:
- Documento de arquitectura de aplicacion.
- Mapa de dependencias y responsabilidades.
- Catalogo de endpoints/casos de uso.

Criterio de exito:
- Arquitectura entendible, testeable y extensible.

## Fase 3 - Infraestructura como codigo avanzada
Objetivo:
Modelar infraestructura cloud de forma modular y reutilizable.

Actividades:
- Separar modulos por dominio de infraestructura.
- Definir estados remotos y politicas de bloqueo.
- Definir estrategia de entornos: dev, staging, prod.
- Definir naming convention y etiquetado obligatorio.

Entregables:
- Modulos reutilizables.
- Configuracion por entorno.
- Plan de evolucion de infraestructura.

Criterio de exito:
- Provisionamiento reproducible y consistente entre entornos.

## Fase 4 - Contenedores y plataforma de ejecucion
Objetivo:
Ejecutar la aplicacion con estandares de seguridad, recursos y disponibilidad.

Actividades:
- Definir imagenes optimizadas y estrategia de versionado.
- Definir despliegues por entorno.
- Definir escalado horizontal y politicas de disponibilidad.
- Definir controles de red y aislamiento.

Entregables:
- Estrategia de despliegue en contenedores.
- Politicas de recursos, salud y disponibilidad.
- Matriz de configuracion por entorno.

Criterio de exito:
- Despliegue estable, escalable y seguro.

## Fase 5 - CI/CD profesional
Objetivo:
Automatizar calidad, seguridad y entrega continua con puertas de control.

Actividades:
- Diseñar pipeline por etapas.
- Definir validaciones de calidad obligatorias.
- Definir escaneos de seguridad en pipeline.
- Definir estrategia de promociones entre entornos.

Entregables:
- Flujo CI/CD end-to-end.
- Politicas de aprobacion y bloqueo.
- Evidencia de trazabilidad de releases.

Criterio de exito:
- Ningun cambio llega a produccion sin pasar controles definidos.

## Fase 6 - Observabilidad y operacion
Objetivo:
Tener visibilidad completa de salud, rendimiento y fallos.

Actividades:
- Definir estandar de logs estructurados.
- Definir metricas de negocio y tecnicas.
- Definir trazabilidad distribuida.
- Definir alertas accionables con umbrales claros.

Entregables:
- Dashboards operativos.
- Catalogo de alertas priorizadas.
- Runbook de respuesta a incidentes.

Criterio de exito:
- Ante incidentes, existe deteccion temprana y respuesta guiada.

## Fase 7 - Seguridad integral
Objetivo:
Aplicar seguridad de identidad, red, datos y cadena de suministro.

Actividades:
- Definir estrategia de identidades administradas.
- Definir gestion segura de secretos y certificados.
- Definir hardening de plataforma y workloads.
- Definir controles de compliance basicos.

Entregables:
- Modelo de acceso y permisos por rol.
- Politicas de secretos y rotacion.
- Lista de controles de seguridad implementados.

Criterio de exito:
- Superficie de ataque reducida y controles auditables.

## Fase 8 - Resiliencia y continuidad
Objetivo:
Garantizar continuidad operativa ante fallos tecnicos y errores humanos.

Actividades:
- Definir objetivos de disponibilidad.
- Definir estrategia de backups y restauracion.
- Definir escenarios de desastre y recuperacion.
- Ejecutar simulacros de fallos.

Entregables:
- Plan de continuidad operativa.
- Procedimiento de restauracion validado.
- Informe de simulacro y mejoras.

Criterio de exito:
- Recuperacion validada dentro de objetivos definidos.

## Fase 9 - Gobierno, costos y madurez
Objetivo:
Controlar costos, estandarizar gobierno tecnico y medir madurez.

Actividades:
- Definir presupuesto y alertas de costo.
- Definir etiquetado para trazabilidad financiera.
- Definir tablero de madurez tecnica.
- Definir backlog de mejoras trimestral.

Entregables:
- Politica de costo y optimizacion.
- Indicadores de madurez del proyecto.
- Roadmap de mejora continua.

Criterio de exito:
- Costos predecibles y evolucion tecnica planificada.

## Fase 10 - Cierre profesional y portafolio
Objetivo:
Presentar el proyecto como caso real de ingenieria senior.

Actividades:
- Consolidar documentacion final.
- Preparar narrativa tecnica de decisiones clave.
- Preparar demo guiada de extremo a extremo.
- Preparar preguntas y respuestas de entrevista.

Entregables:
- Dossier tecnico final.
- Guion de demo profesional.
- Lista de lecciones aprendidas y siguientes pasos.

Criterio de exito:
- Capacidad de defender decisiones tecnicas con evidencia.

## 7. Definicion de terminado (DoD) por fase
Cada fase se considera terminada cuando:
- Existe evidencia verificable de lo implementado.
- Existe documentacion minima obligatoria.
- Existe validacion tecnica (calidad, seguridad, operacion).
- Existe decision de cierre y backlog residual identificado.

## 8. Riesgos clave y mitigacion
Riesgo:
Sobrecargar el proyecto con demasiadas funciones.
Mitigacion:
Proteger el alcance tecnico y priorizar profundidad sobre cantidad.

Riesgo:
Pipeline fragil o sin controles reales.
Mitigacion:
Definir puertas obligatorias y criterios de bloqueo claros.

Riesgo:
Complejidad sin documentacion.
Mitigacion:
Documentar cada decision relevante en el momento.

Riesgo:
Gastos cloud descontrolados.
Mitigacion:
Presupuesto, etiquetas y revisiones de costo periodicas.

## 9. Recomendaciones practicas (muy importantes)
- Mantener este repo senior separado del repo de aprendizaje basico.
- Trabajar por iteraciones cortas con entregables visibles.
- Priorizar calidad tecnica sobre volumen de features.
- Escribir decisiones tecnicas, no solo pasos.
- Medir todo: calidad, seguridad, costo, disponibilidad.
- Preparar el proyecto para ser auditado por terceros.
- Tratar cada fase como evidencia para entrevistas senior.

## 10. Orden de ejecucion recomendado
1. Fase 0 y Fase 1.
2. Fase 2 y Fase 3.
3. Fase 4 y Fase 5.
4. Fase 6 y Fase 7.
5. Fase 8 y Fase 9.
6. Fase 10.

## 11. Como usar esta guia
- Usar este documento como contrato tecnico del proyecto.
- Al iniciar una fase, copiar sus objetivos y entregables al backlog.
- Al cerrar una fase, registrar evidencia y lecciones aprendidas.
- No avanzar a la siguiente fase sin criterio de exito cumplido.

## 12. Estado del proyecto

### Fases completadas:
- ✅ Fase 0 - Definicion y diseño (COMPLETADA)
  - Documentos: fase0-decisiones-tecnicas.md, fase0-arquitectura-v1.md, fase0-resumen.md
  - Fecha de cierre: Completada

- ✅ Fase 1 - Fundacion del repositorio senior (COMPLETADA)
  - Documentos: fase1-estructura-repositorio.md, fase1-git-workflow.md, fase1-plantillas.md, fase1-archivos-configuracion.md, fase1-resumen.md
  - Estructura de carpetas creada
  - Git configurado con main y develop
  - Repositorio en GitHub: https://github.com/fabianbele2605/ops-incident-hub
  - Fecha de cierre: Completada

- ✅ Fase 2 - Arquitectura de Aplicacion (COMPLETADA)
  - Documentos: fase2-modelo-dominio.md, fase2-repositorios.md, fase2-casos-uso.md, fase2-handlers-http.md, fase2-resumen.md
  - Implementacion completa de Clean Architecture
  - Capa de dominio: Entidades (Incident, User), Errores tipados, Interfaces de repositorios
  - Capa de casos de uso: Create, Assign, List incidents
  - Capa de API: Handlers HTTP, DTOs, Router
  - Compilacion exitosa de todos los modulos
  - Fecha de cierre: 26 de febrero de 2025

- ✅ Fase 3 - Infraestructura como codigo avanzada (COMPLETADA)
  - Documentos: fase3-configuracion.md, fase3-repositorios-postgresql.md, fase3-migraciones.md, fase3-main-app.md, fase3-resumen.md
  - Sistema de configuracion con variables de entorno
  - Repositorios PostgreSQL: Incident y User con CRUD completo
  - Migraciones automaticas de base de datos
  - Main.go con inyeccion de dependencias
  - Docker Compose con PostgreSQL y API
  - Aplicacion funcional end-to-end
  - Fecha de cierre: 26 de febrero de 2025

- ✅ Fase 4 - Contenedores y plataforma de ejecucion (COMPLETADA)
  - Documentos: fase4-dockerfile-optimizado.md, fase4-health-checks.md, fase4-resumen.md
  - Dockerfile optimizado con multi-stage build (scratch)
  - Reduccion de imagen de ~300MB a ~10-15MB (95% reduccion)
  - Health checks en 3 niveles: /health, /health/live, /health/ready
  - Resource limits y restart policies configurados
  - Health check integrado en Dockerfile y docker-compose
  - Comando health para verificacion automatica
  - Aplicacion lista para produccion con auto-recovery
  - Fecha de cierre: 26 de febrero de 2025

- ✅ Fase 5 - CI/CD profesional (COMPLETADA)
  - Documentos: fase5-tests-unitarios.md, fase5-tests-completos.md, fase5-resumen.md
  - Tests unitarios de domain: Incident y User entities ✅
  - Tests unitarios de use cases: Create, Assign, List ✅
  - Tests de integración de repositorios con PostgreSQL ✅
  - Tests de handlers HTTP con mocks ✅
  - Security scanning automatizado (gosec, govulncheck, gitleaks) ✅
  - Mocks compartidos refactorizados ✅
  - Correcciones de linting (9 errores resueltos) ✅
  - CI/CD con GitHub Actions funcionando ✅
  - PRs #12-#19 mergeados exitosamente
  - Métricas: 51 casos de prueba, 887+ líneas de test, 3 security scanners
  - Pipeline con 5 jobs: lint, test, security, build, docker-build
  - Fecha de cierre: 27 de febrero de 2025

- ✅ Fase 6 - Observabilidad y operacion (COMPLETADA)
  - Documentos: fase6-observabilidad-diseno.md, fase6-resumen.md
  - Structured logging con slog (Go standard library) ✅
  - Request ID único por request con propagación en contexto ✅
  - Middleware HTTP para logging automático ✅
  - Métricas Prometheus (negocio y técnicas) ✅
  - Endpoint /metrics para scraping ✅
  - Middleware de métricas automático ✅
  - Trazabilidad end-to-end con request_id ✅
  - PRs #21-#22 mergeados exitosamente
  - Métricas: +1,657 líneas, 2 PRs, logging + metrics completos
  - Stack: slog, Prometheus, promhttp
  - Fecha de cierre: 27 de febrero de 2025

### Fase actual:
- 🔄 Fase 7 - Seguridad integral (PENDIENTE)
  - Definir estrategia de identidades administradas
  - Definir gestión segura de secretos y certificados
  - Definir hardening de plataforma y workloads
  - Definir controles de compliance básicos

### Proxima accion inmediata:
Iniciar Fase 7 - Seguridad Integral:
- Diseñar estrategia de gestión de secretos
- Evaluar herramientas (HashiCorp Vault, Azure Key Vault)
- Planificar hardening de contenedores
- Definir controles de seguridad

## 13. Decisiones iniciales de stack (acordadas)
- Backend principal: Go.
- Frontend: TypeScript.
- Base de datos relacional principal: PostgreSQL.
- Regla de complejidad: maximo 2 lenguajes principales en este repo.

## 14. Definicion oficial de la aplicacion
Nombre del proyecto:
- Ops Incident Hub.

Descripcion funcional:
- Plataforma para gestion de incidentes operativos y seguimiento de su ciclo de vida.

Objetivo funcional:
- Registrar, priorizar, asignar, monitorear y cerrar incidentes con trazabilidad completa.

Capacidades base esperadas:
- Gestion de incidentes con estados y prioridades.
- Asignacion de responsables y tiempos objetivo.
- Historial/auditoria de cambios.
- Dashboard operativo con metricas clave.
- Alertas e integracion con flujo de operacion.
