---
name: go-task-model-evolution
description: Estandariza cambios del modelo Task y su salida de terminal para evitar regresiones de formato o estructura, exigiendo pruebas de modelo y renderizado. Usar cuando se modifique internal/task/.
---

# Go Task Model Evolution

## Objetivo

Permitir evolucionar `internal/task/` sin romper el comportamiento esperado en CLI.

## Flujo recomendado

1. Identificar el cambio de estructura en `Task` (campos, tags, reglas).
2. Evaluar impacto en funciones de salida y serializacion.
3. Mantener formato de salida estable, salvo cambio intencional.
4. Si cambia formato, documentar en `README.md` y reflejar en pruebas.

## Reglas

- Evitar cambios ambiguos en nombres de campos.
- Mantener compatibilidad hacia atras cuando sea posible.
- Aislar el formateo de salida para facilitar pruebas.

## Pruebas minimas obligatorias

- [ ] Prueba de creacion/uso del modelo `Task`.
- [ ] Prueba de serializacion o estructura (si aplica).
- [ ] Prueba del formato de salida de tareas.
- [ ] Prueba para escenario con y sin `Project`.
- [ ] `go test ./...` ejecutado.

## Resultado esperado

Evolucion del modelo con cambios controlados, documentados y validados con pruebas.
