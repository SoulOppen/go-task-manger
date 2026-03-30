#!/usr/bin/env bash
set -euo pipefail

echo "==> Verificando entorno..."

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go no esta instalado. Instala Go 1.21+ y vuelve a ejecutar."
  exit 1
fi

GO_VERSION="$(go version | awk '{print $3}')"
echo "Go detectado: ${GO_VERSION}"

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

CONFIG_FILE="${PROJECT_ROOT}/internal/config/config.go"
if [[ ! -f "${CONFIG_FILE}" ]]; then
  echo "Error: no se encontro ${CONFIG_FILE}"
  exit 1
fi

APP_NAME="$(grep -E 'DefaultName[[:space:]]*=' "${CONFIG_FILE}" | head -1 | sed -E 's/.*DefaultName[[:space:]]*=[[:space:]]*"([^"]+)".*/\1/' | tr ' ' '-')"
if [[ -z "${APP_NAME}" ]]; then
  echo "Error: no se pudo leer DefaultName de internal/config/config.go"
  exit 1
fi

echo "==> Descargando dependencias..."
go mod tidy

echo "==> Compilando binario..."
mkdir -p bin
go build -o "bin/${APP_NAME}" main.go

echo "==> Instalacion finalizada."
echo "Binario generado en: ${PROJECT_ROOT}/bin/${APP_NAME}"
echo "Ejemplo de uso: ./bin/${APP_NAME} version"
