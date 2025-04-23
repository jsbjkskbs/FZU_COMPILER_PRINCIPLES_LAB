@echo off
if "%1"=="lexer" (
    go run main.go lexer
) else (
    echo "Usage: make.cmd [lexer]"
)