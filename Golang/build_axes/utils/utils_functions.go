package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

// ReadDeltaVectorFromFile reads the Delta vector of base pairs from a file
func ReadDeltaVectorFromFile(filename string) (map[string]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	DeltaVector := make(map[string]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line) // Split by whitespace
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in DeltaVector file: %s", line)
		}
		key := parts[0]
		var value float64
		_, err := fmt.Sscan(parts[1], &value) // Parse the float value
		if err != nil {
			return nil, fmt.Errorf("error parsing value in DeltaVector file: %s", line)
		}
		DeltaVector[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return DeltaVector, nil
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
				if !seq.Alphabet().IsValid(letter) && letter != '.' && letter != '-' {
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
// func CheckMSAforWobbles(seqs []*linear.Seq, numCores int) []*linear.Seq {
// }

// Sum of square roots of delta values
// The vector of delta values as compared to the aligned reference sequence.
// The score (numeric value) of the vector, as the sum of square root of each delta value. This is the unstandardized calculation. The score is not bound
//
// 666716.5  score of a reference (deltas are all 0) that is 10000-nt long
// 671.4629  score of a reference that is 100-nt long
// 22.46828  score of a reference that is 10-nt long

// Round to specific precision
func RoundNber(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// CalculateDelta2seqs calculates the delta vector between two sequences of characters
func CalculateDelta2seqs(REF string, SEQ string, DeltaVector map[string]float64) []float64 {
	V := make([]float64, len(REF))
	for i := 0; i < len(REF); i++ {
		x := string(REF[i]) + string(SEQ[i])
		if x == "AA" || x == "GG" || x == "CC" || x == "TT" || x == "--" {
			V[i] = 0
		} else {
			V[i] = DeltaVector[x] // retrieves the corresponding value for the pair from the DeltaVector
		}
	}
	return V
}

// ReturnHighestDelta takes one letter and returns the corresponding max value for that letter in DeltaVector
func ReturnHighestDelta1base(Base string, DeltaVector map[string]float64) float64 { // what about ties? => not possible if the DeltaVector is not redundant
	var max float64
	for pair, value := range DeltaVector {
		// fmt.Printf("%s -> %f\n", pair, value)
		if strings.Contains(pair, Base) {
			if value > max {
				max = value
			}
		}
	}
	return max // "Max value"
}

// -----------------------------------------------------------
func CalculateSumSqrtScore(vector []float64) float64 {
	// Create a new slice to store the score
	AugmentedVector := make([]float64, len(vector))
	for i := 0; i < len(vector); i++ {
		AugmentedVector[i] = math.Sqrt(float64(i+1) + vector[i]) // adding the delta to the index, and taking the sqrt of the whole
		// fmt.Println(i, " => ", AugmentedVector[i])
	}

	var sum float64
	for _, value := range AugmentedVector {
		sum += value
	}
	return sum
}

// CalculateMostDivergentScore calculates the most divergent SumSqrtScore that a REF sequence may produce, based on the given DeltaVector. This is used for standardization
func CalculateMostDivergentScore(REF string, DeltaVector map[string]float64) float64 {
	MostDivergentDeltaValues := make([]float64, len(REF))
	for i := 0; i < len(REF); i++ {
		Base := string(REF[i])
		MostDivergentDeltaValues[i] = ReturnHighestDelta1base(Base, DeltaVector)
		// fmt.Println(MostDivergentDeltaValues[i])
	}
	MostDivergentScore := CalculateSumSqrtScore(MostDivergentDeltaValues)
	return MostDivergentScore
}

// CalculateLessDivergentScore calculates the min SumSqrtScore for a given REF sequence
func CalculateLessDivergentScore(REF string) float64 {
	var LessDivergentScore float64
	for i := 0; i < len(REF); i++ {
		LessDivergentScore += math.Sqrt(float64(i + 1))
	}
	// LessDivergentScore := CalculateSumSqrtScore(AugmentedVector)
	return LessDivergentScore
}

// AsymDist1Seq calculates the asymdist between 1 sequence and 1 reference, given a a deltaVecto
// minScore and maxScore are used for the standardization of SumSqrtScore to asymdist
func AsymDist1Seq(REF string, SEQ string, minScore float64, maxScore float64, DeltaVector map[string]float64) float64 {
	SEQDelta2seqs := CalculateDelta2seqs(REF, SEQ, DeltaVector)
	SEQscore := CalculateSumSqrtScore(SEQDelta2seqs)
	ScoreRange := maxScore - minScore
	if ScoreRange <= 0 {
		fmt.Println("[error] AsymDist1Seq- The Score range value for is negative or null! Exiting")
		os.Exit(1)
	}
	asymdist := (SEQscore - minScore) / ScoreRange
	return (asymdist)
}

// Max finds the highest number in a slice
func Max(input []float64) (max float64, err error) {

	// Return an error if there are no numbers
	if len(input) == 0 {
		return math.NaN(), err
	}

	// Get the first value as the starting point
	max = input[0]

	// Loop and replace higher values
	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
		}
	}

	return max, nil
}

// Mean calculates the mean values of a slice
func Mean(input []float64) (max float64, err error) {

	// Return an error if there are no numbers
	if len(input) == 0 {
		return math.NaN(), err
	}

	sum := 0.0
	for _, v := range input {
		sum += v
	}

	return sum / float64(len(input)), nil

}

// SD calculates the standard deviation of a slice
func SD(data []float64) (max float64, err error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("cannot calculate SD of empty slice")
	}

	mean, err := Mean(data)
	if err != nil {
		return 0, err
	}

	sumSqDev := 0.0
	for _, v := range data {
		sumSqDev += (v - mean) * (v - mean)
	}

	variance := sumSqDev / float64(len(data)-1) // Use n-1 for sample variance
	return math.Sqrt(variance), nil
}

// Function to explore the production of asymdist value for specific delta and sequence length
// n: sequence length
// MutatedPosition: which position to be mutated by delta value
func CalculateAsymdistSimpleSNP(n int, delta float64, MutatedPosition int) float64 {
	SequenceNoMutation := make([]float64, n)

	MutatedSequence1last := make([]float64, n)
	MutatedSequence1last[MutatedPosition-1] = delta

	FullyMutated := make([]float64, n)
	for i := range FullyMutated {
		FullyMutated[i] = delta
	}

	// ...
	Score1SNP := CalculateSumSqrtScore(MutatedSequence1last)
	AbsMinScore := CalculateSumSqrtScore(SequenceNoMutation)
	AbsMaxScore := CalculateSumSqrtScore(FullyMutated)
	asymdistValue := (Score1SNP - AbsMinScore) / (AbsMaxScore - AbsMinScore)
	// ...
	// fmt.Println("SequenceNoMutation: ", SequenceNoMutation)
	// fmt.Println("MutatedSequence1last: ", MutatedSequence1last)
	// fmt.Println("FullyMutated sequence: ", FullyMutated)
	// fmt.Println("AbsMinScore: ", AbsMinScore)
	// fmt.Println("AbsMaxScore: ", AbsMaxScore)
	// fmt.Println("asymdistValue: ", asymdistValue) // 0.06302773329183516
	return asymdistValue
}
