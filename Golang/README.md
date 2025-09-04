# Instructions for using the Golang scripts



## 0. Install go on your machine
After [Go installation](https://go.dev/doc/install), check that go indeed runs:
```
go version
```
should return the current version of the installed Go on your machine.

## 1. Clone the repository 
```
git clone https://github.com/RametteLab/CAPASYDIS/
```
and navigate to the folder `Golang/`
```
cd CAPASYDIS/Golang
```

Let's take the example of e.g. `build_axes` folder and scripts.

```
cd build_axes

go mod init build_axes
go mod tidy
```
This will install all the dependencies needed for running your application.


## 2. Option 1: Compile the Golang scripts
In the directory where a "main.go" in the build_axes directory is found do the following:
```
go mod tidy
go build -o build_axes main.go # -o is for the name of the executable
```
This will create a "build_axes" executable file in your current directory. This executable is platform-dependent.

## 3. Option 2: Run without compiling
This may be useful if you want to change the main.go code (at your own risk).

```
go run main.go -h
```
See then the *help* for further instructions about the possible flags and parameters.


## Disclaimer
a) This is my first Golang project. The code base and code structure will most likely need some revision.
b) The scripts have been developed to be used on a Linux servers.
Other operating systems have not been tested. 

## Could be done next:
- use a CLI such as cobra to bundle  all Golang scripts together.
- a TUI would also be good