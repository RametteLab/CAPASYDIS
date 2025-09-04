package main

// conda activate go_1.19.7

// A. Ramette 2025

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"select_seqs_go/utils"
	"strings"
	"sync"
	"time"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/fatih/color"
)

// //////////////////////////////////////////////////////////////////////////////////////
func main() {
	version := "0.3.1"
	inputFilePath := flag.String("i", "data_test/inputMSA1.fasta", "Path to the input FASTA file")
	numCores := flag.Int("j", runtime.NumCPU()-1, "Number of CPU cores to use (default all available on the machine -1)")
	threshold := flag.Float64("k", 0.9, "Minimum proportion of bases at the start and end positions of the MSA without dots \".\" to keep the positions")
	outputDir := flag.String("o", "output", "Output directory")
	forceOutput := flag.Bool("f", false, "Force output (overwrite existing files)")
	vers := flag.Bool("v", false, "Version of the software")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse() // Parse the flags

	// ............................................................
	if *vers {
		fmt.Print("version:")
		color.Red(version)
		return // Exit the program
	}
	if *help {
		color.Red("----------------------------------\n")
		color.Red("          truncAte\n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("\nTruncAte: Truncates a SILVA MSA to keep only the positions of the MSA that are not . for a given k proportion of bases.")
		fmt.Println("Note: The order of the final sequences may be variable, as it depends on the synchronisation of the threads.")
		fmt.Println("The script does not deduplicate the MSA.")
		fmt.Println("DNA bases considered as valid are A, U, T, G, C, - (dash). The \".\" is considered as missing information.")
		fmt.Println("Positions only consisting of - or . are removed from the MSA in a second stage.")

		return // Exit the program
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................
	// Set the number of CPU cores to use
	// runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(*numCores)
	// Check if output directory exists
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.Mkdir(*outputDir, 0755) // Create with read/write/execute permissions for the owner
		if err != nil {
			panic(err)
		}
	} else if !*forceOutput {
		// If the directory exists and -f is not set, stop the program
		fmt.Println("Output directory already exists. Use -f to force overwrite.")
		return
	}

	// Open the FASTA file
	b, err := os.ReadFile(*inputFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: Input file '%s' does not exist.\n", *inputFilePath)
		} else {
			fmt.Printf("Error reading file: %v\n", err)
		}
		os.Exit(1) // Exit the program with an error code
	}

	msa := string(b)
	data := strings.NewReader(msa)

	template := linear.NewSeq("", nil, alphabet.DNAredundant)
	r := fasta.NewReader(data, template)
	sc := seqio.NewScanner(r)

	// Collect sequences in a slice of type linear.Seq
	seqs := make([]*linear.Seq, 0)
	for sc.Next() {
		s := sc.Seq().(*linear.Seq)
		seqs = append(seqs, s)
	}
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}
	// length of the alignment (number of columns)
	seqLen0 := len(seqs[0].Seq)
	// number of sequences at the beginning
	nseqs0 := len(seqs) //number of sequences (rows)
	// fmt.Println("ID:", seqs[0].ID) //KX009750.1.1381
	// fmt.Println("Desc:", seqs[0].Desc) // Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Klebsiella;Klebsiella aerogenes
	// fmt.Println("Alpha:", seqs[0].Alpha) // &{0xc000169500 0xc0000023c0}
	// fmt.Println("Annotation:", seqs[0].Annotation) //{KX009750.1.1381 Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Klebsiella;Klebsiella aerogenes <nil> + linear 0xc00013a2a0 0}

	// concatenating the space between ID and Desc with "__"
	for i := 0; i < nseqs0; i++ { //for each seqs
		if seqs[i].Desc != "" {
			seqs[i].ID = seqs[i].ID + "__" + seqs[i].Desc
		}
	}

	// working on the sequences ----------------------------
	//Conversion of wobbles to N, and U to T
	seqs = utils.CheckMSAforWobblesandConverttoN(seqs, *numCores)

	// utils.DebugSeq(seqs1, "debug.fasta")
	// os.Exit(0)

	//Removing sequences with  N
	seqs = utils.CheckMSAforNandRemoveSeqs(seqs, *numCores)
	// utils.DebugSeq(seqs, "debug.fasta")
	// os.Exit(0)

	//Removal of the - and . that are conserved vertically
	// seqs = utils.RemoveMotifColumnsCores(seqs, '-', *numCores)
	// seqs = utils.RemoveMotifColumnsCores(seqs, '.', *numCores)

	// for i := 0; i < len(seqs); i++ { //for each seqs
	// 	fmt.Println(seqs[i])
	// }

	// update nseqs !!!
	nseqs := len(seqs)

	// Store the counts of good bases i.e. not "." in a slice
	baseCountsNoDots := utils.CountNotAnomalousBases(seqs, ".")
	seqLen1 := len(seqs[0].Seq)
	// Create a new slice to store the proportions
	basePropNoDot := make([]float64, seqLen1)
	for i, count := range baseCountsNoDots {
		basePropNoDot[i] = float64(count) / float64(nseqs)
	}
	// Print the base counts
	// fmt.Println("Base counts no dots per column:\t\t", baseCountsNoDots)
	// fmt.Printf("Proportion of \"no dot\" per column:\t%0.1f", basePropNoDot)
	// const threshold = 0.9

	// Find indices of values larger than the threshold
	indices := make([]int, 0)
	for i, num := range basePropNoDot {
		if num >= *threshold {
			indices = append(indices, i)
		}
	}
	firstPosition := indices[0]             // Get the first element of the slice
	lastPosition := indices[len(indices)-1] // Get the last element of the slice

	// fmt.Println("=> New alignment length:\t\t", lastPosition-firstPosition+1) // Output: 12

	// Initialize slices with a defined size for parallelization
	goodSeqtrunc := make([]alphabet.Letters, 0, nseqs)
	goodIDtrunc := make([]string, 0, nseqs)
	// badSeqtrunc := make([]alphabet.Letters, 0, nseqs)
	// badIDtrunc := make([]string, 0, nseqs)

	// Create a wait group to synchronize goroutines
	var wg1 sync.WaitGroup
	wg1.Add(nseqs)

	// Use a mutex to protect the slices from concurrent access
	var mutex1 sync.Mutex

	for _, seq := range seqs {
		go func(seq *linear.Seq) {
			defer wg1.Done()

			var seqtrunc = seq.Seq[firstPosition : lastPosition+1]
			Lseqtrunc := len(seqtrunc)
			firstBase := string(seqtrunc[0])
			lastBase := string(seqtrunc[Lseqtrunc-1])
			const dot = "."
			if firstBase != dot && lastBase != dot {
				mutex1.Lock()
				goodSeqtrunc = append(goodSeqtrunc, seqtrunc)
				goodIDtrunc = append(goodIDtrunc, seq.ID)
				mutex1.Unlock()
			} else {
				// mutex1.Lock()
				// badSeqtrunc = append(badSeqtrunc, seqtrunc)
				// badIDtrunc = append(badIDtrunc, seq.ID)
				// mutex1.Unlock()
			}
		}(seq)
	}

	wg1.Wait() // Wait for all goroutines to finish

	// ------------------------------------------------------------
	// removing the columns of only - or .
	// utils.DebugSeqTrunc(goodIDtrunc, goodSeqtrunc, "debug.fasta")

	// fmt.Println(goodIDtrunc[0]) // the id
	// fmt.Println(goodSeqtrunc[0]) // the sequence
	seqLen := len(goodSeqtrunc[0])
	badPositions := make([]bool, seqLen)

	// Find positions with gaps in all sequences
	for i := 0; i < seqLen; i++ {
		allGaps := true
		for _, seq := range goodSeqtrunc {
			if seq[i] != alphabet.Letter('-') && seq[i] != alphabet.Letter('.') {
				allGaps = false
				break
			}
		}
		badPositions[i] = allGaps
	}
	// fmt.Println(badPositions)  /// these are the positions to remove now
	// fmt.Println(firstPosition) ///
	// fmt.Println(lastPosition)  ///

	// firstGoodPosition := 0
	// for i, p := range badPositions {
	// 	if !p {
	// 		firstGoodPosition = i
	// 		break
	// 	}
	// }
	// lastGoodPosition := len(badPositions)
	// for i := len(badPositions) - 1; i >= 0; i-- {
	// 	if !badPositions[i] {
	// 		lastGoodPosition = i
	// 		break
	// 	}

	// }

	// Create new sequences without gap positions
	result := make([][]alphabet.Letter, len(goodSeqtrunc))
	for i, seq := range goodSeqtrunc {
		newSeq := []alphabet.Letter{}
		for j, letter := range seq {
			if !badPositions[j] {
				// fmt.Println("letter:", letter)
				// fmt.Println("j:", j)
				newSeq = append(newSeq, letter)
			}
		}
		result[i] = newSeq

	}

	// Open the output file
	filename := *outputDir + "/output_new_MSA.fasta"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i, id := range goodIDtrunc {
		fmt.Fprintf(file, ">%s\n%s\n", id, utils.LettersToString(result[i]))
	}
	FinalAlignmentLength := len(result[0])
	nberFinalSeq := len(goodIDtrunc)
	roundedStr := float64(nberFinalSeq) * 100 / float64(nseqs0)

	// ............................................................
	// Reporting
	// ............................................................
	fmt.Println("=======================================================================")
	fmt.Println("(Select_seqs version: ", version, ")")
	fmt.Println("= MSA input file:\t\t\t", *inputFilePath)
	fmt.Println("= Output directory:\t\t\t", *outputDir)
	fmt.Println("=> Number of seqs:\t\t\t", nseqs0)             //
	fmt.Println("=> Number of aligned positions:\t\t", seqLen0) //
	fmt.Println("=> Threshold applied to keep a position:", *threshold)
	fmt.Println("=======================================================================")
	color.HiBlue("Results:")
	fmt.Println("1) After the analysis of the MSA:")
	fmt.Printf("    - the first position matching k is:\t%d\n", firstPosition+1)
	fmt.Printf("    - the last matching position is:\t%d\n", lastPosition+1)
	fmt.Println("(this is before removing columns of only . or -)")
	// fmt.Printf("    - the first kept position is:\t%d\n", firstGoodPosition+1)
	// fmt.Printf("    - the last kept position is:\t%d\n", lastGoodPosition+1)

	fmt.Println()
	fmt.Println("2) After removing sequences with N or wobbles, or starting or ending with .")
	fmt.Println("    - Final alignment length:\t\t", FinalAlignmentLength)
	fmt.Println("    - Final number of sequences:\t", nberFinalSeq)
	fmt.Printf("(%d / %d = %.1f %s", nberFinalSeq, nseqs0, roundedStr, "% of the initial number of sequences)\n")

	fmt.Println("=======================================================================")
	fmt.Println("1 FASTA file  written successfully to: ", filename)
	fmt.Println("Started at: ", formattedTime)
	endTime := time.Now()
	formattedTime = endTime.Format("2006-01-02 15:04:05") // Use a specific format
	fmt.Println("Finished at:", formattedTime)
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Elapsed time", elapsedTime)
	// ............................................................

}
