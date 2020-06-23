[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/EricNeid/go-mergkol)

# About

Merging multiple kotlin source code files into single file, removing all package declaration.
Usefull when using kotlin in coding game. <https://www.codingame.com/>

## Installation

Download tool.

```bash
go get github.com/EricNeid/go-mergkol/cmd/sleep
go install github.com/EricNeid/go-mergkol/cmd/sleep
```

## Usage

Make sure that $GOPATH/bin is in your path

```bash
go-mergkol.exe -h

go-mergkol.exe -dir input -o merged.kt
```
