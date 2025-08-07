@echo off
setlocal enabledelayedexpansion

echo Installing Moesic...

:: Set install directory
set "INSTALL_DIR=%LOCALAPPDATA%\Programs\moesic"
set "TARGET=%INSTALL_DIR%\moesic.exe"

:: Create directory if not exists
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

:: Get latest release tag from GitHub API
for /f "tokens=*" %%i in ('curl -s https://api.github.com/repos/angga7togk/moesic/releases/latest ^| findstr /i "tag_name"') do (
    set "TAG_LINE=%%i"
)

:: Extract tag name
for /f "tokens=2 delims=:" %%a in ("!TAG_LINE!") do (
    set "TAG_NAME=%%a"
)
set "TAG_NAME=%TAG_NAME:~2,-2%"

:: Compose download URL
set "DOWNLOAD_URL=https://github.com/angga7togk/moesic/releases/download/%TAG_NAME%/moesic-windows.exe"

:: Download file
echo Downloading Moesic from %DOWNLOAD_URL%...
curl -L %DOWNLOAD_URL% -o "%INSTALL_DIR%\moesic.exe"
if %ERRORLEVEL% NEQ 0 (
    echo Failed to download Moesic.
    exit /b 1
)

:: Add install dir to PATH (if not already there)
set "REG_KEY=HKCU\Environment"
for /f "tokens=3*" %%a in ('reg query "%REG_KEY%" /v PATH 2^>nul ^| find "REG_SZ"') do (
    set "USER_PATH=%%a %%b"
)

echo %USER_PATH% | find /i "%INSTALL_DIR%" >nul
if errorlevel 1 (
    echo Adding %INSTALL_DIR% to PATH...
    setx PATH "%USER_PATH%;%INSTALL_DIR%"
    echo You may need to restart your terminal or PC to use 'moesic' globally.
) else (
    echo PATH already contains Moesic install directory.
)

echo.
echo ✅ Moesic installed successfully!
echo You can now run: moesic

exit /b 0
