package main

// conda activate go_1.19.7

// A. Ramette 2025

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
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
func dereplicateSeqsParallel(seqs []*linear.Seq, numCores int) ([]*linear.Seq, []string) {
	if numCores <= 0 {
		numCores = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(numCores)

	seen := make(map[string]int)
	derep := []*linear.Seq{}
	names := []string{}
	var seenMutex sync.Mutex
	var derepMutex sync.Mutex
	var namesMutex sync.Mutex
	var wg sync.WaitGroup

	for _, seq := range seqs {
		wg.Add(1)
		go func(seq *linear.Seq) {
			defer wg.Done()
			seqString := string(seq.Seq)
			seenMutex.Lock()
			seen[seqString]++
			count := seen[seqString]
			seenMutex.Unlock()

			if count == 1 {
				derepMutex.Lock()
				derep = append(derep, seq)
				namesMutex.Lock()
				names = append(names, seq.Name()+";count_1")
				namesMutex.Unlock()
				derepMutex.Unlock()
			} else if count > 1 {
				derepMutex.Lock()
				for i, uniqueSeq := range derep {
					if string(uniqueSeq.Seq) == seqString {
						namesMutex.Lock()
						names[i] = uniqueSeq.Name() + ";count_" + fmt.Sprint(count)
						namesMutex.Unlock()
						break
					}
				}
				derepMutex.Unlock()
			}
		}(seq)
	}
	wg.Wait()

	return derep, names
}

func dereplicateSeqs(seqs []*linear.Seq) ([]*linear.Seq, []string) {

	seen := make(map[string]int)
	derep := []*linear.Seq{}
	names := []string{}

	for _, seq := range seqs {
		seqString := string(seq.Seq) // Convert sequence to string for comparison
		seen[seqString]++

		if seen[seqString] == 1 { // First time seeing this sequence
			derep = append(derep, seq)
			names = append(names, seq.Name()+";count_1")

		} else if seen[seqString] > 1 {
			for i, uniqueSeq := range derep {
				if string(uniqueSeq.Seq) == seqString {
					names[i] = uniqueSeq.Name() + ";count_" + fmt.Sprint(seen[seqString])
					break
				}
			}
		}
	}

	return derep, names
}

// //////////////////////////////////////////////////////////////////////////////////////
func main() {
	version := "0.1.1"
	inputFilePath := flag.String("i", "data_test/input.fasta", "Path to the input FASTA file")
	numCores := flag.Int("j", runtime.NumCPU()-1, "Number of CPU cores to use (default all available on the machine -1)")
	outputDir := flag.String("o", "output", "Output directory")
	forceOutput := flag.Bool("f", false, "Force output (overwrite existing files)")
	vers := flag.Bool("v", false, "Version of the software)")
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
		color.Red("          deduplicateseq\n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		// fmt.Println("\nNote:...)")

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
	// seqLen0 := len(seqs[0].Seq)
	// number of sequences at the beginning
	nseqs0 := len(seqs) //number of sequences (rows)

	// ----------------------------------------------------
	// dereplicating
	// dereplicatedSeqs, nameCounts := dereplicateSeqs(seqs)
	dereplicatedSeqs, nameCounts := dereplicateSeqsParallel(seqs, *numCores)

	// ----------------------------------------------------
	// Open the output file
	filename := *outputDir + "/dedup_MSA.fasta"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < len(dereplicatedSeqs); i++ {
		fmt.Fprintf(file, ">%s\n%s\n", nameCounts[i], dereplicatedSeqs[i].Seq)

	}
	FinalNberSeq := len(nameCounts)
	// ............................................................
	// Reporting
	// // ............................................................
	fmt.Println("=======================================================================")
	fmt.Println("(deduplicateseq version: ", version, ")")
	fmt.Println("= MSA input file:\t\t\t", *inputFilePath)
	fmt.Println("= Output directory:\t\t\t", *outputDir)
	fmt.Println("=> Number of initial seqs:\t\t\t", nseqs0)
	fmt.Println("=======================================================================")
	color.HiBlue("Results:")

	fmt.Println("=> Number of final seqs:\t\t\t", FinalNberSeq)
	fmt.Println("after deduplication of the sequences")

	fmt.Println("=======================================================================")
	fmt.Println("Files 1 FASTA file  written successfully to: ", filename)
	fmt.Println("Started at: ", formattedTime)
	endTime := time.Now()
	formattedTime = endTime.Format("2006-01-02 15:04:05") // Use a specific format
	fmt.Println("Finished at:", formattedTime)
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Elapsed time", elapsedTime)
	// ............................................................

}
