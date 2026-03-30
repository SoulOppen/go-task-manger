# AGENTS

## Contexto del proyecto

Este repositorio implementa `Go Task Manager (GTM)`, un CLI en Go orientado a gestion de tareas desde terminal.

Referencia funcional principal: `README.md`.

Estructura actual relevante:
- `cmd/`: comandos Cobra (`root`, `login`, `logout`, `version`).
- `internal/auth/`: flujo de login/signup y persistencia basica de usuarios.
- `internal/task/`: modelo `Task` y salida de tareas.
- `internal/name/`: delegado a `internal/config` para el nombre del CLI.

## Reglas de trabajo para agentes

1. Hacer cambios pequenos y enfocados por tarea.
2. Mantener consistencia con el estilo existente del repo.
3. No introducir dependencias nuevas salvo necesidad real.
4. Documentar cambios de UX/CLI en `README.md` si aplica.
5. No usar comandos destructivos de git sin autorizacion explicita.

## Guia rapida de implementacion

### Nuevo subcomando Cobra

1. Crear o actualizar archivo en `cmd/`.
2. Registrar comando en `rootCmd` o comando padre.
3. Agregar flags con nombres claros.
4. Enviar logica de negocio a `internal/` (evitar logica pesada en `cmd/`).
5. Ajustar documentacion de uso en `README.md` si cambia el flujo.

### Cambios en auth

1. Mantener separacion entre lectura/escritura y experiencia CLI.
2. Tratar errores de I/O de forma explicita.
3. Evitar exponer datos sensibles en logs o salida.
4. Validar entradas de usuario antes de persistir.

### Cambios en task

1. Conservar compatibilidad del modelo `Task` cuando sea posible.
2. Mantener salida CLI legible y estable.
3. Si cambia formato de salida, actualizar docs y pruebas.

## Politica obligatoria de pruebas

Sin pruebas, una tarea no se considera terminada.

Todo cambio debe incluir pruebas nuevas o ajustadas, segun corresponda:

- Cambios en `cmd/`:
  - pruebas de comandos, flags y salida esperada.
- Cambios en `internal/auth/`:
  - pruebas de flujo feliz y errores de I/O.
- Cambios en `internal/task/`:
  - pruebas de modelo y formato de salida.
- Cambios transversales:
  - cobertura de integracion minima del flujo afectado.

## Checklist de cierre

- [ ] Codigo compila sin errores.
- [ ] Pruebas relevantes agregadas/actualizadas.
- [ ] `go test ./...` ejecutado.
- [ ] Documentacion ajustada si cambia comportamiento.
- [ ] Diff final limpio y enfocado al alcance.
