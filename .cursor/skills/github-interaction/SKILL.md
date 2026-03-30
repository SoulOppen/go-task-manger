---
name: github-interaction
description: Gestiona flujo de trabajo con GitHub CLI para issues, ramas, commits, pull requests y revisiones. Usar cuando el usuario pida interaccion con GitHub, crear PR, revisar PR, consultar issue o ejecutar acciones con gh.
---

# GitHub Interaction

## Objetivo

Estandarizar la interaccion con GitHub en este repositorio usando `git` y `gh` de forma segura y repetible.

## Flujo recomendado

1. Revisar estado local:
   - `git status`
   - `git diff`
   - `git log --oneline -n 10`
2. Verificar rama activa y sincronizacion con remoto.
3. Antes de PR, correr pruebas del proyecto:
   - `go test ./...`
4. Crear o actualizar rama de trabajo con nombre claro.
5. Crear PR con resumen y plan de pruebas.

## Reglas operativas

- No usar comandos destructivos (`reset --hard`, `push --force`) salvo instruccion explicita.
- No hacer commit de secretos (`.env`, tokens, llaves).
- No hacer commit vacio.
- Al abrir PR, incluir siempre:
  - contexto del cambio,
  - alcance,
  - resultado de pruebas,
  - riesgos o pendientes.

## Plantilla breve para PR

Usar esta estructura en el cuerpo:

```markdown
## Resumen
- Cambio principal 1
- Cambio principal 2

## Pruebas
- [ ] go test ./...
- [ ] Prueba manual de CLI (comandos modificados)

## Riesgos
- Ninguno / describir riesgo
```

## Checklist obligatorio

- [ ] Cambios entendidos en `git diff`
- [ ] Pruebas ejecutadas (`go test ./...`)
- [ ] Commit con mensaje claro
- [ ] PR con resumen y test plan
