Param()

$ErrorActionPreference = "Stop"

Write-Host "==> Verificando entorno..."

$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
  Write-Error "Go no esta instalado. Instala Go 1.21+ y vuelve a ejecutar."
  exit 1
}

$goVersion = go version
Write-Host "Go detectado: $goVersion"

$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "==> Descargando dependencias..."
go mod tidy

Write-Host "==> Compilando binario..."
$binDir = Join-Path $projectRoot "bin"
if (-not (Test-Path $binDir)) {
  New-Item -Path $binDir -ItemType Directory | Out-Null
}

go build -o (Join-Path $binDir "gtm.exe") main.go

Write-Host "==> Instalacion finalizada."
Write-Host "Binario generado en: $(Join-Path $projectRoot 'bin\gtm.exe')"
Write-Host "Ejemplo de uso: .\bin\gtm.exe version"
