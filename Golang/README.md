---
title: Instructions for using the Golang scripts
---

# 1. clone the repository

## 2. Option 1: Compile the Golang scripts
In each directory, where a "main.go" is found do the following:
```{sh}
go mod tidy
go mod init

go build main.go
```
This will create a "main" executable file. YOu can rename that executable or move it wherever you want. It is platform-dependent.

## 3. Option 2: Run without compiling
This may be useufl if you want to change the main.go code (at your own risk).

```{sh}
go mod tidy
go mod init

go run main.go -h
```
See then the *help* for further instructions about the possible flags and parameters.


# Disclaimer
This is my first Golang project. The code base and code structure will most likely need some revision.

# To be done next:
- use a CLI such as cobra to bundle  all Golang scripts together.
- a TUI would also be good