package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	version := "0.1.3"
	inputCSVFile := flag.String("i", "./data_test/input2D.csv", "Path to the input CSV file (coordinate file)")
	patternFilePath := flag.String("p", "./data_test/Patterns.txt", "Path to the pattern file. Each  pattern should be provided on a line, such as \"Pattern\" -> color. e.g. \"Bacteria\" -> red")
	DefaultColor := flag.String("c", "black", "default color to use")
	outputFile := flag.String("o", "output2D_colored.csv", "Output CSV file with last column with color")
	vers := flag.Bool("v", false, "Version of the software")
	help := flag.Bool("h", false, "Show help message")
	// Use default DeltaVector if the flag is not set
	flag.Parse() // Parse the flags

	// ............................................................
	if *vers {
		fmt.Print("version:")
		color.Red(version)
		return // Exit the program
	}
	if *help {
		color.Red("----------------------------------\n")
		color.Red("          colorCSVTaxonomy \n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("Usage:  go run main.go -i data_test/input2D.csv -o output_test.csv -p data_test/Patterns.txt")

		return // Exit
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................

	// Open the TSV file
	fileTSV, err := os.Open(*patternFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileTSV.Close()

	// Create a map to store the data
	dataMap := make(map[string]string)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(fileTSV)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line by the space separator
		parts := strings.Split(line, "->")
		// fmt.Println(parts[0])
		// fmt.Println(parts[1])
		// Check if the line has exactly two columns
		if len(parts) == 2 {
			// Store the data in the map, using the first column as the key
			// and the second column as the value
			parts[0] = strings.Replace(parts[0], "\" ", "", 2)
			parts[0] = strings.Replace(parts[0], "\"", "", 2)
			parts[1] = strings.Replace(parts[1], " ", "", 1)
			dataMap[parts[0]] = parts[1]
			// fmt.Printf("-%s-\n", parts[0])
			// fmt.Printf("-%s-\n", parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Pattern used for coloring:")
	// for key, value := range dataMap {
	// 	fmt.Printf("  -%s- [%s]\n", key, value)
	// }
	// ------------------------------------------------------
	// output file
	fileOutput, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOutput.Close()

	// ------------------------------------------------------
	// // Open the CSV file
	if *inputCSVFile == "" {
		log.Fatal("Please provide the path to the CSV file using the flag")
	}

	file, err := os.Open(*inputCSVFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// checks if the first line of input has a header. if not, creates a header
	// Create a scanner to read the file line by line
	///...
	// var Headerparts []string
	// scanner0 := bufio.NewScanner(file)
	// scanned := scanner0.Scan() // Returns true if a line was scanned, false on EOF or error
	// ExistHeader := 0
	// NberHeaderparts := 0
	// // Check if Scan() was successful
	// if scanned {
	// 	// 3Get the text of the first line
	// 	firstLine := scanner0.Text()
	// 	// fmt.Printf("Successfully read first line: \n%s\n", firstLine)
	// 	Headerparts = strings.Split(firstLine, ",")
	// 	if Headerparts[0] == "x" {
	// 		// header exists
	// 		ExistHeader = 1
	// 	} else {
	// 		// fmt.Println("No Header was found. Creating it in the output file.")
	// 	}

	// 	NberHeaderparts = len(Headerparts)
	// 	// fmt.Println("Nber of elements in the first line: ", NberHeaderparts)
	// 	if NberHeaderparts == 3 { // 2D
	// 		fmt.Fprintf(fileOutput, "%s\n", "x,y,label,color")
	// 	}
	// 	if NberHeaderparts == 4 { // 3D
	// 		fmt.Fprintf(fileOutput, "%s\n", "x,y,z,label,color")
	// 	}

	// } else {
	// 	// Scan returned false. Check for errors or EOF.
	// 	if err := scanner0.Err(); err != nil {
	// 		// An actual error occurred during scanning
	// 		log.Fatalf("Error reading first line: %v", err)
	// 	}
	// }

	// Create a scanner to read the csv file line by line
	scanner1 := bufio.NewScanner(file)
	var Headerparts []string
	ExistHeader := 0
	NberHeaderparts := 0
	scanned := scanner1.Scan()
	counterColor := 0
	if scanned {
		firstline := scanner1.Text()
		// analyse the first line
		// fmt.Println(firstline)
		Headerparts = strings.Split(firstline, ",")
		NberHeaderparts = len(Headerparts)
		//
		if NberHeaderparts == 3 { // 2D
			fmt.Fprintf(fileOutput, "%s\n", "x,y,label,color")
		} else if NberHeaderparts == 4 { // 3D
			fmt.Fprintf(fileOutput, "%s\n", "x,y,z,label,color")
		} else {
			fmt.Println("NberHeaderparts: ", NberHeaderparts)
			log.Fatalln("Number of parts in the CSV file uncompatible. Exiting")
		}
		//

		if Headerparts[0] == "x" && Headerparts[1] == "y" {
			// header exists
			ExistHeader = 1
		} else { // the header does not exist.
			//The data of the first line needs to be processed, as it is not a header
			for key, value := range dataMap {
				if regexp.MustCompile(key).MatchString(firstline) {
					fmt.Fprintf(fileOutput, "%s,%s\n", firstline, value)
					counterColor++
				}
			}
			if counterColor == 0 {
				fmt.Fprintf(fileOutput, "%s,%s\n", firstline, *DefaultColor)

			}
		}

	} else {
		// Scan returned false. Check for errors or EOF.
		if err := scanner1.Err(); err != nil {
			// An actual error occurred during scanning
			log.Fatalf("Error reading first line: %v", err)
		} else {
			// No error, Scan() returned false because the input was empty (EOF immediately)
			fmt.Println("Input was empty, no lines to read.")
		}
	}

	///// =============================
	// continue scanning
	for scanner1.Scan() {
		line := scanner1.Text()
		counterColor := 0
		for key, value := range dataMap {
			if regexp.MustCompile(key).MatchString(line) {
				// fmt.Printf("%s=> %s\n", line, value)
				fmt.Fprintf(fileOutput, "%s,%s\n", line, value)
				counterColor++
			}
		}
		if counterColor == 0 {
			// fmt.Printf("%s,%s\n", line, *DefaultColor)
			fmt.Fprintf(fileOutput, "%s,%s\n", line, *DefaultColor)

		}
	}

	if err := scanner1.Err(); err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------
	fmt.Println("=======================================================================")
	color.HiBlue("Files:")
	fmt.Println("-Input file:\t", *inputCSVFile)
	if ExistHeader == 0 {
		fmt.Println("   - no header detected")
	} else {
		fmt.Println("   - header detected")
	}
	if NberHeaderparts == 3 {
		fmt.Println("   - 2D csv file (x,y,label)")
	} else if NberHeaderparts == 4 {
		fmt.Println("   - 3D csv file (x,y,z,label)")
	}

	fmt.Println("-Output file:\t", *outputFile)
	// Print the map
	fmt.Println()
	color.HiBlue("Pattern used for coloring:")
	for key, value := range dataMap {
		fmt.Printf(" \"%s\" .... [%s]\n", key, value)
	}
	fmt.Println()
	// fmt.Println("=======================================================================")
	color.HiBlue("Timing:")

	// fmt.Println("=======================================================================")
	fmt.Println("Started at: ", formattedTime)
	endTime := time.Now()
	formattedTime = endTime.Format("2006-01-02 15:04:05") // Use a specific format
	fmt.Println("Finished at:", formattedTime)
	elapsedTime := endTime.Sub(startTime)
	red := "\033[31m"
	reset := "\033[0m" // Reset color
	fmt.Printf("Elapsed time: %s%v%s\n", red, elapsedTime, reset)
	fmt.Println("=======================================================================")
}
