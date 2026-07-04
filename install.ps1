[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$ErrorActionPreference = "Stop"
$repoUrl = "https://raw.githubusercontent.com/Dxrmy/Wormhole/main"
$binaryName = "proxy-windows-amd64.exe"
$url = "$repoUrl/$binaryName"

$installDir = "$env:USERPROFILE\.wormhole"
if (-not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

$dest = "$installDir\wormhole.exe"

Write-Host "Downloading Terraria Proxy (Wormhole)..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $url -OutFile $dest

Write-Host "Starting Wormhole Proxy..." -ForegroundColor Green
Start-Process -FilePath $dest
