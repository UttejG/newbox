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

$Arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq 'Arm64') { 'arm64' } else { 'amd64' }
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
} catch {
    Write-Err "Download failed: $_"
}

# Extract
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Copy-Item (Join-Path $TmpDir "newbox.exe") (Join-Path $InstallDir "newbox.exe") -Force

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    Write-Info "Added $InstallDir to PATH (restart terminal to take effect)"
}

# Cleanup
Remove-Item $TmpDir -Recurse -Force

Write-Success "newbox $Version installed to $InstallDir\newbox.exe"
Write-Success "Run 'newbox' to get started!"
