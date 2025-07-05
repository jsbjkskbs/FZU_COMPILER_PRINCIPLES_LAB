<div align="center">

# Compiler Principles Lab - Fuzhou University

</div>

<div align="center">
English | <a href="README.zh.md">中文</a>
</div>

<div style="display: flex; flex-direction: row; isolation: isolate; width: 100%; height: 8rem; background: #29292c; border-radius: 0.75rem; overflow: hidden; font-size: 1.25rem; padding: 0.5rem 0.5rem; margin-top: 2rem;">
   <div style="width: 0.5rem; height: 100%; background: linear-gradient(to bottom, #2eadff, #3d83ff, #7e61ff); border-radius: 0.5rem;">
   </div>
   <div style="flex: 1; padding: 0.5rem; color: #ffffff;line-height: 1.5rem;">
         <div style="font-size: 1.5rem; font-weight: bold; color: #2eadff;"> 
            Notification to FZUer
         </div>
         <div style="flex: 1; padding: 2rem 0.5rem; color: #ffffff; line-height: 1.5rem;">
            FZUer, click there! 😋👉 <a href="/docs/note.md"> Note.md </a>
         </div>
   </div>
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