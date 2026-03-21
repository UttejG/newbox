#Requires -Version 5.1
[CmdletBinding()]
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\newbox"
)

$ErrorActionPreference = 'Stop'
$Repo = "UttejG/newbox"

function Write-Info    { Write-Host "[newbox] $args" -ForegroundColor Blue }
function Write-Success { Write-Host "[newbox] $args" -ForegroundColor Green }
function Write-Err     { Write-Host "[newbox] $args" -ForegroundColor Red; exit 1 }

# Get latest version
Write-Info "Fetching latest release..."
try {
    $Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name
} catch {
    Write-Err "Failed to fetch release info: $_"
}

$RawArch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
# Windows ARM64 is not yet supported natively; fall back to amd64 (runs via x64 emulation on Windows 11).
$Arch = if ($RawArch -eq 'Arm64') {
    Write-Info "Windows ARM64 detected; using amd64 build (runs via x64 emulation)"
    'amd64'
} else {
    'amd64'
}
$VersionNum = $Version.TrimStart('v')
$FileName = "newbox_${VersionNum}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$FileName"

Write-Info "Installing newbox $Version (windows/$Arch)..."

# Download
$TmpDir = Join-Path $env:TEMP "newbox-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir | Out-Null
$ZipPath = Join-Path $TmpDir $FileName

try {
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing
    $ChecksumsUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"
    $ChecksumsPath = Join-Path $TmpDir "checksums.txt"
    Invoke-WebRequest -Uri $ChecksumsUrl -OutFile $ChecksumsPath -UseBasicParsing
} catch {
    Write-Err "Download failed: $_"
}

# Verify checksum
Write-Info "Verifying checksum..."
$Expected = (Get-Content $ChecksumsPath | Where-Object { $_ -match [regex]::Escape($FileName) }) -split '\s+' | Select-Object -First 1
$Actual = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
if ($Actual -ne $Expected.ToLower()) {
    Write-Err "Checksum mismatch! Expected $Expected, got $Actual"
}

# Extract
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Copy-Item (Join-Path $TmpDir "newbox.exe") (Join-Path $InstallDir "newbox.exe") -Force

# Add to PATH if not already there (exact entry match to avoid false positives)
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$PathArray = $UserPath -split ';' | Where-Object { $_ -ne '' }
if (-not ($PathArray -contains $InstallDir)) {
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    Write-Info "Added $InstallDir to PATH (restart terminal to take effect)"
}

# Cleanup
Remove-Item $TmpDir -Recurse -Force

Write-Success "newbox $Version installed to $InstallDir\newbox.exe"
Write-Success "Run 'newbox' to get started!"
