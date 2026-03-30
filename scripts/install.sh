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

echo "==> Verificando cliente MySQL..."
if ! command -v mysql >/dev/null 2>&1; then
  echo "Error: no se encontro mysql en PATH. Instala el cliente MySQL."
  exit 1
fi
mysql --version

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

echo "==> Descargando dependencias..."
go mod tidy

echo "==> Compilando binario..."
mkdir -p bin
go build -o "bin/${APP_NAME}" main.go

BIN_PATH="${PROJECT_ROOT}/bin/${APP_NAME}"
INSTALL_BIN="${HOME}/.local/bin"
mkdir -p "${INSTALL_BIN}"
cp "${BIN_PATH}" "${INSTALL_BIN}/${APP_NAME}"
chmod +x "${INSTALL_BIN}/${APP_NAME}"

echo "==> Instalacion finalizada."
echo "Binario en el repo: ${BIN_PATH}"
echo "Copia para uso global: ${INSTALL_BIN}/${APP_NAME}"
if [[ ":${PATH}:" != *":${INSTALL_BIN}:"* ]]; then
  echo ""
  echo ">>> Anade esta linea a ~/.bashrc, ~/.zshrc o ~/.profile y abre una terminal nueva:"
  echo "    export PATH=\"${INSTALL_BIN}:\${PATH}\""
  echo ">>> O ejecuta con ruta completa hasta entonces:"
  echo "    ${INSTALL_BIN}/${APP_NAME} -v"
fi
echo ""
echo "==> ${APP_NAME} -v"
"${INSTALL_BIN}/${APP_NAME}" -v
