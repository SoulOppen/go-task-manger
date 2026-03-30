Param()

$ErrorActionPreference = "Stop"

$appName = "gtm"

Write-Host "==> Verificando entorno..."

$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
  Write-Error "Go no esta instalado. Instala Go 1.21+ y vuelve a ejecutar."
  exit 1
}

$goVersion = go version
Write-Host "Go detectado: $goVersion"

Write-Host "==> Verificando cliente MySQL/MariaDB..."
$mysql = Get-Command mysql -ErrorAction SilentlyContinue
$mariadb = Get-Command mariadb -ErrorAction SilentlyContinue
if ($mysql) {
  mysql --version
} elseif ($mariadb) {
  mariadb --version
} else {
  Write-Error "No se encontro mysql ni mariadb en PATH. Instala el cliente MySQL o MariaDB."
  exit 1
}

$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "==> Descargando dependencias..."
go mod tidy

Write-Host "==> Compilando binario..."
$binDir = Join-Path $projectRoot "bin"
if (-not (Test-Path $binDir)) {
  New-Item -Path $binDir -ItemType Directory | Out-Null
}

$binaryPath = Join-Path $binDir "$appName.exe"
go build -o $binaryPath main.go

Write-Host "==> Verificando binario (version)..."
$verOut = & $binaryPath version 2>&1
if ($LASTEXITCODE -ne 0) {
  Write-Error "El comando version fallo (codigo $LASTEXITCODE)."
  exit 1
}
if ([string]::IsNullOrWhiteSpace([string]$verOut)) {
  Write-Error "El comando version no produjo salida."
  exit 1
}
Write-Host "Salida: $verOut"

Write-Host "==> Instalacion finalizada."
Write-Host "Binario generado en: $binaryPath"
Write-Host "Ejemplo de uso: .\bin\$appName.exe version"
