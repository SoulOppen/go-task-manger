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

Write-Host "==> Verificando cliente MySQL..."
$mysql = Get-Command mysql -ErrorAction SilentlyContinue
if (-not $mysql) {
  Write-Error "No se encontro mysql en PATH. Instala el cliente MySQL."
  exit 1
}
mysql --version

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

Write-Host "==> Instalacion finalizada."
Write-Host "Binario generado en: $binaryPath"
Write-Host ""
Write-Host "==> $appName version"
& $binaryPath version
if ($LASTEXITCODE -ne 0) {
  exit $LASTEXITCODE
}
