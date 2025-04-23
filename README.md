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
3. Run the program
   - Linux/Unix
     ```bash
      make [lexer]
     ```
    - Windows
      ```bash
      make.bat [lexer]
      ```
