Instructions for using the Golang scripts



# 0. Make sure go is installed on your machine
```
go version
```
should return the current version of the installed Go on your machine.

# 1. Clone the repository 
```
git clone https://github.com/RametteLab/CAPASYDIS/
```
and navigate to the folder `Golang/`

## 2. Option 1: Compile the Golang scripts
In each directory, where a "main.go" is found do the following:
```
go mod tidy
go build main.go
```
This will create a "main" executable file in your current directory. You can rename that executable or move it wherever you want. The executable is platform-dependent.

## 3. Option 2: Run without compiling
This may be useful if you want to change the main.go code (at your own risk).

```
go run main.go -h
```
See then the *help* for further instructions about the possible flags and parameters.


# Disclaimer
a) This is my first Golang project. The code base and code structure will most likely need some revision.
b) The scripts have been developed to be used on a Linux servers.
Other operating systems have not been tested. 

# To be done next:
- use a CLI such as cobra to bundle  all Golang scripts together.
- a TUI would also be good