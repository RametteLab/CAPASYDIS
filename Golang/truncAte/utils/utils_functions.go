// Package utils provides utility functions for select_seqs.
package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

// DeduplicateSlice deduplicates a slice.
func DeduplicateSlice(slice []alphabet.Letters) ([]int, []alphabet.Letters) {
	encountered := map[string]bool{}
	result := make([]alphabet.Letters, 0, len(slice))
	indices := make([]int, 0, len(slice)) // Slice to store indices

	for i, item := range slice {
		key := string(item)
		if !encountered[key] {
			encountered[key] = true
			result = append(result, item)
			indices = append(indices, i) // Store the index
		}
	}

	return indices, result
}

// convertSeqsUtoT converts all 'U' bases to 'T' in a slice of linear.Seq using multiple cores.
func ConvertSeqsUtoT(seqs []*linear.Seq, numCores int) []*linear.Seq {
	runtime.GOMAXPROCS(numCores)

	var wg sync.WaitGroup
	wg.Add(len(seqs))

	result := make([]*linear.Seq, len(seqs))
	for i, seq := range seqs {
		go func(i int, seq *linear.Seq) {
			defer wg.Done()

			// Create a new sequence and copy the original data
			newSeq := linear.NewSeq(seq.ID, make([]alphabet.Letter, len(seq.Seq)), seq.Alphabet())
			copy(newSeq.Seq, seq.Seq)

			// Convert 'U' to 'T' in the new sequence
			for j, letter := range newSeq.Seq {
				if letter == 'U' {
					newSeq.Seq[j] = 'T'
				}
			}

			result[i] = newSeq
		}(i, seq)
	}

	wg.Wait()
	return result
}

// CountNotAnomalousBases counts the bases not matching a motif across sequences of a MSA.
func CountNotAnomalousBases(seqs []*linear.Seq, motif string) []int {
	if motif != "-" && motif != "." {
		fmt.Println("invalid argument: motif should be - or . Exiting.")
		os.Exit(1)
	}
	// length of the alignment (number of columns)
	seqLen := len(seqs[0].Seq)

	// Store the counts in a slice (protected by a mutex for parallelization)
	var mutex sync.Mutex
	baseCounts := make([]int, seqLen)

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(seqLen)

	for i := 0; i < seqLen; i++ { //for each base
		go func(i int) {
			defer wg.Done()
			for _, seq := range seqs { // Use seqs here
				letter := seq.Seq[i]
				if motif == "." {
					if letter != '.' {
						mutex.Lock()
						baseCounts[i]++
						mutex.Unlock()
					}
				} else if motif == "-" {
					if letter != '-' {
						mutex.Lock()
						baseCounts[i]++
						mutex.Unlock()
					}
				}
			}
		}(i)
	}

	wg.Wait()
	// return baseCounts
	return (baseCounts)
}

// func DeduplicateSlice(slice []alphabet.Letters) ([]alphabet.Letters, []int) {
// 	encountered := map[string]bool{}
// 	result := make([]alphabet.Letters, 0, len(slice))

// RemoveMotifColumns removes the columns that are conserved and that match a motif, typically . or -.
func RemoveMotifColumns(seqs []*linear.Seq, motif byte) []*linear.Seq {
	// Determine the indices of columns to remove
	seqLen := len(seqs[0].Seq)
	indicesToRemove := make([]int, 0)
	for i := 0; i < seqLen; i++ {
		allMatching := true
		for _, seq := range seqs { // keeping pos i fixed, loop across all sequences
			if seq.Seq[i] != alphabet.Letter(motif) { // if not matching the motif, don't remove! and stop iterating over all seqs on that position
				allMatching = false
				break
			}
		}
		if allMatching {
			indicesToRemove = append(indicesToRemove, i)
		}
	}

	// Create a map to store the indices to remove
	removeMap := make(map[int]bool) // the keys are integers (representing the indices), and the values are booleans
	for _, index := range indicesToRemove {
		removeMap[index] = true
	}

	// Create new sequences with the columns removed
	result := make([]*linear.Seq, len(seqs))
	for i, seq := range seqs {
		newSeq := linear.NewSeq(seq.ID, nil, seq.Alphabet())
		for j, letter := range seq.Seq {
			if !removeMap[j] {
				newSeq.Seq = append(newSeq.Seq, letter)
			}
		}
		result[i] = newSeq
	}

	return result
}

