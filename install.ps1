# NEXUS-V Windows Install Script
Write-Host "Installing NEXUS-V..." -ForegroundColor Cyan

$destDir = "$env:USERPROFILE\bin"
if (-not (Test-Path $destDir)) {
    New-Item -Path $destDir -ItemType Directory
}

$url = "https://github.com/billy-kidd-dev/nexus-v/releases/latest/download/nexus-v_windows_amd64.zip"
$zipFile = "$destDir\nexus-v.zip"
Invoke-WebRequest -Uri $url -OutFile $zipFile
Expand-Archive -Path $zipFile -DestinationPath $destDir -Force
Remove-Item $zipFile

# Add to PATH if not already there
$path = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($path -notlike "*$destDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$path;$destDir", "User")
    $env:PATH += ";$destDir"
    Write-Host "Added $destDir to User PATH." -ForegroundColor Yellow
}

Write-Host "NEXUS-V installed successfully! Restart your terminal and run 'nexus-v version' to verify." -ForegroundColor Green
