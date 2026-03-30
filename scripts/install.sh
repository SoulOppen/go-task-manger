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

if [[ ! -f ".env" ]]; then
  echo "Error: no existe .env. Crea .env con la variable NAME."
  exit 1
fi

RAW_NAME="$(awk -F'=' '/^NAME=/{print $2; exit}' .env | tr -d '\r' | sed 's/^ *//;s/ *$//;s/^"//;s/"$//')"
if [[ -z "${RAW_NAME}" ]]; then
  echo "Error: NAME no esta definido en .env."
  exit 1
fi

APP_NAME="${RAW_NAME// /-}"

echo "==> Descargando dependencias..."
go mod tidy

echo "==> Compilando binario..."
mkdir -p bin
go build -o "bin/${APP_NAME}" main.go

echo "==> Instalacion finalizada."
echo "Binario generado en: ${PROJECT_ROOT}/bin/${APP_NAME}"
echo "Ejemplo de uso: ./bin/${APP_NAME} version"
