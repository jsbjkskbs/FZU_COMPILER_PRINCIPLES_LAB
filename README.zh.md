<div align="center">

# 编译原理实验 - 福州大学

</div>

<div align="center">
<a href="README.md">English</a> | 中文
</div>


<div style="display: flex; flex-direction: row; isolation: isolate; width: 100%; height: 8rem; background: #29292c; border-radius: 0.75rem; overflow: hidden; font-size: 1.25rem; padding: 0.5rem 0.5rem; margin-top: 2rem;">
   <div style="width: 0.5rem; height: 100%; background: linear-gradient(to bottom, #2eadff, #3d83ff, #7e61ff); border-radius: 0.5rem;">
   </div>
   <div style="flex: 1; padding: 0.5rem; color: #ffffff;line-height: 1.5rem;">
         <div style="font-size: 1.5rem; font-weight: bold; color: #2eadff;"> 
            看到这个仓库的 FZUer 们可以看看这个👇
         </div>
         <div style="flex: 1; padding: 2rem 0.5rem; color: #ffffff; line-height: 1.5rem;">
            点这点这！😋👉 <a href="/docs/note.md"> Note.md </a>
         </div>
   </div>
</div>

# 简介

## 这是什么？

这是福州大学2025学年编译原理实验的课程设计，旨在实现一个简单的编译器前端。词法分析参考了Golang的规则，但并未完全覆盖（可能存在一些Bug）。

## 如何运行

1. 安装 Golang（1.23.6或更高版本）
2. 安装 Go modules（其实不依赖任何第三方库）
    ```bash
    go mod tidy
    ```
3. 构建程序
    - Linux/Unix
        ```bash
        make build
        ```
    - Windows
        ```bash
        make.bat build
        ```
4. 在`bin`目录下运行生成的二进制文件
5. 使用`-h`标志查看帮助信息
    ```bash
    ./bin/xxx -h
    ```

## 文档

1. [词法分析器](/docs/lexer.zh.md)
2. [语法分析器](/docs/parser.zh.md)
3. [中间代码生成](/docs/intermediate-code-generation.zh.md)