// RemoveMotifColumnsCores removes the columns that are conserved and that match a motif using the cores available.
func RemoveMotifColumnsCores(seqs []*linear.Seq, motif byte, numCores int) []*linear.Seq {
	// Set the number of CPU cores to use
	runtime.GOMAXPROCS(numCores)

	// Determine the indices of columns to remove in parallel
	seqLen := len(seqs[0].Seq)
	indicesToRemove := make([]int, 0)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(seqLen)

	for i := 0; i < seqLen; i++ {
		go func(i int) {
			defer wg.Done()
			allMatching := true
			for _, seq := range seqs {
				if seq.Seq[i] != alphabet.Letter(motif) {
					allMatching = false
					break
				}
			}
			if allMatching {
				mutex.Lock()
				indicesToRemove = append(indicesToRemove, i)
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Create a map to store the indices to remove
	removeMap := make(map[int]bool)
	for _, index := range indicesToRemove {
		removeMap[index] = true
	}

	// Create new sequences with the columns removed
	result := make([]*linear.Seq, len(seqs))
	for i, seq := range seqs {
		newSeq := linear.NewSeq(seq.ID, nil, seq.Alphabet())
		for j, letter := range seq.Seq {
			if !removeMap[j] {
				newSeq.Seq = append(newSeq.Seq, letter)
			}
		}
		result[i] = newSeq
	}

	return result
}

// CheckMSAforWobblesandConverttoN converts wobbles to N, converts U to T.
func CheckMSAforWobblesandConverttoN(seqs []*linear.Seq, numCores int) []*linear.Seq {
	runtime.GOMAXPROCS(numCores) // Set the number of CPU cores to use

	var wg sync.WaitGroup
	wg.Add(len(seqs))

	// Create a new slice to store the modified sequences
	result := make([]*linear.Seq, len(seqs))

	for i, seq := range seqs {
		go func(i int, seq *linear.Seq) {
			defer wg.Done()

			newSeq := linear.NewSeq(seq.ID, make([]alphabet.Letter, len(seq.Seq)), seq.Alphabet())
			copy(newSeq.Seq, seq.Seq) // Copy the original sequence

			for j, letter := range newSeq.Seq {
				if letter != 'A' && letter != 'T' && letter != 'G' && letter != 'C' && letter != '.' && letter != '-' {
					newSeq.Seq[j] = 'N'
				}
				if letter == 'U' {
					newSeq.Seq[j] = 'T'
				}
			}
			result[i] = newSeq
		}(i, seq)
	}

	wg.Wait()
	return result
}

// CheckMSAforWobblesandRemoveSeqs checks for the presence of wobbles in the seqs, and if any present, removes the sequences with non-alphabet
// bases, except for '.' and '-'.
func CheckMSAforNandRemoveSeqs(seqs []*linear.Seq, numCores int) []*linear.Seq {

	runtime.GOMAXPROCS(numCores) // Set the number of CPU cores to use

	var wg sync.WaitGroup
	wg.Add(len(seqs))

	var mutex sync.Mutex
	result := make([]*linear.Seq, 0) // Use a slice with initial length 0

	for _, seq := range seqs {
		go func(seq *linear.Seq) {
			defer wg.Done()

			isValid := true
			for _, letter := range seq.Seq {
				if letter == 'N' {
					isValid = false
					break
				}
			}

			if isValid {
				newSeq := linear.NewSeq(seq.ID, make([]alphabet.Letter, len(seq.Seq)), seq.Alphabet())
				copy(newSeq.Seq, seq.Seq) // Copy the original sequence
				mutex.Lock()
				result = append(result, newSeq)
				mutex.Unlock()
			}
		}(seq)
	}

	wg.Wait()
	return result
}

// truncate a string to a specific length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 0 {
		return ""
	}

	var result string
	count := 0
	for _, char := range s {
		charLen := utf8.RuneLen(char)
		if count+charLen > maxLength {
			break
		}
		result += string(char)
		count += charLen

	}
	return result
}

func DebugSeq(seqs []*linear.Seq, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, s := range seqs {
		if strings.Contains(s.ID, "Escherichia") {
			fmt.Fprintf(file, ">%s\n%s\n", s.ID, s.Seq)
		}

	}

	return nil
}

// LettersToString converts a slice of alphabet.Letter to a string.
func LettersToString(letters []alphabet.Letter) string {
	runes := make([]rune, len(letters))
	for i, letter := range letters {
		runes[i] = rune(letter) // Convert letter to rune
	}
	return string(runes)
}
