package main
// This script verifies that the fundamental building blocks are independent given a specific set of delta values. 
// It simulates all 100 possible position-delta combinations and checks for any shared square-free cores, which would indicate potential structural dependencies under Besicovitch's Theorem.
import (
	"fmt"
	"math"
)

// getSquareFreeCore extracts the square-free core of an integer 
// by stripping out all perfect square factors.
func getSquareFreeCore(n int64) int64 {
	var core int64 = 1

	// Check for factors of 2
	if n%2 == 0 {
		count := 0
		for n%2 == 0 {
			count++
			n /= 2
		}
		if count%2 != 0 {
			core *= 2
		}
	}

	// Check for odd factors
	var i int64
	limit := int64(math.Sqrt(float64(n)))
	for i = 3; i <= limit; i += 2 {
		if n%i == 0 {
			count := 0
			for n%i == 0 {
				count++
				n /= i
			}
			if count%2 != 0 {
				core *= i
			}
			// Recalculate limit as n shrinks to optimize loop
			limit = int64(math.Sqrt(float64(n)))
		}
	}

	// If n is still greater than 1, the remaining n is prime
	if n > 1 {
		core *= n
	}

	return core
}

func globalMatrixAudit(deltaPool []float64) {
	// Map to track which combination owns a specific square-free core
	coreMap := make(map[int64]string)
	
	type Collision struct {
		core  int64
		compA string
		compB string
	}
	var collisions []Collision

	fmt.Println("Analysis of all possible position-delta combinations ...")
	fmt.Println(string(make([]byte, 85))) // Print a separator line

	// Iterate through all 10 positions
	var pos int64
	for pos = 1; pos <= 10; pos++ {
		// Iterate through all available delta values
		for _, delta := range deltaPool {
			// Convert to integer radicand (multiplying by 100 to clear 2 decimal places)
			// Math.Round handles precision safety during float-to-int conversion
			radicand := int64(math.Round((float64(pos) + delta) * 100))
			core := getSquareFreeCore(radicand)

			combinationLabel := fmt.Sprintf("Pos %d with δ=%.2f", pos, delta)

			if existingLabel, exists := coreMap[core]; exists {
				// Avoid logging duplicates of the exact same pairings
				collisions = append(collisions, Collision{
					core:  core,
					compA: existingLabel,
					compB: combinationLabel,
				})
			} else {
				coreMap[core] = combinationLabel
			}
		}
	}

	// Report the Results
	if len(collisions) > 0 {
		fmt.Printf("⚠️ Audit failed: Found %d potential structural dependencies.\n\n", len(collisions))
		fmt.Printf("%-18s %-25s %-25s\n", "Square-Free Core", "Combination A", "Combination B")
		fmt.Println(string(make([]byte, 85)))
		for _, c := range collisions {
			fmt.Printf("%-18d %-25s %-25s\n", c.core, c.compA, c.compB)
		}
		fmt.Println(string(make([]byte, 85)))
		fmt.Println("Insight: Because these combinations share a base, specific cross-position mutations could mimic each other.")
	} else {
		fmt.Println("✅ Audit passed!")
		fmt.Println(string(make([]byte, 85)))
		fmt.Println("Result: All 100 possible position-delta combinations generate unique square-free cores.")
		fmt.Println("Conclusion: Dynamic assignment of these delta values cannot cause structural collisions under Besicovitch's Theorem.")
	}
}

func main() {
	//  pool of 10 available delta values (Table 1 in CAPASYDIS manuscript)
	deltaPool := []float64{0.01, 0.02, 0.06, 0.07, 0.08, 0.09, 0.11, 0.12, 0.13, 0.14}

	globalMatrixAudit(deltaPool)
}
