@echo off
if "%1"=="help" (
    echo Usage: make.bat [command]
    echo Commands:
    echo   tidy     Run go mod tidy
    echo   build    Build the project with mmap and traditional methods
    echo   help     Show this help message
    echo How to run program:
    echo   start ./bin/mmap or ./bin/traditional
    echo   ./bin/mmap -h or ./bin/traditional -h for help
    exit /b 0
)

if "%1"=="tidy" (
    echo Running go mod tidy...
    go mod tidy
) else if "%1"=="build" (
    echo Running go mod tidy...
    go mod tidy
    echo "Building the project with mmap method..."
    go build -o ./bin/mmap.exe -tags using_mmap_io
    echo "Building the project with traditional method..."
    go build -o ./bin/traditional.exe -tags using_traditional_io
) else (
    echo Invalid argument. Use "tidy", "build".
)