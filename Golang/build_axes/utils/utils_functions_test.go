package utils

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestRoundNber(t *testing.T) {
	// We can get this number to about 15-17 decimal places with float64 (see Sprintf)
	// After that, the representation will lose precision.
	SmallNber := 0.123456789012345678901
	//             +++++++++++++++++    17 decimals max
	// 0.0000000001
	// actual := RoundNber(verySmallNber, 0.000_000_000_1)
	actual := RoundNber(SmallNber, 1e-16)
	// fmt.Println(actual)
	expected := 0.12345678901234569 //68 at the end would be expected. But error!!!
	actual1 := RoundNber(SmallNber, 0.1)
	// fmt.Println(actual)
	expected1 := 0.1
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", SmallNber, actual, expected)
		t.Errorf("ProcessData(%v) = %v; want %v", SmallNber, actual1, expected1)
	}
	actual2 := fmt.Sprintf("%.20f", SmallNber)
	// expected2 := 3.1
	fmt.Println("SmallNber:", SmallNber)
	fmt.Println("Sprintf:  ", actual2)
	fmt.Printf("with 20f:  %.20f\n", SmallNber)
	fmt.Printf("with 17f:  %.17f\n", SmallNber)

	fmt.Println("\nMaximum float64 value:", math.MaxFloat64)

}

func TestCalculateDelta2seqs(t *testing.T) {
	REF1 := "AA"
	SEQ1 := "AA"
	REF2 := "AA"
	SEQ2 := "GC"
	DeltaVector := map[string]float64{
		"AG": 0.01, "GA": 0.01,
		"CT": 0.02, "TC": 0.02,
		"AC": 0.06, "CA": 0.06,
		"GT": 0.07, "TG": 0.07,
		"AT": 0.08, "TA": 0.08,
		"GC": 0.09, "CG": 0.09,
		"A-": 0.11, "-A": 0.11,
		"T-": 0.12, "-T": 0.12,
		"G-": 0.13, "-G": 0.13,
		"C-": 0.14, "-C": 0.14,
	}
	actual1 := CalculateDelta2seqs(REF1, SEQ1, DeltaVector)
	expected1 := []float64{0.0, 0.0}
	actual2 := CalculateDelta2seqs(REF2, SEQ2, DeltaVector)
	expected2 := []float64{0.0, 0.06}
	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("ProcessData(%v) = %v; want %v", DeltaVector, actual1, expected1)
		t.Errorf("ProcessData(%v) = %v; want %v", DeltaVector, actual2, expected2)
	}

}

func TestReturnHighestDelta1base(t *testing.T) {
	Base := "A"
	DV := map[string]float64{
		"AG": 0.01, "GA": 0.01,
		"CT": 0.02, "TC": 0.02,
		"AC": 0.06, "CA": 0.06,
		"GT": 0.07, "TG": 0.07,
		"AT": 0.08, "TA": 0.08,
		"GC": 0.09, "CG": 0.09,
		"A-": 0.11, "-A": 0.11,
		"T-": 0.12, "-T": 0.12,
		"G-": 0.13, "-G": 0.13,
		"C-": 0.14, "-C": 0.14,
	}
	max := ReturnHighestDelta1base(Base, DV)
	expected := 0.11
	if !reflect.DeepEqual(max, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", DV, max, expected)
	}

}

// func TestCheckMSAforWobbles(t *testing.T) {
// Example sequences
// seqs := []*linear.Seq{
// 	linear.NewSeq("Seq1", []alphabet.Letter("YCGTACGTR"), alphabet.DNA),
// 	linear.NewSeq("Seq2", []alphabet.Letter("TGCWAG-CT"), alphabet.DNA),
// 	linear.NewSeq("Seq3", []alphabet.Letter("AAAA.CCCC"), alphabet.DNA),
// 	linear.NewSeq("Seq4", []alphabet.Letter("AAAAUUUUU"), alphabet.DNA),
// 	linear.NewSeq("Seq5", []alphabet.Letter("NNNNNNNNN"), alphabet.DNA),
// }
// // Process the sequences with 5 cores
// actual := CheckMSAforWobbles(seqs, 5)
// actual_length := len(actual)
// // Print the modified sequences
// // for _, seq := range actual {
// // 	fmt.Println(seq)
// // }
// expected := []*linear.Seq{
// 	linear.NewSeq("Seq3", []alphabet.Letter("AAAA.CCCC"), alphabet.DNA),
// 	linear.NewSeq("Seq4", []alphabet.Letter("AAAAUUUUU"), alphabet.DNA),
// }

// expected_length := len(expected) //2
// if !reflect.DeepEqual(actual_length, expected_length) {
// 	t.Errorf("ProcessData(%v) = %v; want %v", seqs, actual_length, expected_length)
// }

// }

