package main

// conda activate go_1.19.7

// A. Ramette 2025

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/fatih/color"
)

// retrieve sequences from a a MSA based on the provided pattern
// //////////////////////////////////////////////////////////////////////////////////////
func main() {
	version := "0.0.1"
	inputFilePath := flag.String("i", "data_test/inputMSA.fasta", "Path to the input FASTA file")
	pattern := flag.String("p", "", "Pattern to find in the FASTA header")
	patternseq := flag.String("s", "", "Pattern to find in the FASTA sequence")
	convert := flag.Bool("u", false, "convert U to T? The conversion is done after searching for the pattern.")
	clean := flag.Bool("c", false, "Remove the - and . symbols from the sequences")
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
		color.Red("          findseq\n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("\nfindseq: Find sequences from a Multiple Sequence Alignment (MSA) FASTA file.")
		fmt.Println("The script does not deduplicate the MSA.")
		fmt.Println("usage: findseq -p S1 -o output/test1.fasta -c")
		fmt.Println("usage: findseq -s TGGTGAT -o output/test2.fasta -c")

		return // Exit the program
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................
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

	nseqs0 := len(seqs) // number of sequences at the beginning

	// Open the output file
	filename := *outputDir + "/output_selected_sequences.fasta"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	nseqoutput := 0

	if *pattern != "" && *patternseq != "" {
		fmt.Println("[error] Either -p or -s can be defined. Not both together. Exiting.")
		os.Exit(0)

	}

	if *pattern != "" { // Finding the pattern in headers
		for i := 0; i < nseqs0; i++ { //for each seq
			Header := seqs[i].ID + " " + seqs[i].Desc
			if regexp.MustCompile(*pattern).MatchString(Header) {
				fmt.Fprintf(file, ">%s\n", Header)
				nseqoutput++
				Seq := string(seqs[i].Seq)
				if *clean {
					Seq = strings.ReplaceAll(Seq, ".", "")
					Seq = strings.ReplaceAll(Seq, "-", "")

				}
				if *convert {
					Seq = strings.ReplaceAll(Seq, "U", "T")
					Seq = strings.ReplaceAll(Seq, "u", "T")
				}
				fmt.Fprintf(file, "%s\n", Seq)
			}

		}

	} else if *patternseq != "" { // Finding the pattern in seqs
		for i := 0; i < nseqs0; i++ { //for each seq
			Header := seqs[i].ID + " " + seqs[i].Desc
			Seq := string(seqs[i].Seq)
			if regexp.MustCompile(*patternseq).MatchString(Seq) {
				fmt.Fprintf(file, ">%s\n", Header)
				nseqoutput++

				if *clean {
					Seq = strings.ReplaceAll(Seq, ".", "")
					Seq = strings.ReplaceAll(Seq, "-", "")

				}
				if *convert {
					Seq = strings.ReplaceAll(Seq, "U", "T")
					Seq = strings.ReplaceAll(Seq, "u", "T")
				}
				fmt.Fprintf(file, "%s\n", Seq)
			}
		}
	} else {
		fmt.Println("[error] No pattern was defined. Exiting.")
		os.Exit(0)

	}
	// ............................................................
	// Reporting
	// ............................................................
	fmt.Println("=======================================================================")
	fmt.Println("(findseq version: ", version, ")")
	fmt.Println("= MSA input file:\t\t\t", *inputFilePath)
	fmt.Println("=> Number of seqs in the input file:\t", nseqs0) //
	fmt.Println("= Output file:\t\t\t", filename)
	fmt.Println("=> Number of seqs in the output file:\t", nseqoutput) //
	if *clean {
		fmt.Println("(Sequences cleaned for - and . symbols)") //
	}

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
