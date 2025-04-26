@echo off
if "%1"=="lexer" (
    go mod tidy
    go run main.go lexer
) else if "%1"=="parser" (
    go mod tidy
    go run main.go parser
) else if "%1"=="tidy" {
    go mod tidy
) else (
    echo Invalid argument. Use "lexer", "parser", "tidy", "test", "build", or "run".
)