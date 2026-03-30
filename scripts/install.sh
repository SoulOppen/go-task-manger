#!/usr/bin/env bash
set -euo pipefail

readonly APP_NAME="gtm"

echo "==> Verificando entorno..."

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go no esta instalado. Instala Go 1.21+ y vuelve a ejecutar."
  exit 1
fi

GO_VERSION="$(go version | awk '{print $3}')"
echo "Go detectado: ${GO_VERSION}"

echo "==> Verificando cliente MySQL/MariaDB..."
if command -v mysql >/dev/null 2>&1; then
  mysql --version
elif command -v mariadb >/dev/null 2>&1; then
  mariadb --version
else
  echo "Error: no se encontro mysql ni mariadb en PATH. Instala el cliente MySQL o MariaDB."
  exit 1
fi

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

echo "==> Descargando dependencias..."
go mod tidy

echo "==> Compilando binario..."
mkdir -p bin
go build -o "bin/${APP_NAME}" main.go

BIN_PATH="${PROJECT_ROOT}/bin/${APP_NAME}"
echo "==> Verificando binario (version)..."
if ! OUT="$("$BIN_PATH" version 2>&1)"; then
  echo "Error: \"${BIN_PATH} version\" fallo."
  exit 1
fi
if [[ -z "${OUT//[$'\t\r\n ']}" ]]; then
  echo "Error: \"version\" no produjo salida."
  exit 1
fi
echo "Salida: ${OUT}"

echo "==> Instalacion finalizada."
echo "Binario generado en: ${BIN_PATH}"
echo "Ejemplo de uso: ./bin/${APP_NAME} version"
