# Introduction

## What's this?

This is the course design for the Compiler Principles Lab of Fuzhou University for the 2025 academic year, 
implementing a simple compiler front-end. The lexical analysis refers to Golang's rules but does not fully 
cover them (there may be some bugs).

## How to run

1. Install Golang (1.23.6 or higher)
2. Install Go modules
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
   Here is help information for the `-h` flag:
   ```
   Usage of ./bin/mmap:
    -b    Enable benchmark mode
    -lexer--no-buffered
    Use no buffered reader for lexer
    -s    Stop writing results to file
    -t string
    Target to run: lexer or parser (default "lexer")
    ```