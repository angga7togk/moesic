Write-Host "Installing Moesic..."

# $HOME\AppData\Local\Programs\moesic
$InstallDir = "$env:LOCALAPPDATA\Programs\moesic"
$Target = "$InstallDir\moesic.exe"

if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

$ReleaseInfo = Invoke-RestMethod https://api.github.com/repos/angga7togk/moesic/releases/latest
$LatestTag = $ReleaseInfo.tag_name
$DownloadUrl = "https://github.com/angga7togk/moesic/releases/download/$LatestTag/moesic-windows.exe"

$TempPath = "$env:TEMP\moesic.exe"
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempPath

Move-Item -Path $TempPath -Destination $Target -Force

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "PATH updated. You may need to restart your terminal."
} else {
    Write-Host "ℹPATH already includes install directory."
}

Write-Host "✅ Moesic installed successfully!"
Write-Host "You can now run: moesic"
