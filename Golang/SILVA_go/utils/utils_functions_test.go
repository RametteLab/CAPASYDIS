package utils

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// TestReadFileToSlice is a test function for readFileToSlice.
func TestReadFileToSlice(t *testing.T) {
	// Define a test file name and content.
	testFile := "test_keywords.txt"
	testContent := "keyword1\nkeyword2\nkeyword3"

	// Write the test content to a temporary file.
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Defer a cleanup function to remove the test file after the test.
	defer os.Remove(testFile)

	// Call the function under test.
	keywords, err := ReadFileToSlice(testFile)
	fmt.Println(keywords)
	// Check for any errors returned by the function.
	if err != nil {
		t.Errorf("readFileToSlice returned an unexpected error: %v", err)
	}

	// Define the expected output.
	expected := []string{"keyword1", "keyword2", "keyword3"}

	// Compare the actual output with the expected output.
	if !reflect.DeepEqual(keywords, expected) {
		t.Errorf("readFileToSlice produced unexpected output.\nExpected: %v\nGot: %v", expected, keywords)
	}
}

func TestCountNberFieldinLine(t *testing.T) {
	L1 := "HA782847.3.1866__Eukaryota;Amorphea"
	L2 := "HA782847.3.1866__Eukaryota;Amorphea;Obazoa;Opisthokonta;Holozoa;Choanozoa;Metazoa;Animalia;BCP;count_1"
	actual1 := CountNberFieldinLine(L1)
	expected1 := 2
	actual2 := CountNberFieldinLine(L2)
	expected2 := 10
	//
	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("ProcessData(%v) = %v; want %v", L1, actual1, expected1)
		t.Errorf("ProcessData(%v) = %v; want %v", L2, actual2, expected2)

	}
}

func TestReplaceKeywordsinString(t *testing.T) {
	replacement := "NA"
	L1 := "HA782847.3.1866__Eukaryota;unknown;Obazoa;Opisthokonta;Group;Choanozoa;Metazoa;insertae;BCP;count_1"
	BadNames := []string{"unknown", "Group", "insertae"}
	actual1 := ReplaceKeywordsinString(L1, replacement, BadNames)
	expected1 := "HA782847.3.1866__Eukaryota;NA;Obazoa;Opisthokonta;NA;Choanozoa;Metazoa;NA;BCP;count_1"
	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("ProcessData(%v) = %v; want %v", L1, actual1, expected1)
		fmt.Println("expected1:", expected1)
		fmt.Println("actual1:  ", actual1)
	}
}

func TestCleanString(t *testing.T) {
	A1 := CleanString("; one;")
	E1 := ";one;"
	A2 := CleanString(" one ' [two] (three)'")
	E2 := "one two three"
	// fmt.Println(A2)
	if !reflect.DeepEqual(A1, E1) {
		t.Errorf("ProcessData(%v) = %v; want %v", A1, A1, E1)
		t.Errorf("ProcessData(%v) = %v; want %v", A1, A2, E2)
	}
}

func TestExtendMissingFieldswithNA(t *testing.T) {
	// Given a taxonomy string S, add the N number of fields as ;NA;
	S := "keep;keep1"
	A1 := ExtendMissingFieldswithNA(S, 4)
	E1 := "keep;keep1;NA;NA;NA;NA"
	if !reflect.DeepEqual(A1, E1) {
		t.Errorf("ProcessData(%v) = %v; want %v", S, A1, E1)
	}

}

func TestDeleteFieldByPos(t *testing.T) {
	// Given a taxonomy string delete completely the field indicated by position in base 1
	S := "keep;keep1;to_be_deleted;keep2"
	P := 3
	A1 := DeleteFieldByPos(S, P)
	E1 := "keep;keep1;keep2"
	if !reflect.DeepEqual(A1, E1) {
		t.Errorf("ProcessData(%v) = %v; want %v", S, A1, E1)
	}

}
