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

echo "==> Descargando dependencias..."
go mod tidy

echo "==> Compilando binario..."
mkdir -p bin
go build -o bin/gtm main.go

echo "==> Instalacion finalizada."
echo "Binario generado en: ${PROJECT_ROOT}/bin/gtm"
echo "Ejemplo de uso: ./bin/gtm version"
