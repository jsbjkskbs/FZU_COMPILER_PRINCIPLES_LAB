<div align="center">

# Compiler Principles Lab - Fuzhou University

</div>

<div align="center">
English | <a href="README.zh.md">中文</a>
</div>

# Introduction

## What's this?

This is the course design for the Compiler Principles Lab of Fuzhou University for the 2025 academic year, 
implementing a simple compiler front-end. The lexical analysis refers to Golang's rules but does not fully 
cover them (there may be some bugs).

## How to run

1. Install Golang (1.23.6 or higher)
2. Install Go modules(actually does not depend on any third-party libraries)
   ```bash
   go mod tidy
   ```
3. Build the program
   - Linux/Unix
      ```bash
      make build
      ```
   - Windows
      ```bash
      make.bat build
      ```
4. Run the binary in the `bin` directory
5. Use the `-h` flag to see the help information
   ```bash
   ./bin/xxx -h
   ```

## Documentation

1. [Lexer](/docs/lexer.md)
2. [Parser](/docs/parser.md)
3. [Intermediate Code Generation](/docs/intermediate-code-generation.md)