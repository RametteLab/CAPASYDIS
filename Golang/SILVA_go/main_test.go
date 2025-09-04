package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
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
	if bytes.Contains(out.Bytes(), []byte("version: 0.0.3")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -v")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestInputSmall(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/small.csv"
	outputFile := tempDir + "/output_test.csv"
	// go run main.go -i data_test/small.csv -o output_test.csv -p BadNames.txt -d -F 7
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-o", outputFile, "-p", "BadNames.txt", "-d", "-F", "7")
	// fmt.Println("Temporary directory", tempDir, "was created.")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// fmt.Println(out.String())
	if bytes.Contains(out.Bytes(), []byte("Elapsed time:")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected!")
	}

	// open output file to check
	file2, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a scanner to read the csv file line by line
	scanner := bufio.NewScanner(file2)

	E := []string{"0.010368,0.102936,0.057836,format__Bacteria;Pseudomonadota;Gammaproteobacteria;Coxiellales;Coxiellaceae;Coxiella;Ornithodoros",
		"0.004106,0.098028,0.052806,badNames__Bacteria;NA;NA;NA;NA;NA;NA",
		"0.010368,0.102936,0.057836,2__Bacteria;Pseudomonadota;NA;NA;NA;NA;NA",
		"0.004106,0.098028,0.052806,3__Bacteria;Pseudomonadota;Gammaproteobacteria;NA;NA;NA;NA",
		"0.004106,0.098028,0.052806,Normal__Bacteria;Pseudomonadota;Gammaproteobacteria;Pseudomonadales;Pseudomonadaceae;Pseudomonas;Pseudomonas"}

	// i:=0
	// Create a slice to store the keywords.
	var obtainedStrings []string
	for scanner.Scan() {
		line := scanner.Text()
		obtainedStrings = append(obtainedStrings, line)
	}
	//
	for i := 0; i < 5; i++ {
		if obtainedStrings[i] == E[i] {
			// fmt.Println(i, " OK")
		} else {
			// fmt.Println(i, " NOT")
			t.Errorf("\n[error] The output was not produced as expected!")
			fmt.Println("	Expected:", E[i])
			fmt.Println("	Obtained:", obtainedStrings[i])

		}
	}

	// Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)
	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {
		// fmt.Println("Temporary directory", tempDir, "was removed.")
		fmt.Println()
	}
}

func TestInputSmall2(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	// go run main.go -i ./data_test/inputMSA.fasta -o ./output_test -f -j 5
	inputFile := "./data_test/small2.csv"
	outputFile := tempDir + "/output_test2.csv"
	// go run main.go -i data_test/small.csv -o output_test.csv -p BadNames.txt  -F 7 # no "d" here
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-o", outputFile, "-p", "BadNames.txt", "-F", "7")
	// fmt.Println("Temporary directory", tempDir, "was created.")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// fmt.Println(out.String())
	if bytes.Contains(out.Bytes(), []byte("Elapsed time:")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected!")
	}
	// open input file to check
	fileIn, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileIn.Close()

	// open output file to check
	file2, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a scanner to read the csv file line by line
	scannerIn := bufio.NewScanner(fileIn)
	// Create a slice to store the obtained Strings.
	var InputStrings []string
	for scannerIn.Scan() {
		lineIn := scannerIn.Text()
		InputStrings = append(InputStrings, lineIn)
	}

	// Create a scanner to read the csv file line by line
	scanner2 := bufio.NewScanner(file2)
	// Create a slice to store the obtained Strings.
	var obtainedStrings []string
	for scanner2.Scan() {
		line2 := scanner2.Text()
		obtainedStrings = append(obtainedStrings, line2)
	}
	// Expected
	E := []string{"0.010368,0.102936,0.057836,format__Bacteria;Pseudomonadota;Gammaproteobacteria;Coxiellales;Coxiellaceae;Coxiella;Ornithodoros",
		"0.004106,0.098028,0.052806,badNames__Bacteria;NA;NA;NA;NA;NA;NA",
		"0.010368,0.102936,0.057836,2__Bacteria;Pseudomonadota;NA;NA;NA;NA;NA",
		"0.004106,0.098028,0.052806,3__Bacteria;Pseudomonadota;Gammaproteobacteria;NA;NA;NA;NA",
		"0.004106,0.098028,0.052806,Normal__Bacteria;Pseudomonadota;Gammaproteobacteria;Pseudomonadales;Pseudomonadaceae;Pseudomonas;Pseudomonasfluorescens"}

	//
	for i := 0; i < 5; i++ {
		if obtainedStrings[i] == E[i] {
			// fmt.Println(i, " OK")
		} else {
			t.Errorf("\n[error]")
			fmt.Println("Case: ", i)
			fmt.Println("	Input:", InputStrings[i])
			fmt.Println("	Expected:", E[i])
			fmt.Println("	Obtained:", obtainedStrings[i])

		}
	}

	// Clean up the temporary directory and its contents when you're done.
	errT := os.RemoveAll(tempDir)
	if errT != nil {
		fmt.Println("Error removing temporary directory:", errT)
		return
	} else {
		// fmt.Println("Temporary directory", tempDir, "was removed.")
		// fmt.Println()
	}
}
