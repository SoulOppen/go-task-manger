---
name: go-cli-command-workflow
description: Implementa y modifica comandos Cobra en este repositorio Go CLI, con enfoque en separacion cmd/internal, flags claros y pruebas de comportamiento. Usar cuando se agregan o cambian comandos en cmd/.
---

# Go CLI Command Workflow

## Objetivo

Estandarizar cambios en comandos de `cmd/` para mantener el CLI coherente y facil de probar.

## Flujo recomendado

1. Identificar comando afectado en `cmd/`.
2. Mantener el handler del comando liviano.
3. Mover logica no trivial a `internal/`.
4. Verificar flags, ayudas y salida de texto.
5. Actualizar `README.md` cuando cambie la UX CLI.

## Reglas

- Evitar logica de negocio compleja dentro de `Run`.
- Nombrar flags y comandos de forma consistente.
- Reusar funciones de `internal/` en lugar de duplicar logica.

## Pruebas minimas obligatorias

- [ ] Prueba de ejecucion del comando principal afectado.
- [ ] Prueba de flags nuevos/modificados.
- [ ] Prueba de salida esperada (stdout/stderr) en casos relevantes.
- [ ] `go test ./...` ejecutado.

## Resultado esperado

El comando funciona con sus flags, mantiene convenciones del repo y queda cubierto por pruebas.