func TestCalculateSumSqrtScore(t *testing.T) {
	Vector := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Score := CalculateSumSqrtScore(Vector)
	Score_rounded := math.Floor(Score*100000) / 100000
	// fmt.Println("score: ", Score_rounded)
	// fmt.Println("score: ", Score)

	expected := 22.46827

	if !reflect.DeepEqual(Score_rounded, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", Vector, Score_rounded, expected)
	}
}

func TestCalculateMostDivergentScore(t *testing.T) {
	DeltaVectorConst := map[string]float64{
		"AG": 0.1, "GA": 0.1,
		"CT": 0.1, "TC": 0.1,
		"AC": 0.1, "CA": 0.1,
		"GT": 0.1, "TG": 0.1,
		"AT": 0.1, "TA": 0.1,
		"GC": 0.1, "CG": 0.1,
		"A-": 0.1, "-A": 0.1,
		"T-": 0.1, "-T": 0.1,
		"G-": 0.1, "-G": 0.1,
		"C-": 0.1, "-C": 0.1,
	}
	REF1 := "ACGTATGCAT"
	MostDivertscore := CalculateMostDivergentScore(REF1, DeltaVectorConst)
	MostDivertscore_rounded := math.Floor(MostDivertscore*100000) / 100000
	expected := 22.71691

	if !reflect.DeepEqual(MostDivertscore_rounded, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", REF1, MostDivertscore_rounded, expected)
	}
}

func TestCalculateLessDivergentScore(t *testing.T) {
	REF := "ACGTATGCAT"
	LessDivertscore := CalculateLessDivergentScore(REF)
	LessDivertscore_rounded := math.Floor(LessDivertscore*100000) / 100000
	expected := 22.46827
	if !reflect.DeepEqual(LessDivertscore_rounded, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", REF, LessDivertscore_rounded, expected)
	}
}

func TestAsymDist1Seq(t *testing.T) {
	DeltaVectorConst := map[string]float64{
		"AG": 0.1, "GA": 0.1,
		"CT": 0.1, "TC": 0.1,
		"AC": 0.1, "CA": 0.1,
		"GT": 0.1, "TG": 0.1,
		"AT": 0.1, "TA": 0.1,
		"GC": 0.1, "CG": 0.1,
		"A-": 0.1, "-A": 0.1,
		"T-": 0.1, "-T": 0.1,
		"G-": 0.1, "-G": 0.1,
		"C-": 0.1, "-C": 0.1,
	}
	REF1 := "ACGTATGCAT"
	SEQ1 := "AAGTATGCAT"
	// SEQ3 := "ACGTATGCTT"
	LessDivergScore := CalculateLessDivergentScore(REF1)
	MostDivergScore := CalculateMostDivergentScore(REF1, DeltaVectorConst)
	ASD1 := AsymDist1Seq(REF1, SEQ1, LessDivergScore, MostDivergScore, DeltaVectorConst)
	ASD1_rounded := math.Floor(ASD1*100000) / 100000
	expected1 := 0.14046
	if !reflect.DeepEqual(ASD1_rounded, expected1) {
		t.Errorf("ProcessData(%v) = %v; want %v", REF1, ASD1_rounded, expected1)
	}
	//
	// 	ASD2 := AsymDist1Seq(REF1, SEQ3, DeltaVectorConst)
	// 	ASD2_rounded := math.Floor(ASD2*100000) / 100000
	// 	expected2 := 0.06684
	// 	if !reflect.DeepEqual(ASD2_rounded, expected2) {
	// 		t.Errorf("ProcessData(%v) = %v; want %v", REF1, ASD2_rounded, expected2)
	// 	}
}

