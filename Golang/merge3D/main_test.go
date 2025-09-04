package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMainVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-v")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// Add assertions to check the output.
	if bytes.Contains(out.Bytes(), []byte("version:0.0.1")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -v")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestMainOK(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//  go run main.go -i ./data_test/input1.csv -j ./data_test/input2.csv
	inputFile1 := "./data_test/input1.csv"
	inputFile2 := "./data_test/input2.csv"

	outputFile := tempDir + "/output_test.csv"
	// cmd := exec.Command("go", "run", "main.go", "-i", "data_test/input1.csv", "-j", "data_test/input2.csv", "-o", "output_test.csv")
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile1, "-j", inputFile2, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -i ./data_test/input1.csv -j ./data_test/input2.csv returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input1.csv -j data_test/input2.csv")
	fmt.Println("=> input file1:", inputFile1)

	dat, err := os.ReadFile(inputFile1)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile1)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> input file2:", inputFile2)

	dat2, err := os.ReadFile(inputFile2)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(inputFile2)
	if err != nil {
		panic(err)
	}

	f2.Close()
	fmt.Println("\n=> output file:", outputFile)
	dato, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dato))

	fo, err := os.Open(outputFile)
	if err != nil {
		panic(err)
	}

	fo.Close()
	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	In file 1, x is used as the reference for R1.
	In file 2, "x" and the "label" columns should be identical as that in file 1. 
	Only the second column (i.e. z) can be different. Otherwise, the script is exiting.
	The output file merges the data from the two files, and adds a header to the resulting csv file.
		`)

	if bytes.Contains(out.Bytes(), []byte("Elapsed time")) {
		// fmt.Println("correct first position reported")
	} else {
		// fmt.Println("Checking the accuracy of the reported start and end positions:")
		// fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong ending. Expecting: \"Elapsed time\".")

	}
	fmt.Println("***********************************************************")
	// Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)
	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {
		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}

}

func TestMaininfo(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// testing different info in one file
	// go run main.go -i data_test/input1.csv -j data_test/input2_info.csv
	inputFile1 := "./data_test/input1.csv"
	inputFile2 := "./data_test/input2_info.csv"

	outputFile := tempDir + "/output_test.csv"
	// cmd := exec.Command("go", "run", "main.go", "-i", "data_test/input1.csv", "-j", "data_test/input2.csv", "-o", "output_test.csv")
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile1, "-j", inputFile2, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {

	} else {

		t.Errorf("go run main.go -i ./data_test/input1.csv -j ./data_test/input2_info.csv returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input1.csv -j data_test/input2_info.csv")
	fmt.Println("=> input file1:", inputFile1)

	dat, err := os.ReadFile(inputFile1)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile1)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> input file2:", inputFile2)

	dat2, err := os.ReadFile(inputFile2)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(inputFile2)
	if err != nil {
		panic(err)
	}

	f2.Close()
	fmt.Println("\n=> output file:", outputFile)
	dato, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dato))

	fo, err := os.Open(outputFile)
	if err != nil {
		panic(err)
	}

	fo.Close()
	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	In file 1, x is used as the reference for R1.
	In file 2, the label at the row 3 is different between the two files. The script exits. 
	The output file is empty.
	
		`)

	fmt.Println("***********************************************************")
	// Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)
	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {
		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}

}

func TestMainx(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// tesing different x info in one file
	// go run main.go -i data_test/input1.csv -j data_test/input2_x.csv

	inputFile1 := "./data_test/input1.csv"
	inputFile2 := "./data_test/input2_x.csv"

	outputFile := tempDir + "/output_test.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile1, "-j", inputFile2, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {

	} else {

		t.Errorf("go run main.go -i ./data_test/input1.csv -j ./data_test/input2_x.csv returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input1.csv -j data_test/input2_x.csv")
	fmt.Println("=> input file1:", inputFile1)

	dat, err := os.ReadFile(inputFile1)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile1)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> input file2:", inputFile2)

	dat2, err := os.ReadFile(inputFile2)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(inputFile2)
	if err != nil {
		panic(err)
	}

	f2.Close()
	fmt.Println("\n=> output file:", outputFile)
	dato, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dato))

	fo, err := os.Open(outputFile)
	if err != nil {
		panic(err)
	}

	fo.Close()
	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	In file 1, x is used as the reference for R1.
	In file 2, the x value of the 3rd row is different between the two files. The script exits. 
	The output file is empty.
	
		`)

	fmt.Println("***********************************************************")
	// Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)
	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {
		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}

}
