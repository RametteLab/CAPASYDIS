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
	if bytes.Contains(out.Bytes(), []byte("version:0.1.8")) {
		// Test passes, help message found
	} else if bytes.Contains(errOut.Bytes(), []byte("flag provided but not defined: -v")) {
		//Test passes, -v flag is not defined, and correct error message is returned.
	} else {
		t.Errorf("Expected output not found. Output: %s, stderr: %s", out.String(), errOut.String())
	}

}

func TestLongMSA(t *testing.T) {
	// testing the basic script with default parameters
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}

	// go run main.go -i data_test/longMSA.fasta -o output_test -f -R1 1 -R2 2
	inputFile := "./data_test/longMSA.fasta"
	outputFile := tempDir + "/output.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-o", tempDir, "-f", "-R1", "1", "-R2", "2")
	// fmt.Println("Temporary directory", tempDir, "was created.")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// fmt.Println(out.String()) // to console print the terminal output
	if bytes.Contains(out.Bytes(), []byte("Elapsed time:")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.2480322982")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R1!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.3024777272")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R2!")
	}

	// open output file to check
	file2, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a scanner to read the csv file line by line
	scanner := bufio.NewScanner(file2)

	E := []string{"0,0.18887337270000001,DQ017911.1.1452__Bacteria;Chloroflexota;Anaerolineae;Aggregatilineales;A4b;Incertae;count_1__",
		"0.1885109831,0,JF429345.1.1501__Bacteria;Pseudomonadota;Gammaproteobacteria;Burkholderiales;Rhodocyclaceae;Zoogloea;uncultured;count_1__",
		"0.1659146621,0.1624251412,CP017112.924837.926386__Bacteria;Bacillota;Bacilli;Bacillales;Bacillaceae;Bacillus;Bacillus;count_1__",
		"0.1730019602,0.18596455180000002,GU305747.1.1458__Bacteria;Actinomycetota;Acidimicrobiia;Microtrichales;Ilumatobacteraceae;CL500-29;count_1__",
		"0.2480322982,0.3024777272,CU916528.1.1311__Archaea;Halobacteriota;Methanosarcinia;Methanosarcinales;Methanosaetaceae;Methanothrix;uncultured;count_1__"}

	// Create a slice to store the strings.
	var obtainedStrings []string
	for scanner.Scan() {
		line := scanner.Text()
		obtainedStrings = append(obtainedStrings, line)
	}

	for i := range 5 {
		if obtainedStrings[i] == E[i] {
			// fmt.Println(i, " OK")
		} else {
			fmt.Println(i, " NOT")
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

func TestLongMSAPrecision(t *testing.T) {
	// testing the basic script with default parameters
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}

	// go run main.go -i data_test/longMSA.fasta -o output_test -f -R1 1 -R2 2
	inputFile := "./data_test/longMSA.fasta"
	outputFile := tempDir + "/output.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-o", tempDir, "-f", "-R1", "1", "-R2", "2", "-f", "-r", "0.0001")
	// fmt.Println("Temporary directory", tempDir, "was created.")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// fmt.Println(out.String()) // to console print the terminal output
	if bytes.Contains(out.Bytes(), []byte("Elapsed time:")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.248")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R1!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.3025")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R2!")
	}

	// open output file to check
	file2, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a scanner to read the csv file line by line
	scanner := bufio.NewScanner(file2)

	E := []string{"0,0.1889,DQ017911.1.1452__Bacteria;Chloroflexota;Anaerolineae;Aggregatilineales;A4b;Incertae;count_1__",
		"0.1885,0,JF429345.1.1501__Bacteria;Pseudomonadota;Gammaproteobacteria;Burkholderiales;Rhodocyclaceae;Zoogloea;uncultured;count_1__",
		"0.16590000000000002,0.16240000000000002,CP017112.924837.926386__Bacteria;Bacillota;Bacilli;Bacillales;Bacillaceae;Bacillus;Bacillus;count_1__",
		"0.17300000000000001,0.186,GU305747.1.1458__Bacteria;Actinomycetota;Acidimicrobiia;Microtrichales;Ilumatobacteraceae;CL500-29;count_1__",
		"0.248,0.3025,CU916528.1.1311__Archaea;Halobacteriota;Methanosarcinia;Methanosarcinales;Methanosaetaceae;Methanothrix;uncultured;count_1__"}

	// Create a slice to store the strings.
	var obtainedStrings []string
	for scanner.Scan() {
		line := scanner.Text()
		obtainedStrings = append(obtainedStrings, line)
	}

	for i := range 5 {
		if obtainedStrings[i] == E[i] {
			// fmt.Println(i, " OK")
		} else {
			fmt.Println(i, " NOT")
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

func TestPrecision2(t *testing.T) {
	// testing the basic script with default parameters
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "outdir_temp") // "" means default temp dir
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}

	// go run main.go -i data_test/longMSA.fasta -o output_test -f -R1 1 -R2 2
	inputFile := "./data_test/longMSA.fasta"
	outputFile := tempDir + "/output.csv"
	cmd := exec.Command("go", "run", "main.go", "-i", inputFile, "-o", tempDir, "-f", "-R1", "1", "-R2", "2", "-f", "-r", "0.000_000_01")
	// fmt.Println("Temporary directory", tempDir, "was created.")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		t.Errorf("go run main.go -v returned an error: %v, stderr: %s", err, errOut.String())
	}

	// fmt.Println(out.String()) // to console print the terminal output
	if bytes.Contains(out.Bytes(), []byte("Elapsed time:")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.2480323")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R1!")
	}
	if bytes.Contains(out.Bytes(), []byte("value: 0.30247773")) {
	} else {
		fmt.Println("ERROR!")
		t.Errorf("\n[error] The output was not produced as expected for R2!")
	}

	// open output file to check
	file2, err := os.Open(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// Create a scanner to read the csv file line by line
	scanner := bufio.NewScanner(file2)

	E := []string{"0,0.18887337,DQ017911.1.1452__Bacteria;Chloroflexota;Anaerolineae;Aggregatilineales;A4b;Incertae;count_1__",
		"0.18851098,0,JF429345.1.1501__Bacteria;Pseudomonadota;Gammaproteobacteria;Burkholderiales;Rhodocyclaceae;Zoogloea;uncultured;count_1__",
		"0.16591466,0.16242514,CP017112.924837.926386__Bacteria;Bacillota;Bacilli;Bacillales;Bacillaceae;Bacillus;Bacillus;count_1__",
		"0.17300196,0.18596455,GU305747.1.1458__Bacteria;Actinomycetota;Acidimicrobiia;Microtrichales;Ilumatobacteraceae;CL500-29;count_1__",
		"0.2480323,0.30247773,CU916528.1.1311__Archaea;Halobacteriota;Methanosarcinia;Methanosarcinales;Methanosaetaceae;Methanothrix;uncultured;count_1__"}

	// Create a slice to store the strings.
	var obtainedStrings []string
	for scanner.Scan() {
		line := scanner.Text()
		obtainedStrings = append(obtainedStrings, line)
	}

	for i := range 5 {
		if obtainedStrings[i] == E[i] {
			// fmt.Println(i, " OK")
		} else {
			fmt.Println(i, " NOT")
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
