package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestInputMSA(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/input.fasta"
	outputFile := tempDir + "/dedup_MSA.fasta"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-f", "-o", tempDir)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("**************************************************")
	fmt.Println("=> input FASTA:", inputFile)
	// Open the FASTA file
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
	//
	fmt.Println("\n=> output FASTA:", outputFile)
	// Open the FASTA file
	dat1, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat1))

	f1, err := os.Open(outputFile)
	if err != nil {
		panic(err)
	}

	f1.Close()
	fmt.Println()
	fmt.Println(`
Interpretation: 
***************
Only the first unique name is kept. The header of the FASTA sequences are updated with ";count_i", where 
i represents the number of occurrence of the same sequence in the MSA.
	`)

	if bytes.Contains(out.Bytes(), []byte("Number of initial seqs:			 6")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported Number of initial seq:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong. Expecting 6.")

	}
	// checking for accuracy of the reported start and end positins
	if bytes.Contains(out.Bytes(), []byte("Number of final seqs:			 3")) {
		// fmt.Println("correct last position reported")
	} else {
		fmt.Println("Checking the Number of final seqs:")
		fmt.Println("ERROR!")
		t.Errorf("Wrong. Expecting 3.")
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
