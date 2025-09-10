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


## 2. Option 1: Compile specific Golang scripts
In the directory where a "main.go" in the build_axes directory is found do the following:
```
go mod tidy
go build -o build_axes main.go # -o is for the name of the executable
```
This will create a "build_axes" executable file in your current directory. This executable is platform-dependent.

## 3. Option 2: Compile all Golang scripts

```
cd CAPASYDIS/Golang

mkdir -p bin

for FOLDER in build_axes colorCSVTaxonomy deduplicateseq degap findseq merge3D SILVA_go truncAte
do
	cd $FOLDER
	go mod init $FOLDER
	go mod tidy
	go build -o $FOLDER main.go
	cp 	$FOLDER ../bin/
	cd ..
done	
```
The `Golang/bin` directory should now contain the compiled binaries.
Check with one of the binaries:
```
cd Golang/bin
./build_axes -h
```
## 4. Option 3: Run without compiling
This may be useful if you want to change the main.go code (at your own risk).

```
go run main.go -h
```
See then the *help* for further instructions about the possible flags and parameters.

## 5. Check the test to see how the functions work
To inspect the implemented functions in each of the package directories, you can run the following:

```
cd colorCSVTaxonomy/ # for instance
go test -v  -run Test ./... -count=1
```
This will print the detail output of the tests that were developed to create the main functions. 
It can be used as a documentation of the expected behavior of the implemented functions.

## Disclaimer
a) This is my first (large) Golang project. The code base and code structure will most likely need some revision.   
b) The scripts have been developed to be used on Linux servers.
Other operating systems have not been tested. 

## Could be done next:
- use a CLI such as cobra to bundle  all Golang scripts together.
- a TUI would also be good