func TestFunctions(t *testing.T) {
	REF := "ACGTATGCAT"
	SEQ1 := "AAGTATGCAT"
	// SEQ2 := "ACGTATGCTT"
	DeltaVector := map[string]float64{
		"AG": 0.01, "GA": 0.01,
		"CT": 0.02, "TC": 0.02,
		"AC": 0.06, "CA": 0.06,
		"GT": 0.07, "TG": 0.07,
		"AT": 0.08, "TA": 0.08,
		"GC": 0.09, "CG": 0.09,
		"A-": 0.11, "-A": 0.11,
		"T-": 0.12, "-T": 0.12,
		"G-": 0.13, "-G": 0.13,
		"C-": 0.14, "-C": 0.14,
	}
	result1 := CalculateDelta2seqs(REF, SEQ1, DeltaVector) //  [0 0.06 0 0 0 0 0 0 0 0]
	Scores1 := CalculateSumSqrtScore(result1)              //  [0 0.06 0 0 0 0 0 0 0 0]
	Expct1 := float64(22.48933463327174)
	// fmt.Println("CalculateDelta2seqs result1:", Scores1)
	if !reflect.DeepEqual(Scores1, Expct1) {
		t.Errorf("ProcessData(%v) = %v; want %v", SEQ1, Scores1, Expct1)
	}
	// // testing constant deltavector as in the ms
	DeltaVectorConst := map[string]float64{
		"AG": 0.1, "GA": 0.1,
		"CT": 0.1, "TC": 0.1,
		"AC": 0.1, "CA": 0.1,
		"GT": 0.1, "TG": 0.1,
		"AT": 0.1, "TA": 0.1,
		"GC": 0.1, "CG": 0.1,
		"A-": 0.1, "-A": 0.1,
		"T-": 0.1, "-T": 0.1,
		"G-": 0.1, "-G": 0.1,
		"C-": 0.1, "-C": 0.1,
	}
	// // to calculate the Asymdist, one needs to standardize the scores, given the max divergence that the REF can produce
	result3 := CalculateDelta2seqs(REF, REF, DeltaVectorConst) //  [0 0.01 0 0 0 0 0 0 0 0]
	// fmt.Println("CalculateDelta2seqs result3 (REF constant):", result3)
	Scores3 := CalculateSumSqrtScore(result3) //
	// fmt.Println("CalculateSumSqrtScore result3:", Scores3)
	Expct3 := float64(22.4682781862041)
	// fmt.Println("CalculateDelta2seqs result1:", Scores1)
	if !reflect.DeepEqual(Scores3, Expct3) {
		t.Errorf("ProcessData(%v) = %v; want %v", REF, Scores3, Expct3)
	}

	LessDivergScore := CalculateLessDivergentScore(REF)
	MostDivergScore := CalculateMostDivergentScore(REF, DeltaVector)
	AsymDistSEQ2 := AsymDist1Seq(REF, SEQ1, LessDivergScore, MostDivergScore, DeltaVector)
	Expct4 := float64(0.06920517114739912)
	// fmt.Println("AsymDistSEQ2:", AsymDistSEQ2)
	if !reflect.DeepEqual(AsymDistSEQ2, Expct4) {
		t.Errorf("ProcessData(%v) = %v; want %v", REF, AsymDistSEQ2, Expct4)
	}
}

func TestMax(t *testing.T) {
	V := []float64{0.0, 10.0, 2}
	M, _ := Max(V)
	Expct := float64(10.0)
	// fmt.Println("Result:", M)
	if !reflect.DeepEqual(M, Expct) {
		t.Errorf("ProcessData(%v) = %v; want %v", V, M, Expct)
	}
}

func TestMean(t *testing.T) {
	V := []float64{0.0, 10.0, 2}
	M, _ := Mean(V)
	Expct := float64(4.0)
	// fmt.Println("Result:", M)
	if !reflect.DeepEqual(M, Expct) {
		t.Errorf("ProcessData(%v) = %v; want %v", V, M, Expct)
	}
}

func TestSD(t *testing.T) {
	V := []float64{1, 2.0, 3, 4, 1, 2, 3, 4}
	M, _ := SD(V)
	Expct := float64(1.1952286093343936)
	// fmt.Println("Result:", M)
	if !reflect.DeepEqual(M, Expct) {
		t.Errorf("ProcessData(%v) = %v; want %v", V, M, Expct)
	}
}

// testing decimal precision
func TestDecimalPrecision(t *testing.T) {
	// calculateAsymdistSimpleSNP
	observed := CalculateAsymdistSimpleSNP(10, 0.01, 10)
	expected := 0.06302773329183516
	// fmt.Println(observed)
	if !reflect.DeepEqual(observed, expected) {
		t.Errorf("ProcessData(%v) = %v; want %v", "0.01", observed, expected)
	}
	//
	fmt.Println("Testing long sequences")
	L := 1_000_000
	O1 := CalculateAsymdistSimpleSNP(L, 0.01, L)
	O2 := CalculateAsymdistSimpleSNP(L, 0.02, L)
	fmt.Println(O1)
	fmt.Println(O2)
	// 5.010482167288356e-07
	// 5.010492388602767e-07
	//       ^- i.e 12th position after the comma
	fmt.Println("The numbers in decimal notation are:")
	fmt.Printf("%.20f\n", O1)
	fmt.Printf("%.20f\n", O2)
	// 0.00000050104821672884
	// 0.00000050104923886028
	//              ^
	fmt.Println("The rounded numbers (using my function):")
	fmt.Println(RoundNber(O1, 0.000_000_001))
	fmt.Println(RoundNber(O2, 0.000_000_001))
	fmt.Println("The rounded numbers (using Sprintf function):")
	O1s := fmt.Sprintf("%.20f", O1)
	O2s := fmt.Sprintf("%.20f", O2)
	fmt.Println(O1s)
	fmt.Println(O2s)
}
