📝 Go Task Manager (GTM)
========================

**Go Task Manager** es una herramienta de línea de comandos (CLI) ligera y rápida diseñada para ayudarte a organizar tu flujo de trabajo diario sin salir de la terminal.

🚀 Características principales
------------------------------

*   **Gestión Rápida:** Añade, lista y elimina tareas en segundos.
    
*   **Estados:** Marca tareas como pendientes o completadas.
    
*   **Persistencia de Datos:** Almacenamiento local mediante archivos JSON o bases de datos embebidas (como BoltDB).
    
*   **Arquitectura Limpia:** Código organizado siguiendo las mejores prácticas de Go.
    

🛠️ Instalación
---------------

### Requisitos previos

*   **Go** (versión 1.21 o superior recomendada)
    

### Pasos para instalar

1.  Clona el repositorio:
    `git clone https://github.com/tu-usuario/go-task-manager.git`
2.  Entra al proyecto:
    `cd go-task-manager`
3.  Define `NAME` en `.env` (puedes copiar `.envExample`).
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

Una vez compilado, puedes usar los siguientes comandos:

**ComandoDescripciónEjemplo**addAñade una nueva tarea./gtm add "Estudiar Go"listMuestra todas las tareas./gtm listdoMarca una tarea como completada./gtm do 1rmElimina una tarea de la lista./gtm rm 2

Si usaste los scripts de instalacion, el binario queda con el nombre de `NAME` en tu `.env` (espacios reemplazados por `-`).

- Linux/macOS: `./bin/<NAME-normalizado>`
- Windows: `.\bin\<NAME-normalizado>.exe`

🔧 Stack Tecnológico
--------------------

*   **Lenguaje:** [Go](https://go.dev/)
    
*   **Librería CLI:** [Cobra](https://github.com/spf13/cobra) (opcional, para subcomandos complejos)
    
*   **Persistencia:** JSON / SQLite (dependiendo de tu implementación)
    

🤝 Contribuciones
-----------------

Si quieres contribuir:

1.  Haz un **Fork** del proyecto.
    
2.  Crea una rama con tu mejora: git checkout -b feature/nueva-mejora.
    
3.  Haz un **Commit**: git commit -m 'Añadida funcionalidad X'.
    
4.  Haz **Push**: git push origin feature/nueva-mejora.
    
5.  Abre un **Pull Request**.