package main

// read CSV of the x,y,z labels
// split into fields
// clean the fields for usual issues: space beginning, end/ [ ] , \\'
// read badNames from file
// process the input file Reformating line by line

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"SILVA_go/utils"

	"github.com/fatih/color"
)

func main() {
	version := "0.0.3"
	inputCSVFile := flag.String("i", "./data_test/output_R1_R2_R3.csv", "Path to the input CSV file (coordinate file)")
	badNamesFilePath := flag.String("p", "./BadNames.txt", "Path to the bad names found in the taxonomy file. Each  pattern should be provided on a line")
	outputFile := flag.String("o", "output_cleaned_Taxo.csv", "Taxonomy fields cleaned.")
	TotalDesideredFields := flag.Int("F", 7, "The number of fields that are expected at the end.")
	deleteLast := flag.Bool("d", false, "delete the last field of the taxonomic string. (For instance, to remove \"count_i\")")
	vers := flag.Bool("v", false, "Version of the software")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse() // Parse the flags

	// ............................................................
	if *vers {
		fmt.Print("version: ")
		color.Red(version)
		return // Exit the program
	}
	if *help {
		color.Red("----------------------------------\n")
		color.Red("          SILVA_taxonomy_cleaner_go \n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette")
		fmt.Println("Date: 2025-08-02")
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("Usage:  go run main.go -i data_test/small.csv -o output_test.csv -p BadNames.txt -d -F 7")

		return // Exit
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................

	// Open the badNames file
	BadNamesSlice, err := utils.ReadFileToSlice(*badNamesFilePath)
	if err != nil {
		log.Fatal(err)
	}

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

	// Create a scanner to read the csv file line by line
	scanner := bufio.NewScanner(file)

	lineCounter := 0
	for scanner.Scan() {
		line := scanner.Text()

		// fmt.Println(line)
		// clean
		line = utils.CleanString(line)

		// Replace bad ";names;" by ";NA;"
		line = utils.ReplaceKeywordsinString(line, "NA", BadNamesSlice)

		// Count and Extend the missing fields
		Nfields := 1
		Nfields = utils.CountNberFieldinLine(line)
		if Nfields == 1 {
			fmt.Println("[Warning] Line ", lineCounter+1, ": No \";\" was found at all. Is this correct?")
		} else {

			ExtraFieldsNeeded := *TotalDesideredFields - Nfields
			// delete last if needed
			if *deleteLast {
				line = utils.DeleteFieldByPos(line, Nfields)
				line = utils.ExtendMissingFieldswithNA(line, ExtraFieldsNeeded+1)
				// fmt.Println(Nfields)
			} else {
				line = utils.ExtendMissingFieldswithNA(line, ExtraFieldsNeeded)
			}
		}

		// fmt.Println(line)
		// fmt.Println(Nfields)
		// fmt.Println("----------------------------")
		fmt.Fprintln(fileOutput, line)
		lineCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------
	fmt.Println("=======================================================================")
	color.HiBlue("Files:")
	fmt.Println("- Input file:\t\t", *inputCSVFile)
	fmt.Println("- Bad Names file:\t", *badNamesFilePath)
	fmt.Println("- Output file:\t\t", *outputFile)
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
