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
	if bytes.Contains(out.Bytes(), []byte("version:0.1.3")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -v")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestMain2D(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//  go run main.go -i data_test/input2D.csv -o output_test -p data_test/Patterns.txt
	inputFile := "./data_test/input2D.csv"
	PatternFile := "./data_test/Patterns.txt"

	outputFile := tempDir + "/output_test.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-p", PatternFile, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf(" \"go run main.go -i data_test/input2D.csv -o output_test -p data_test/Patterns.txt\" returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input1.csv -j data_test/input2.csv")
	fmt.Println("=> input csv file:", inputFile)

	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> Pattern File:", PatternFile)

	dat2, err := os.ReadFile(PatternFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(PatternFile)
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
	In the input file, no header was provided. The file starts directly with the data.
	The Pattern file must be exactly formatted as shown in this example. 
	The output file contains the data, with the color indicated in the Pattern file
	associated to each line of the input file. An header has been added to the new file.
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

func TestMain2DwHeader(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//  go run main.go -i data_test/input2D_withHeader.csv -o output_test -p data_test/Patterns.txt
	inputFile := "./data_test/input2D_withHeader.csv"
	PatternFile := "./data_test/Patterns.txt"

	outputFile := tempDir + "/output_test.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-p", PatternFile, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf(" \"go run main.go -i data_test/input2D_withHeader.csv -o output_test -p data_test/Patterns.txt\" returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input2D_withHeader.csv -j data_test/input2.csv")
	fmt.Println("=> input csv file:", inputFile)

	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> Pattern File:", PatternFile)

	dat2, err := os.ReadFile(PatternFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(PatternFile)
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
	In the input CSV file, a header was provided.
	The Pattern file must be exactly formatted as shown in this example. 
	The output file contains the data, with the color indicated in the Pattern file
	associated to each line of the input file. An header has been added to the new file.
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

func TestMain3D(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//  go run main.go -i data_test/input3D.csv -o output_test -p data_test/Patterns.txt
	inputFile := "./data_test/input3D.csv"
	PatternFile := "./data_test/Patterns.txt"

	outputFile := tempDir + "/output_test.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-p", PatternFile, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf(" \"go run main.go -i data_test/input3D.csv -o output_test -p data_test/Patterns.txt\" returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input3D.csv -j data_test/input2.csv")
	fmt.Println("=> input csv file:", inputFile)

	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> Pattern File:", PatternFile)

	dat2, err := os.ReadFile(PatternFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(PatternFile)
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
	In the 3D input CSV file, a header was not provided.
	The Pattern file must be exactly formatted as shown in this example. 
	The output file contains the data, with the color indicated in the Pattern file
	associated to each line of the input file. An header has been added to the new file.
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

func TestMain3DwHeader(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//  go run main.go -i data_test/input3D_withHeader.csv -o output_test -p data_test/Patterns.txt
	inputFile := "./data_test/input3D_withHeader.csv"
	PatternFile := "./data_test/Patterns.txt"

	outputFile := tempDir + "/output_test.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-p", PatternFile, "-o", outputFile)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf(" \"go run main.go -i data_test/input3D_withHeader.csv -o output_test -p data_test/Patterns.txt\" returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("running: go run main.go -i data_test/input3D_withHeader.csv -j data_test/input2.csv")
	fmt.Println("=> input csv file:", inputFile)

	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	f.Close()
	// 	//
	fmt.Println("=> Pattern File:", PatternFile)

	dat2, err := os.ReadFile(PatternFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat2))
	fmt.Println()
	f2, err := os.Open(PatternFile)
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
	In the input CSV file, a header was provided.
	The Pattern file must be exactly formatted as shown in this example. 
	The output file contains the data, with the color indicated in the Pattern file
	associated to each line of the input file. An header has been added to the new file.
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
