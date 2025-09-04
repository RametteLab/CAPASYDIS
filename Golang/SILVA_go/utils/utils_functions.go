package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readFileToSlice reads a file line by line and returns the lines as a slice of strings.
// It is designed for files where each line is a single keyword.
func ReadFileToSlice(filePath string) ([]string, error) {
	// Open the file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	// Ensure the file is closed when the function returns.
	defer file.Close()

	// Use a scanner to read the file line by line efficiently.
	scanner := bufio.NewScanner(file)

	// Create a slice to store the keywords.
	var keywords []string

	// Loop through each line of the file.
	for scanner.Scan() {
		line := scanner.Text()
		keywords = append(keywords, line)
	}

	// Check for any errors that occurred during scanning.
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	// Return the populated slice and a nil error.
	return keywords, nil
}

func CountNberFieldinLine(LineData string) int {
	semicolonCount := 0
	semicolonCount = strings.Count(LineData, ";")
	return semicolonCount + 1
}

// replaceKeyword checks if a keyword is present in a string and replaces it.
// It takes the original string as input and returns the modified string.
// List to replace is case-sensitive!
func ReplaceKeywordsinString(inputString string, replacement string, ListToReplace []string) string {
	modifiedString := inputString // Start with a copy of the original string
	replacement1 := ";" + replacement + ";"
	replacement2 := ";" + replacement // for ending bad word
	// Iterate over the slice of keywords
	for _, keyword := range ListToReplace {
		keyword1 := ";" + keyword + ";"
		keyword2 := ";" + keyword
		if strings.Contains(modifiedString, keyword) {
			modifiedString = strings.ReplaceAll(modifiedString, keyword1, replacement1)
		}
		if strings.HasSuffix(modifiedString, keyword2) {
			modifiedString = strings.Replace(modifiedString, keyword2, replacement2, 1)
		}

	}

	return modifiedString
}

func CleanString(inputString string) string {

	replacer := strings.NewReplacer(
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"'", "",
		" ", "")
	// Use the Replacer's Replace method to get the modified string.
	modifiedString := replacer.Replace(inputString)
	return modifiedString
}

// Given a taxonomy string S, add the N number of fields as ;NA;
// N:=4
// S:= "keep;keep1"
// Result="keep;keep1;NA;NA;NA;NA"
func ExtendMissingFieldswithNA(S string, N int) string {
	// A strings.Builder is more efficient than repeated string concatenation
	// as it minimizes memory allocations.
	var builder strings.Builder

	// Write the initial string to the builder.
	builder.WriteString(S)

	// Append ";NA" for N number of times.
	for i := 0; i < N; i++ {
		builder.WriteString(";NA")
	}

	// Return the final string.
	return builder.String()

}

// Given a taxonomy string delete completely the field indicated by position in base 1
// pos:=3
// s:=keep;keep1;to_be_deleted;keep2
// result:= keep;keep1;keep2
func DeleteFieldByPos(S string, pos int) string {
	// Split the string into fields.
	fields := strings.Split(S, ";")
	numFields := len(fields)
	// Validate the position. It must be a 1-based index within the bounds of the fields.
	if pos <= 0 || pos > numFields {
		return S // Return the original string for an invalid position.
	}
	// Convert the 1-based position to a 0-based slice index.
	indexToDelete := pos - 1
	// Create a new slice by appending the parts before and after the
	// element to be deleted. This is the idiomatic way to remove an element
	// from a slice in Go.
	// The `...` operator unpacks the second slice so it can be appended.
	newFields := append(fields[:indexToDelete], fields[indexToDelete+1:]...)

	// Join the new slice of fields back into a single string.
	return strings.Join(newFields, ";")
}
