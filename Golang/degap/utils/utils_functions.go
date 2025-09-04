// Package utils provides utility functions for select_seqs.
package utils

import (
	"runtime"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

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
