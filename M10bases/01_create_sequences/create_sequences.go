package main

import (
	"bufio"
	"fmt"
	"os"
)

// Helper to convert integer to a unique lowercase letter label
func indexToLabel(index int, length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	label := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		label[i] = letters[index%26]
		index /= 26
	}
	return string(label)
}

func main() {
	// Configuration
	const sequenceLength = 10
	const labelLength = 5
	bases := []byte{'A', 'T', 'G', 'C'}
	referenceBase := byte('A')
	fileName := "sequences10bases.fasta"

	// Create the output file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Use buffered writing to handle 1+ million lines efficiently
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	currentSeq := make([]byte, sequenceLength)
	count := 0

	// Recursive function
	var generate func(int)
	generate = func(position int) {
		// Base case: sequence is fully populated
		if position == -1 {
			// 1. Calculate mutations relative to "AAAAAAAAAA"
			mutations := 0
			for _, b := range currentSeq {
				if b != referenceBase {
					mutations++
				}
			}

			// 2. Generate unique lowercase ID
			label := indexToLabel(count, labelLength)

			// 3. Write to file: >[label]_muts_[count]
			fmt.Fprintf(writer, ">S%d%s\n", mutations, label)
			writer.Write(currentSeq)
			writer.WriteByte('\n')

			count++
			return
		}

		// Iterate through bases.
		// By recursing from the end of the string to the front (position - 1),
		// the "mutations" will appear at the start of the string first in the file.
		for _, base := range bases {
			currentSeq[position] = base
			generate(position - 1)
		}
	}

	fmt.Println("Generating sequences...")

	// Start the recursion at the last index to make the first index change fastest
	generate(sequenceLength - 1)

	fmt.Printf("Done! Results saved to %s\n", fileName)
}
