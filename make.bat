@echo off
if "%1"=="help" (
    echo Usage: make.bat [command]
    echo Commands:
    echo   help     Show this help message
    echo   tidy     Run go mod tidy
    echo   build    Build the project
    echo   test-io  Build the project with mmap and traditional methods for testing I/O performance
    echo   fmt      Format the code and imports
    echo How to run program:
    echo   start ./bin/main.exe
    echo   run ./bin/main.exe -h # for help
    exit /b 0
)

if "%1"=="tidy" (
    echo Running go mod tidy...
    go mod tidy
) else if "%1"=="build" (
    echo Running go mod tidy and formatting the code...
    call "%~f0" tidy
    call "%~f0" fmt
    echo Building the project...
    go build -o ./bin/main.exe -tags using_mmap_io
) else if "%1"=="test-io" (
    echo Running go mod tidy and formatting the code...
    call "%~f0" tidy
    call "%~f0" fmt
    echo "Building the project with mmap method..."
    go build -o ./bin/mmap.exe -tags using_mmap_io
    echo "Building the project with traditional method..."
    go build -o ./bin/traditional.exe -tags using_traditional_io
) else if "%1"=="fmt" (
    echo Checking if goimports is installed...
    where goimports >nul 2>nul || (
        echo goimports is not installed. Please install it first: go install golang.org/x/tools/cmd/goimports@latest.
        exit /b 1
    )
    echo goimports is installed.
    echo Checking if gofumpt is installed...
    where gofumpt >nul 2>nul || (
        echo gofumpt is not installed. Please install it first: go install mvdan.cc/gofumpt@latest.
        exit /b 1
    )
    echo gofumpt is installed.
    echo Formatting imports...
    goimports -w -local app .
    echo Formatting the code...
    gofumpt -l -w .
) else (
    echo Invalid argument. Use "tidy", "build", "test-io".
    exit /b 1
)