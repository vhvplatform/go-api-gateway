@echo off
REM Build script for Windows (CMD/Batch)
REM API Gateway Service Build Script

setlocal enabledelayedexpansion

REM Variables
set APP_NAME=api-gateway
set GOBIN=%CD%\bin
set GOCMD=go

REM Parse command line arguments
if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="build" goto build
if "%1"=="test" goto test
if "%1"=="clean" goto clean
if "%1"=="deps" goto deps
if "%1"=="run" goto run
if "%1"=="fmt" goto fmt
if "%1"=="vet" goto vet
if "%1"=="validate" goto validate
if "%1"=="version" goto version
goto help

:help
echo Available commands:
echo   build      - Build the application
echo   test       - Run tests
echo   clean      - Clean build artifacts
echo   deps       - Download dependencies
echo   run        - Run the application locally
echo   fmt        - Format Go code
echo   vet        - Run go vet
echo   validate   - Run all validation checks (fmt, vet, test)
echo   version    - Display version information
echo   help       - Display this help screen
echo.
echo Usage: build.bat [command]
echo Example: build.bat build
goto end

:build
echo Building %APP_NAME% for Windows...
if not exist "%GOBIN%" mkdir "%GOBIN%"

REM Get version info
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "!VERSION!"=="" set VERSION=dev

REM Get build time
for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss'"') do set BUILD_TIME=%%i

echo Building version: !VERSION!
%GOCMD% build -ldflags "-X main.Version=!VERSION! -X main.BuildTime=!BUILD_TIME!" -o "%GOBIN%\%APP_NAME%.exe" .\cmd\main.go
if %ERRORLEVEL% EQU 0 (
    echo Build complete: %GOBIN%\%APP_NAME%.exe
    exit /b 0
) else (
    echo Build failed!
    exit /b 1
)

:test
echo Running tests...
%GOCMD% test -v -race -coverprofile=coverage.txt -covermode=atomic .\...
if %ERRORLEVEL% EQU 0 (
    echo Tests complete
    exit /b 0
) else (
    echo Tests failed!
    exit /b 1
)

:clean
echo Cleaning build artifacts...
if exist "%GOBIN%" rd /s /q "%GOBIN%"
if exist coverage.txt del /f coverage.txt
if exist coverage.html del /f coverage.html
if exist test-output.log del /f test-output.log
echo Clean complete
goto end

:deps
echo Downloading dependencies...
%GOCMD% mod download
if %ERRORLEVEL% EQU 0 (
    echo Dependencies downloaded
    exit /b 0
) else (
    echo Failed to download dependencies!
    exit /b 1
)

:run
echo Running %APP_NAME%...
%GOCMD% run .\cmd\main.go
goto end

:fmt
echo Formatting code...
%GOCMD% fmt .\...
echo Format complete
goto end

:vet
echo Running go vet...
%GOCMD% vet .\...
if %ERRORLEVEL% EQU 0 (
    echo Vet complete
    exit /b 0
) else (
    echo Vet found issues!
    exit /b 1
)

:validate
echo Running validation checks...
call :fmt
if %ERRORLEVEL% NEQ 0 exit /b 1
call :vet
if %ERRORLEVEL% NEQ 0 exit /b 1
call :test
if %ERRORLEVEL% EQU 0 (
    echo All validation checks passed!
    exit /b 0
) else (
    exit /b 1
)

:version
echo Application: %APP_NAME%
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "!VERSION!"=="" set VERSION=dev
echo Version: !VERSION!
for /f "tokens=*" %%i in ('%GOCMD% version') do echo Go Version: %%i
goto end

:end
endlocal
