package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMainRun1(t *testing.T) {
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

func TestMainRunh(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-h")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -h returned an error: %v, stderr: %s", err, errOut.String())
	}

	// Add assertions to check the output.
	if bytes.Contains(out.Bytes(), []byte("findseq")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -h")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestInputMSAp(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//go run main.go -i data_test/inputMSA.fasta -p S1
	inputFile := "./data_test/inputMSA.fasta"
	outputFile := tempDir + "/output_selected_sequences.fasta"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-f", "-o", tempDir, "-p", "S1")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -i data_test/inputMSA.fasta -p S1 returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("\n**************************************************")
	fmt.Println("running: go run main.go -i data_test/inputMSA.fasta -p S1")
	fmt.Println()
	fmt.Println("=> input FASTA:", inputFile)
	// Open the FASTA file
	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	fmt.Println("\n=> output FASTA:", outputFile)
	// Open the FASTA file
	dat1, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat1))

	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	For p S1, 3 sequences matching S1 in the headers are retrieved.
	`)

	if bytes.Contains(out.Bytes(), []byte("Number of seqs in the output file:	 3")) { //tab + space
	} else {
		fmt.Println("Checking the number of the reported sequences:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong number reported. Expecting 3.")

	}

	fmt.Println("***********************************************************")
	// // Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)

	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {

		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}
}

func TestInputMSApc(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//go run main.go -i data_test/inputMSA.fasta -p S1 -c
	inputFile := "./data_test/inputMSA.fasta"
	outputFile := tempDir + "/output_selected_sequences.fasta"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-f", "-o", tempDir, "-p", "S1", "-c")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -i data_test/inputMSA.fasta -p S1 -c returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("\n**************************************************")
	fmt.Println("running: go run main.go -i data_test/inputMSA.fasta -p S1 -c")
	fmt.Println()
	fmt.Println("=> input FASTA:", inputFile)
	// Open the FASTA file
	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	fmt.Println("\n=> output FASTA:", outputFile)
	// Open the FASTA file
	dat1, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat1))

	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	For p S1 -c, 3 sequences matching S1 in the headers are retrieved. 
	-c cleans the sequences for . and -.
	`)

	if bytes.Contains(out.Bytes(), []byte("Number of seqs in the output file:	 3")) { //tab + space
	} else {
		fmt.Println("Checking the number of the reported sequences:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong number reported. Expecting 3.")

	}

	fmt.Println("***********************************************************")
	// // Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)

	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {

		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}
}

func TestInputMSAu(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//go run main.go -i data_test/inputMSA.fasta -p S0 -u
	inputFile := "./data_test/inputMSA.fasta"
	outputFile := tempDir + "/output_selected_sequences.fasta"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-f", "-o", tempDir, "-p", "S0", "-u")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -i data_test/inputMSA.fasta -p S0 -u returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("\n**************************************************")
	fmt.Println("running: go run main.go -i data_test/inputMSA.fasta -p S0 -u")
	fmt.Println()
	fmt.Println("=> input FASTA:", inputFile)
	// Open the FASTA file
	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	fmt.Println("\n=> output FASTA:", outputFile)
	// Open the FASTA file
	dat1, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat1))

	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	For p S0 -u, 1 sequence matching S0 in the header is retrieved. 
	-u converts "U" to "T".
	`)

	if bytes.Contains(out.Bytes(), []byte("Number of seqs in the output file:	 1")) { //tab + space
	} else {
		fmt.Println("Checking the number of the reported sequences:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong number reported. Expecting 1.")

	}

	fmt.Println("***********************************************************")
	// // Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)

	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {

		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}
}

func TestInputMSs(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	//go run main.go -i data_test/inputMSA.fasta -s CTCAACA
	inputFile := "./data_test/inputMSA.fasta"
	outputFile := tempDir + "/output_selected_sequences.fasta"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-f", "-o", tempDir, "-s", "CTCAACA", "-u")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -i data_test/inputMSA.fasta -s CTCAACA returned an error: %v, stderr: %s", err, errOut.String())
	}

	fmt.Println(out.String())
	fmt.Println("\n**************************************************")
	fmt.Println("running: go run main.go -i data_test/inputMSA.fasta -s CTCAACA")
	fmt.Println()
	fmt.Println("=> input FASTA:", inputFile)
	// Open the FASTA file
	dat, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	fmt.Println()
	fmt.Println("\n=> output FASTA:", outputFile)
	// Open the FASTA file
	dat1, err := os.ReadFile(outputFile)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat1))

	fmt.Println()
	fmt.Println(`
	Interpretation:
	***************
	With "-s CTCAACA", 1 sequence matching S1 in the sequence is retrieved. 
	`)

	if bytes.Contains(out.Bytes(), []byte("Number of seqs in the output file:	 1")) { //tab + space
	} else {
		fmt.Println("Checking the number of the reported sequences:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong number reported. Expecting 1.")

	}

	fmt.Println("***********************************************************")
	// // Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)

	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {

		fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}
}
