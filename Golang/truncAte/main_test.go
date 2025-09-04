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
	if bytes.Contains(out.Bytes(), []byte("version:0.3.1")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -v")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestMainRun2(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-h")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// Add assertions to check the output.
	if bytes.Contains(out.Bytes(), []byte("TruncAte")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -h")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestInputMSA(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/inputMSA.fasta"
	outputFile := tempDir + "/output_new_MSA.fasta"
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
For k=0.9, the first and second positions are eliminated, 
but the 3rd one is kept (first "G"). The last position is also eliminated ("T").
From the block remaining, S6 now starts with ".", and S5 ends with ".".
Those two sequences are thus removed from the alignment.
	`)

	if bytes.Contains(out.Bytes(), []byte("the first position matching k is:	3")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start and end positions:")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong first position reported. Expecting 3.")

	}
	// checking for accuracy of the reported start and end positins
	if bytes.Contains(out.Bytes(), []byte("the last matching position is:	14")) {
		// fmt.Println("correct last position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start and end positions:")
		fmt.Println("ERROR!")
		t.Errorf("Wrong last position reported. Expecting 14.")
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

func TestInputMSA2(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA2.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/inputMSA2.fasta"
	outputFile := tempDir + "/output_new_MSA.fasta"
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
As "-" is considered a valid base, conserved columns of "-" are accounted for during the first pass of the analysis that calculates the proportion of non "." bases.
In the second stage, columns consisting of only "-" are removed.
	`)

	if bytes.Contains(out.Bytes(), []byte("the first position matching k is:	1")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong first position reported. Expecting 1.")
	}
	if bytes.Contains(out.Bytes(), []byte("the last matching position is:	78")) {
		// fmt.Println("correct last position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported end position:")
		fmt.Println("ERROR!")
		t.Errorf("Wrong last position reported. Expecting 78.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final alignment length:		 30")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final alignment length. Expecting 30.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final number of sequences:	 5")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final number of sequences. Expecting 5.")

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

//

func TestInputMSAdots(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA_dots.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/inputMSA_dots.fasta"
	outputFile := tempDir + "/output_new_MSA.fasta"
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
#With k=0.9, 90% of sequences start at position 8 and end at position 19. 
Because sequences S5 and S6 start or end with "." at those positions, they are discarded from the final MSA.
	`)

	if bytes.Contains(out.Bytes(), []byte("the first position matching k is:	8")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong first position reported. Expecting 8.")
	}
	if bytes.Contains(out.Bytes(), []byte("the last matching position is:	19")) {
		// fmt.Println("correct last position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported end position:")
		fmt.Println("ERROR!")
		t.Errorf("Wrong last position reported. Expecting 19.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final alignment length:		 12")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final alignment length. Expecting 12.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final number of sequences:	 9")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final number of sequences. Expecting 9.")

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

func TestInputMSA_wobbles(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir, "my-temp-dir" is the prefix
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA_wobbles.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/inputMSA_wobbles.fasta"
	outputFile := tempDir + "/output_new_MSA.fasta"
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
As "-" is considered a valid base, conserved columns of "-" are accounted for during the first pass of the analysis that calculates the proportion of non "." bases.
In the second stage, columns consisting of only "-" or "." are removed, shrinking the number of positions further.
Sequences starting or ending with ".", or containing wobbles are discarded from the final MSA.
	`)

	if bytes.Contains(out.Bytes(), []byte("the first position matching k is:	3")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong first position reported. Expecting 3.")
	}
	if bytes.Contains(out.Bytes(), []byte("the last matching position is:	31")) {
		// fmt.Println("correct last position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported end position:")
		fmt.Println("ERROR!")
		t.Errorf("Wrong last position reported. Expecting 31.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final alignment length:		 13")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final alignment length. Expecting 13.")
	}
	if bytes.Contains(out.Bytes(), []byte("Final number of sequences:	 4")) {
		// fmt.Println("correct first position reported")
	} else {
		fmt.Println("Checking the accuracy of the reported start positions")
		fmt.Println("ERROR!")
		t.Errorf("\n[error] Wrong Final number of sequences. Expecting 4.")

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
