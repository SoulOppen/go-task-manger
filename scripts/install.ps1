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

$configFile = Join-Path $projectRoot "internal\config\config.go"
if (-not (Test-Path $configFile)) {
  Write-Error "No se encontro $configFile"
  exit 1
}

$line = Select-String -Path $configFile -Pattern '^\s*DefaultName\s*=' | Select-Object -First 1
if (-not $line) {
  Write-Error "No se pudo leer DefaultName de internal/config/config.go"
  exit 1
}

if ($line.Line -notmatch 'DefaultName\s*=\s*"([^"]+)"') {
  Write-Error "Formato inesperado de DefaultName en config.go"
  exit 1
}

$appName = $Matches[1] -replace ' ', '-'

Write-Host "==> Descargando dependencias..."
go mod tidy

Write-Host "==> Compilando binario..."
$binDir = Join-Path $projectRoot "bin"
if (-not (Test-Path $binDir)) {
  New-Item -Path $binDir -ItemType Directory | Out-Null
}

$binaryName = "$appName.exe"
go build -o (Join-Path $binDir $binaryName) main.go

Write-Host "==> Instalacion finalizada."
Write-Host "Binario generado en: $(Join-Path $binDir $binaryName)"
Write-Host "Ejemplo de uso: .\bin\$binaryName version"
