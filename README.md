📝 Go Task Manager (GTM)
========================

**Go Task Manager** es una herramienta de línea de comandos (CLI) ligera y rápida diseñada para ayudarte a organizar tu flujo de trabajo diario sin salir de la terminal.

🚀 Características principales
------------------------------

*   **Gestión Rápida:** Añade, lista y elimina tareas en segundos.
    
*   **Estados:** Marca tareas como pendientes o completadas.
    
*   **Persistencia de Datos:** Tareas y **usuarios** (credenciales + quick-connect) en **MySQL**; sesion local (`session.json`) y archivos `quick_connect_*.json` en el directorio de configuracion del usuario.
    
*   **Arquitectura Limpia:** Código organizado siguiendo las mejores prácticas de Go.
    

🛠️ Instalación
---------------

### Requisitos previos

*   **Go** (versión 1.21 o superior recomendada)
*   **Cliente MySQL** en PATH (`mysql --version`; los scripts `install.sh` / `install.ps1` no validan que el servidor este en marcha)
    

### Pasos para instalar

1.  Clona el repositorio:
    `git clone https://github.com/tu-usuario/go-task-manager.git`
2.  Entra al proyecto:
    `cd go-task-manager`
3.  Configura MySQL en `.env` a partir de `.envExample` (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`).
4.  Usa el instalador segun tu sistema operativo:
    - Linux/macOS: `bash scripts/install.sh`
    - Windows (PowerShell): `powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1`
    

📂 Estructura del Proyecto
--------------------------

El proyecto sigue una estructura modular para facilitar su mantenimiento:

*   cmd/: Puntos de entrada de la aplicación.
    
*   internal/: Lógica de negocio (gestión de tareas y almacenamiento).
    
*   pkg/: Paquetes de utilidad compartidos.
    

💻 Uso
------

Los instaladores generan siempre el binario **`gtm`** (`gtm.exe` en Windows), alineado con `DefaultName` en [internal/config/config.go](internal/config/config.go). Tambien puedes usar `go run .`.

### Autenticacion (MySQL)

Los comandos `login`, `login --signup` y `switch` requieren las mismas variables `DB_*` que las tareas. La primera ejecucion crea la tabla `users` si no existe. `logout` solo borra la sesion local.

- `login` / `login --signup` / `logout` / `switch`

### Tareas (MySQL)

La primera ejecucion de un subcomando `task` crea la tabla `tasks` si no existe.

| Comando | Descripcion |
|--------|-------------|
| `task add` (sin flags) | Modo interactivo por terminal. |
| `task add --name "..." --description "..." [--relevance N] [--due YYYY-MM-DD]` | Crea tarea (id UUID); imprime el id. |
| `task list` | Lista ordenada: relevancia mayor primero; con fecha de entrega antes que sin fecha; entrega mas cercana primero. |
| `task get [id]` | Detalle; sin `id` lo pide por terminal. |
| `task update [id]` | Sin flags: modo interactivo. Con flags: actualiza solo lo indicado (`--clear-due` quita entrega). |
| `task delete [id]` | Elimina; sin `id` lo pide por terminal. |

Ejemplo:

```bash
./bin/gtm task add --name "Reunion" --description "Cliente X" --relevance 8 --due 2026-04-15
./bin/gtm task list
```

Tras instalar con el script:

- Linux/macOS: `./bin/gtm`
- Windows: `.\bin\gtm.exe`

🔧 Stack Tecnológico
--------------------

*   **Lenguaje:** [Go](https://go.dev/)
    
*   **Librería CLI:** [Cobra](https://github.com/spf13/cobra) (opcional, para subcomandos complejos)
    
*   **Persistencia:** MySQL (`tasks`, `users`); JSON local (sesion y quick_connect por archivo)
    

🤝 Contribuciones
-----------------

Si quieres contribuir:

1.  Haz un **Fork** del proyecto.
    
2.  Crea una rama con tu mejora: git checkout -b feature/nueva-mejora.
    
3.  Haz un **Commit**: git commit -m 'Añadida funcionalidad X'.
    
4.  Haz **Push**: git push origin feature/nueva-mejora.
    
5.  Abre un **Pull Request**.