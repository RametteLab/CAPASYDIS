package main

// conda activate go_1.19.7

// A. Ramette 2025

import (
	"degap_go/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/fatih/color"
)

// //////////////////////////////////////////////////////////////////////////////////////
func main() {
	version := "0.1.1"
	inputFilePath := flag.String("i", "data_test/input_new_MSA.fasta", "Path to the input FASTA file")
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
		color.Red("          degap\n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("\nNote: The order of the final sequences may be variable, as it depends on the core synchronisation)")
		fmt.Println("Columns of bases consisting of only . or - are removed from the MSA. ")
		fmt.Println()
		fmt.Println("1 FASTA file1 is written to the output directory.")
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
	// concatenating the space between ID and Desc with "__"
	for i := 0; i < nseqs0; i++ { //for each seqs
		seqs[i].ID = seqs[i].ID + "__" + seqs[i].Desc
	}

	// working on the sequences ----------------------------
	//Conversion of wobbles to N, and U to T
	//Removal of the - and . that are conserved vertically
	seqs = utils.RemoveMotifColumnsCores(seqs, '-', *numCores)
	seqs = utils.RemoveMotifColumnsCores(seqs, '.', *numCores)

	// for i := 0; i < len(seqs); i++ { //for each seqs
	// 	fmt.Println(seqs[i])
	// }

	// fmt.Println("---------- end ------------")
	// Open the output file
	filenameDedup := *outputDir + "/MSA_degap.fasta"
	fileDedup, err := os.Create(filenameDedup)
	if err != nil {
		panic(err)
	}
	defer fileDedup.Close()
	for i := 0; i < nseqs0; i++ { //for each seqs
		// seqs[i].ID = seqs[i].ID + "__" + seqs[i].Desc
		fmt.Fprintf(fileDedup, ">%s\n%s\n", seqs[i].ID, seqs[i].Seq)
	}
	newNberPos := len(seqs[0].Seq)
	// ............................................................
	// Reporting
	// ............................................................
	fmt.Println("=======================================================================")
	fmt.Println("(degap version: ", version, ")")
	fmt.Println("= MSA input file:\t\t\t", *inputFilePath)
	fmt.Println("= Output directory:\t\t\t", *outputDir)
	fmt.Println("=> Number of seqs:\t\t\t", nseqs0)
	fmt.Println("=> Number of initial aligned positions:\t\t", seqLen0)
	fmt.Println("=======================================================================")
	color.HiBlue("Results of degapping the MSA for columns containing only . or -:")
	fmt.Println("=> Number of remaining aligned positions:\t", newNberPos)
	PcPosRemaining := float64(newNberPos*100) / float64(seqLen0)
	roundedStr := fmt.Sprintf("%.1f", PcPosRemaining)
	fmt.Println(" (", newNberPos, "/", seqLen0, " = ", roundedStr, "% of positions remaining)")
	fmt.Println(" (the number of sequences is unchanged)")
	fmt.Println("=======================================================================")
	fmt.Println("1 FASTA file written successfully to: ", *outputDir, "/output_new_MSA_degap.fasta")
	fmt.Println("Started at: ", formattedTime)
	endTime := time.Now()
	formattedTime = endTime.Format("2006-01-02 15:04:05") // Use a specific format
	fmt.Println("Finished at:", formattedTime)
	elapsedTime := endTime.Sub(startTime)
	fmt.Println("Elapsed time", elapsedTime)
	// ............................................................

}
