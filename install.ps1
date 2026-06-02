# install.ps1 — Installs env-vault on Windows via PowerShell.
#
# Uso (one-liner):
#   irm https://raw.githubusercontent.com/U23205221-S/env-vault/main/install.ps1 | iex
#
# Variables de entorno opcionales:
#   $env:ENV_VAULT_VERSION = "v1.2.0"   # instala una versión específica
#
# Instala en %LOCALAPPDATA%\Programs\env-vault\ (no requiere permisos de
# administrador) y agrega ese directorio al PATH del usuario.

$ErrorActionPreference = 'Stop'

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

function Write-Info  { param($m) Write-Host "[env-vault] $m" }
function Write-Err   { param($m) Write-Host "[env-vault] ERROR: $m" -ForegroundColor Red; exit 1 }

# ---------------------------------------------------------------------------
# Detección de arquitectura
# ---------------------------------------------------------------------------

switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'ARM64' { $arch = 'arm64' }
    default { Write-Err "Arquitectura no soportada: $env:PROCESSOR_ARCHITECTURE. Soportadas: AMD64, ARM64." }
}

# ---------------------------------------------------------------------------
# Resolución de versión
# ---------------------------------------------------------------------------

$repo   = 'U23205221-S/env-vault'
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

$versionTag = $env:ENV_VAULT_VERSION
if (-not $versionTag) {
    try {
        $release    = Invoke-RestMethod -Uri $apiUrl -ErrorAction Stop
        $versionTag = $release.tag_name
    } catch {
        Write-Err "No se pudo resolver la última versión desde la API de GitHub. Verificá tu conexión a internet."
    }
}

# Aceptamos "v1.2.0" o "1.2.0"; normalizamos.
if ($versionTag -notmatch '^v') { $versionTag = "v$versionTag" }
$versionBare = $versionTag -replace '^v', ''

# ---------------------------------------------------------------------------
# Descarga del archivo y checksums
# ---------------------------------------------------------------------------

$baseUrl      = "https://github.com/$repo/releases/download/$versionTag"
$filename     = "env-vault_${versionBare}_windows_${arch}.zip"
$archiveUrl   = "$baseUrl/$filename"
$checksumUrl  = "$baseUrl/checksums.txt"

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "env-vault-install-$([System.Guid]::NewGuid())"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

try {
    Write-Info "Descargando $archiveUrl"
    $archivePath = Join-Path $tmpDir $filename
    Invoke-WebRequest -Uri $archiveUrl -OutFile $archivePath -ErrorAction Stop

    Write-Info "Descargando $checksumUrl"
    $checksumsPath = Join-Path $tmpDir 'checksums.txt'
    Invoke-WebRequest -Uri $checksumUrl -OutFile $checksumsPath -ErrorAction Stop

    # ---------------------------------------------------------------------
    # Verificación de SHA256
    # ---------------------------------------------------------------------

    $expected = $null
    Get-Content $checksumsPath | ForEach-Object {
        # GoReleaser genera lineas como "<hash>  filename" (dos espacios) o
        # "<hash> *filename" (modo --check). Aceptamos ambos.
        if ($_ -match "^\S+\s+\*?$([regex]::Escape($filename))$") {
            $script:expected = ($_ -split '\s+')[0]
        }
    }

    if (-not $expected) {
        Write-Err "No se encontró checksum para $filename en checksums.txt."
    }

    $actual = (Get-FileHash -Algorithm SHA256 -Path $archivePath).Hash.ToLower()
    if ($actual -ne $expected.ToLower()) {
        Write-Err "Checksum inválido. Esperado: $expected, actual: $actual."
    }
    Write-Info "SHA256 verificado correctamente."

    # ---------------------------------------------------------------------
    # Extracción
    # ---------------------------------------------------------------------

    Write-Info "Extrayendo..."
    $extractDir = Join-Path $tmpDir 'extracted'
    Expand-Archive -Path $archivePath -DestinationPath $extractDir -Force

    $binaryPath = Join-Path $extractDir 'env-vault.exe'
    if (-not (Test-Path $binaryPath)) {
        Write-Err "No se encontró env-vault.exe en el archivo extraído."
    }

    # ---------------------------------------------------------------------
    # Instalación
    # ---------------------------------------------------------------------

    # Directorio de instalación: %LOCALAPPDATA%\Programs\env-vault
    # (estándar de Windows, no requiere admin). Fallback raro a $HOME\bin.
    if ($env:LOCALAPPDATA) {
        $installDir = Join-Path $env:LOCALAPPDATA 'Programs\env-vault'
    } else {
        $installDir = Join-Path $env:USERPROFILE 'bin'
    }
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null

    Copy-Item -Path $binaryPath -Destination (Join-Path $installDir 'env-vault.exe') -Force
    Write-Info "Binario instalado en $installDir\env-vault.exe."

    # Agregar al PATH del usuario (no requiere admin)
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable('Path', "$userPath;$installDir", 'User')
        $env:Path = "$env:Path;$installDir"
        Write-Info "PATH del usuario actualizado."
    } else {
        Write-Info "$installDir ya estaba en el PATH del usuario."
    }

    Write-Info ""
    Write-Info "env-vault $versionTag instalado correctamente."
    Write-Info "Abrí una nueva terminal (PowerShell o cmd) y ejecutá 'env-vault --help'."
    Write-Info "Para actualizar: volvé a correr este mismo comando."

} finally {
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
}
