package utils

import (
	"reflect"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

func TestDeduplicateSlice(t *testing.T) {
	// Correct way to define the slice of alphabet.Letters
	data := []alphabet.Letters{
		{'A', 'T', 'G', 'C'},
		{'A', 'T', 'G', 'C'},
		{'A', 'A', 'T', 'G'},
	}
	// actual
	Goodindices_actual, dedupedGoodSlice_actual := DeduplicateSlice(data)
	// expexted
	dedupedGoodSlice_expected := []alphabet.Letters{ //"[ATGC AATG]"
		{'A', 'T', 'G', 'C'},
		{'A', 'A', 'T', 'G'},
	}

	Goodindices_expected := []int{0, 2}
	// fmt.Println(dedupedGoodSlice_actual) // [GCTGATGCTGGC GGTGATGCTGGC]
	// fmt.Println(Goodindices_actual)      // Indices: [0 2]
	if !reflect.DeepEqual(dedupedGoodSlice_actual, dedupedGoodSlice_expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", data, dedupedGoodSlice_actual, dedupedGoodSlice_expected)
	}
	if !reflect.DeepEqual(Goodindices_actual, Goodindices_expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", data, Goodindices_actual, Goodindices_expected)
	}
}

func TestConvertSeqsUtoT(t *testing.T) {
	// Create some linear.Seq objects
	seqs := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("U.C-"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("A..G"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter(".-CU"), alphabet.DNA),
	}
	// Append the sequences to the slice

	actual := ConvertSeqsUtoT(seqs, 1)
	// fmt.Println(actual)
	expected := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("T.C-"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("A..G"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter(".-CT"), alphabet.DNA),
	}
	// Append the sequences to the slice
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual, expected)
	}

}

func TestCountNotAnomalousBases(t *testing.T) {
	// var seqs []linear.Seq = make([]linear.Seq, 0, 3)

	// Create some linear.Seq objects
	seqs := make([]*linear.Seq, 0)
	seq1 := linear.NewSeq("Seq1", []alphabet.Letter("..C-"), alphabet.DNA)
	seq2 := linear.NewSeq("Seq3", []alphabet.Letter("...G"), alphabet.DNA)
	seq3 := linear.NewSeq("Seq2", []alphabet.Letter(".-CT"), alphabet.DNA)

	// Append the sequences to the slice
	seqs = append(seqs, seq1)
	seqs = append(seqs, seq2)
	seqs = append(seqs, seq3)

	// // Print the sequences
	// for _, s := range seqs {
	// 	fmt.Println(s)
	// }
	// testing for .
	actual := CountNotAnomalousBases(seqs, ".")
	// fmt.Println(actual)
	expected := []int{0, 1, 2, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual, expected)
	}
	// testing for -
	actual1 := CountNotAnomalousBases(seqs, "-")
	// fmt.Println(actual1)
	expected1 := []int{3, 2, 3, 2}
	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual1, expected1)
	}
}

func TestRemoveDotsColumn(t *testing.T) {
	// Example sequences
	seqs := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("AC--T"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TG--A"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AA--A"), alphabet.DNA),
	}

	motif := byte('-') // Motif to remove

	actual := RemoveMotifColumns(seqs, motif)
	expected := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("ACT"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TGA"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AAA"), alphabet.DNA),
	}
	// Print the filtered sequences
	// for _, seq := range filteredSeqs {
	// 	fmt.Println(seq)
	// }
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual, expected)
	}
}

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

func TestCheckMSAforWobblesandConverttoN(t *testing.T) {
	// Example sequences
	seqs := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("YCGTACGTR"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TGCWAG-KT"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("MAAA.CCCC"), alphabet.DNA),
		linear.NewSeq("Seq4", []alphabet.Letter("AAAAUUUUU"), alphabet.DNA),
	}
	// Process the sequences with 5 cores
	actual := CheckMSAforWobblesandConverttoN(seqs, 5)
	// Print the modified sequences
	// for _, seq := range actual {
	// 	fmt.Println(seq)
	// }
	expected := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("NCGTACGTN"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TGCNAG-NT"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("NAAA.CCCC"), alphabet.DNA),
		linear.NewSeq("Seq4", []alphabet.Letter("AAAATTTTT"), alphabet.DNA),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual, expected)
	}

}

func TestCheckMSAforNandRemoveSeqs(t *testing.T) {
	// Example sequences
	seqs := []*linear.Seq{
		linear.NewSeq("Seq1", []alphabet.Letter("NCGTACGTR"), alphabet.DNA),
		linear.NewSeq("Seq2", []alphabet.Letter("TGCWAG-CT"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AAAA.CCCC"), alphabet.DNA),
		linear.NewSeq("Seq4", []alphabet.Letter("AAAAUUUUU"), alphabet.DNA),
		linear.NewSeq("Seq5", []alphabet.Letter("NNNNNNNNN"), alphabet.DNA),
	}
	// Process the sequences with 5 cores
	actual := CheckMSAforNandRemoveSeqs(seqs, 5)
	actual_length := len(actual)
	// Print the modified sequences
	// for _, seq := range actual {
	// 	fmt.Println(seq)
	// }
	expected := []*linear.Seq{
		linear.NewSeq("Seq2", []alphabet.Letter("TGCWAG-CT"), alphabet.DNA),
		linear.NewSeq("Seq3", []alphabet.Letter("AAAA.CCCC"), alphabet.DNA),
		linear.NewSeq("Seq4", []alphabet.Letter("AAAAUUUUU"), alphabet.DNA),
	}

	expected_length := len(expected) //2
	if !reflect.DeepEqual(actual_length, expected_length) {
		t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual_length, expected_length)
	}

}
