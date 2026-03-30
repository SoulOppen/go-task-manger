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

$envFile = Join-Path $projectRoot ".env"
if (-not (Test-Path $envFile)) {
  Write-Error "No existe .env. Crea .env con la variable NAME."
  exit 1
}

$nameLine = Select-String -Path $envFile -Pattern '^NAME=' | Select-Object -First 1
if (-not $nameLine) {
  Write-Error "NAME no esta definido en .env."
  exit 1
}

$rawName = ($nameLine.Line -replace '^NAME=', '').Trim().Trim('"')
if ([string]::IsNullOrWhiteSpace($rawName)) {
  Write-Error "NAME no puede estar vacio en .env."
  exit 1
}

$appName = $rawName -replace ' ', '-'

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
