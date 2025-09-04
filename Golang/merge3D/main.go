package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

// mergeAndValidateCSVs reads two 3-column CSV files (without headers),
// validates that the label (3rd column) matches row-by-row,
// and writes a combined 4-column CSV (with header) to the output file.
// File 1 format: x,y,label
// File 2 format: x,z,label
// Output format: x,y,z,label
func mergeAndValidateCSVs(filePath1, filePath2, outputFilePath string) error {
	// --- Open File 1 ---
	f1, err := os.Open(filePath1)
	if err != nil {
		return fmt.Errorf("failed to open file 1 '%s': %w", filePath1, err)
	}
	defer f1.Close() // Ensure file is closed

	// --- Open File 2 ---
	f2, err := os.Open(filePath2)
	if err != nil {
		return fmt.Errorf("failed to open file 2 '%s': %w", filePath2, err)
	}
	defer f2.Close() // Ensure file is closed

	// --- Create CSV Readers ---
	reader1 := csv.NewReader(f1)
	reader1.FieldsPerRecord = 3 // Expect exactly 3 columns
	// reader1.Comment = '#' // Optional: Ignore lines starting with #

	reader2 := csv.NewReader(f2)
	reader2.FieldsPerRecord = 3 // Expect exactly 3 columns
	// reader2.Comment = '#' // Optional: Ignore lines starting with #

	// --- Create Output File and CSV Writer ---
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFilePath, err)
	}
	defer outFile.Close() // Ensure output file is closed

	writer := csv.NewWriter(outFile)
	defer writer.Flush() // Ensure all buffered data is written

	// --- Write Header to Output ---
	outputHeader := []string{"x", "y", "z", "label"}
	if err := writer.Write(outputHeader); err != nil {
		return fmt.Errorf("failed to write header to output file: %w", err)
	}

	// --- Process Rows ---
	rowNum := 0 // Start from row 0 (or 1 if you prefer 1-based logging)
	processedCount := 0
	mismatchCount := 0

	for {
		rowNum++ // Increment for each pair of rows attempted

		// Read one record from each file
		record1, err1 := reader1.Read()
		record2, err2 := reader2.Read()

		// --- Check for End of File (EOF) conditions ---
		eof1 := err1 == io.EOF
		eof2 := err2 == io.EOF

		if eof1 && eof2 {
			// Both files finished, normal exit
			log.Printf("Reached end of both files after processing %d rows in each file.", rowNum-1)
			break
		} else if eof1 {
			// log.Printf("Warning: Reached end of file '%s' (at row %d) but not '%s'. Files have different lengths.", filePath1, rowNum, filePath2)
			// break // Stop processing
			log.Fatalf("Warning: Reached end of file '%s' (at row %d) but not '%s'. Files have different lengths.", filePath1, rowNum, filePath2)
		} else if eof2 {
			// log.Printf("Warning: Reached end of file '%s' (at row %d) but not '%s'. Files have different lengths.", filePath2, rowNum, filePath1)
			// break // Stop processing
			log.Fatalf("Warning: Reached end of file '%s' (at row %d) but not '%s'. Files have different lengths.", filePath2, rowNum, filePath1)
		}

		// --- Check for other read errors ---
		if err1 != nil {
			// Check specifically for the wrong number of fields error
			if parseErr, ok := err1.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
				log.Printf("Warning: Skipping row %d in '%s' due to incorrect number of fields (expected 3). Content: %v\n", rowNum, filePath1, record1)
				// Try to read the corresponding row from file2 to keep sync, but skip processing this pair
				_, _ = reader2.Read() // Read and discard corresponding row from file2 if possible
				continue              // Skip this row pair
			}
			return fmt.Errorf("error reading record %d from '%s': %w", rowNum, filePath1, err1)
		}
		if err2 != nil {
			// Check specifically for the wrong number of fields error
			if parseErr, ok := err2.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
				log.Printf("Warning: Skipping row %d in '%s' due to incorrect number of fields (expected 3). Content: %v\n", rowNum, filePath2, record2)
				// We already read file1 successfully, but we skip this pair
				continue // Skip this row pair
			}
			return fmt.Errorf("error reading record %d from '%s': %w", rowNum, filePath2, err2)
		}

		// --- Extract data (assuming columns 0=x, 1=y/z, 2=label) ---
		x1 := record1[0]
		y1 := record1[1]
		label1 := record1[2]

		x2 := record2[0]
		z2 := record2[1]
		label2 := record2[2]

		// --- Validate Labels ---
		if label1 == label2 {
			// Labels match, proceed to write output
			processedCount++

			// Optional: Validate if x values match
			if x1 != x2 {
				// log.Printf("Warning: Row %d - 'x' values differ for label '%s'. Using x from '%s' ('%s'). ('%s' had '%s').", rowNum, label1, filePath1, x1, filePath2, x2)
				log.Fatalf("Warning: Row %d - 'x' values differ for label '%s'. Using x from '%s' ('%s'). ('%s' had '%s').", rowNum, label1, filePath1, x1, filePath2, x2)
			}

			// Prepare output row (using x from file1 as decided)
			outputRow := []string{x1, y1, z2, label1}

			// Write the combined row
			if err := writer.Write(outputRow); err != nil {
				return fmt.Errorf("failed to write record %d to output file: %w", rowNum, err)
			}
		} else {
			// Labels do not match for the same row number
			mismatchCount++
			// log.Printf("Warning: Row %d - Labels do not match! File 1 Label: '%s', File 2 Label: '%s'. Skipping row.", rowNum, label1, label2)
			log.Fatalf("Error: Row %d - Labels do not match! File 1 Label: '%s', File 2 Label: '%s'. Exiting.", rowNum, label1, label2)
			// Skip writing this row
		}
	} // End of loop

	// Final flush to ensure everything is written
	writer.Flush()
	if err := writer.Error(); err != nil { // Check for any error during flush
		return fmt.Errorf("error finalizing write to output file: %w", err)
	}

	log.Printf("Processing complete. Wrote %d records to '%s'.", processedCount, outputFilePath)
	if mismatchCount > 0 {
		log.Printf("Warning: Skipped %d rows due to label mismatches.", mismatchCount)
	}

	return nil // Success
}

