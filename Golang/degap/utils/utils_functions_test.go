package utils

import (
	"reflect"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

func TestRemoveMotifColumnsCores(t *testing.T) {
	// Example sequences
	seqs := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("AC--T"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TG--A"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AA--A"), alphabet.DNA),
	}

	motif := byte('-') // Motif to remove
	Ncores := 2
	actual := RemoveMotifColumnsCores(seqs, motif, Ncores)
	expected := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("ACT"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TGA"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AAA"), alphabet.DNA),
	}
	// Print the filtered sequences
	// 	for _, seq := range filteredSeqs {
	// 		fmt.Println(seq)
	// 	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual, expected)
	}
}
