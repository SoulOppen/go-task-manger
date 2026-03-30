---
name: go-auth-flow-safety
description: Guía cambios seguros en el flujo de autenticacion del CLI (signup/login), manejando entrada de usuario, persistencia en archivo y errores de I/O con pruebas obligatorias. Usar cuando se edite internal/auth/.
---

# Go Auth Flow Safety

## Objetivo

Mantener seguro y estable el flujo de autenticacion en `internal/auth/`.

## Flujo recomendado

1. Revisar rutas de lectura/escritura de usuarios.
2. Separar logica de persistencia de la interaccion por terminal.
3. Confirmar mensajes de error claros para usuario final.
4. Evitar dependencias de estado global no controlado.
5. Validar escenarios de usuario existente/no existente.

## Reglas

- Manejar errores de archivo de forma explicita.
- No registrar secretos ni datos sensibles.
- Mantener mensajes CLI consistentes y accionables.

## Pruebas minimas obligatorias

- [ ] Flujo de registro exitoso.
- [ ] Flujo con usuario ya existente.
- [ ] Flujo de login exitoso.
- [ ] Flujo con usuario inexistente.
- [ ] Caso de error de I/O (lectura o escritura).
- [ ] `go test ./...` ejecutado.

## Resultado esperado

Cambios en auth confiables, con cobertura de casos felices y de error.