// ---------------------------------------------------------------------------
func main() {
	version := "0.0.1"
	File1Path := flag.String("i", "./data_test/input1.csv", "Path to the input CSV file1 (coordinate file)")
	File2Path := flag.String("j", "./data_test/input2.csv", "Path to the input CSV file2 (coordinate file)")
	outputFilePath := flag.String("o", "./output_merge3D.csv", "Output CSV file with the columns merged")
	vers := flag.Bool("v", false, "Version of the software)")
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
		color.Red("          merge3D \n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println()
		fmt.Println("This software merges the information from csv files produced by build_axes into a 3D csv file (x,y,z,label).")
		fmt.Println("Each csv file contains three columns (x,y,label), and has no column headers.")
		fmt.Println("The first file, indicated by the flag \"-i\", determines \"x\" and \"label\" for both files.")
		fmt.Println("The second file, indicated by the flag \"-j\", can have only the second colum different. Otherwise, the 2 files are not merged.")
		fmt.Println()
		fmt.Println("Usage: go run main.go -i data_test/input1.csv -j data_test/input2.csv -o output.csv")
		fmt.Println("(run also \"go test\") for additional examples.")

		return // Exit
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................
	// --- Run the merge and validation ---
	log.Println("Starting CSV merge and validation...")
	// err = mergeAndValidateCSVs(*File1Path, *File2Path, *outputFilePath)
	err := mergeAndValidateCSVs(*File1Path, *File2Path, *outputFilePath)
	if err != nil {
		log.Fatalf("Script failed: %v", err)
	}

	log.Println("Script finished successfully.")

	// ------------------------------------------------------

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
