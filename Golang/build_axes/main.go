package main

import (
	"bufio"
	"build_axes/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/fatih/color"
)

func main() {
	version := "0.1.8"
	inputFilePath := flag.String("i", "data_test/inputMSA1.fasta", "Path to the input FASTA file")
	numCores := flag.Int("j", runtime.NumCPU()-1, "Number of CPU cores to use (default all available on the machine -1)")
	outputDir := flag.String("o", "output_MSA", "Output directory")
	R1 := flag.Int("R1", 1, "index of the sequence to be used as reference 1")
	R2 := flag.Int("R2", 0, "index of the sequence to be used as reference 2. Indicate \"0\" ")
	precision := flag.Float64("r", 0.000_000_000_1, "how many decimals to report in the results. The smallest value that can be reported has precision r: 0.000_000_01 (i.e. 1e-08)")
	Maxonly := flag.Bool("max", false, "if set, only the most divergent sequence to -R1 is calculated and reported as a csv output. The rest of the coordinates are not reported. The order of the values is \"REF1 name, REF1 sequence number, Name of most distant sequence, Number of most distant sequence, (max) asymdist value\"")
	Statonly := flag.Bool("stat", false, "if set, only the statistics (mean, sd, max) asymdist values using -R1 as reference are calculated")
	PositionsFile := flag.String("p", "", "Path to the input file containing the positions to use for R1 (one per line). This works in combination with \"-stat\" only") //"data_test/PosFile.txt"
	Statall := flag.Bool("statall", false, "if set, the statistics (mean, sd, max) of asymdist values for all the sequences are calculated.")
	forceOutput := flag.Bool("f", false, "Force output (overwrite existing files)")
	vers := flag.Bool("v", false, "Version of the software)")
	help := flag.Bool("h", false, "Show help message")
	// Define flag for DeltaVector file path
	deltaVectorFilePath := flag.String("d", "", "Path to the DeltaVector file. If not specified, these default values are used for calculations:\nAG: 0.01 GA: 0.01 CT: 0.02 TC: 0.02 AC: 0.06 CA: 0.06 GT: 0.07 TG: 0.07 AT: 0.08 TA: 0.08\nGC: 0.09 CG: 0.09 A-: 0.11 -A: 0.11 T-: 0.12 -T: 0.12 G-: 0.13 -G: 0.13 C-: 0.14 -C: 0.14")
	// Use default DeltaVector if the flag is not set
	var DeltaVector map[string]float64
	if *deltaVectorFilePath == "" {
		DeltaVector = map[string]float64{
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
	} else {
		// Read DeltaVector from file if the flag is set
		var err error
		DeltaVector, err = utils.ReadDeltaVectorFromFile(*deltaVectorFilePath)
		if err != nil {
			panic(err)
		}
	}

	flag.Parse() // Parse the flags

	// ............................................................
	if *R1 <= 0 {
		log.Fatalln("[Error] Wrong input for -R1. The value should be a positive integer!")
	}
	if *R2 < 0 {
		log.Fatalln("[Error] Wrong input for -R2. The value should be set to 0 or a positive integer.")
	}
	if *Maxonly || *Statonly {
		*R2 = 0 // swich off -R2

	}
	if *vers {
		fmt.Print("version:")
		color.Red(version)
		return // Exit the program
	}

	if *help {
		color.Red("----------------------------------\n")
		color.Red("          capasydis (go version)\n")
		color.Red("----------------------------------\n")
		fmt.Print("version: ")
		color.Red(version + "\n")

		fmt.Println("Author: A. Ramette, 2025")
		fmt.Println()
		flag.PrintDefaults() // Print the help message with default values
		fmt.Println("Usage: go run main.go -i data_test/inputMSA1.fasta -o outputMSA -R1 1 -R2 2 -f")

		return // Exit the program
	}
	// ............................................................
	startTime := time.Now()
	formattedTime := startTime.Format("2006-01-02 15:04:05") // Use a specific format

	// ............................................................
	// Set the number of CPU cores to use
	runtime.GOMAXPROCS(*numCores)
	if !*Maxonly && !*Statonly && !*Statall {
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
	seqLen := len(seqs[0].Seq)
	// number of sequences at the beginning
	nseqs := len(seqs) //number of sequences (rows)
	if *R1 > nseqs {
		log.Fatalln("[Error] Wrong input for -R1. The value should be below the total number of sequences in the MSA of:", nseqs)
	}
	if *R2 > nseqs {
		log.Fatalln("[Error] Wrong input for -R2. The value should be below the total number of sequences in the MSA of:", nseqs)
	}
	// ---------------------------------------------
	if *PositionsFile != "" && *Statonly {
		// fmt.Println("Position file detected and -stat flag too")
		// open and validate the file (only 1 column expected) then

		filePos, err := os.Open(*PositionsFile)
		if err != nil {
			fmt.Printf("Error: Input file '%s' does not exist.\n", *PositionsFile)
			log.Fatal(err)
		}
		defer filePos.Close()

		numbers := []int64{}
		scanner := bufio.NewScanner(filePos)
		for scanner.Scan() {
			line := scanner.Text()
			number, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				fmt.Println("[error] in reading the file: ", *PositionsFile)
				fmt.Println("        only 1 number per line is expected.")
				log.Fatal(err)
			}
			numbers = append(numbers, number)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// fmt.Println(numbers)
		// then loop on each provided position as done in *Statall
		for _, Rval := range numbers {
			// fmt.Println("-----------------: REF pos", Rval)

			REFi := string(seqs[Rval-1].Seq)
			// fmt.Println(REF)
			LessDivergScorei := utils.CalculateLessDivergentScore(REFi)
			MostDivergScorei := utils.CalculateMostDivergentScore(REFi, DeltaVector)

			var wg sync.WaitGroup
			wg.Add(nseqs)
			Ai := make([]float64, nseqs)

			for j := 0; j < nseqs; j++ {
				go func(j int) {
					defer wg.Done()
					Ai[j] = utils.AsymDist1Seq(REFi, string(seqs[j].Seq), LessDivergScorei, MostDivergScorei, DeltaVector)
				}(j)
			}
			wg.Wait()

			AiMean, _ := utils.Mean(Ai)
			AiSD, _ := utils.SD(Ai)
			AiMax, _ := utils.Max(Ai)

			AiMean = utils.RoundNber(AiMean, *precision)
			AiSD = utils.RoundNber(AiSD, *precision)
			AiMax = utils.RoundNber(AiMax, *precision)
			fmt.Printf("%v,%s,%g,%g,%g_:_seqNber,ID,Mean,SD,Max\n", Rval, seqs[Rval-1].ID, AiMean, AiSD, AiMax)

		}
		os.Exit(0)
	}
	// ---------------------------------------------
	if *Statall {
		for i := 0; i < nseqs; i++ {
			REFi := string(seqs[i].Seq)
			// fmt.Println(REF)
			LessDivergScorei := utils.CalculateLessDivergentScore(REFi)
			MostDivergScorei := utils.CalculateMostDivergentScore(REFi, DeltaVector)

			var wg sync.WaitGroup
			wg.Add(nseqs)
			Ai := make([]float64, nseqs)

			for j := 0; j < nseqs; j++ {
				go func(j int) {
					defer wg.Done()

					Ai[j] = utils.AsymDist1Seq(REFi, string(seqs[j].Seq), LessDivergScorei, MostDivergScorei, DeltaVector)
				}(j)
			}
			wg.Wait()

			AiMean, _ := utils.Mean(Ai)
			AiSD, _ := utils.SD(Ai)
			AiMax, _ := utils.Max(Ai)

			AiMean = utils.RoundNber(AiMean, *precision)
			AiSD = utils.RoundNber(AiSD, *precision)
			AiMax = utils.RoundNber(AiMax, *precision)

			fmt.Printf("%v,%s,%g,%g,%g_:_seqNber,ID,Mean,SD,Max\n", i+1, seqs[i].ID, AiMean, AiSD, AiMax)

		}
		os.Exit(0)

	} else {

		// loop through the MSA using embarrasingly parallel implementation (no data sync needed in-between)
		// fmt.Println("----------")
		REF1 := string(seqs[*R1-1].Seq)
		// fmt.Println(REF)
		LessDivergScore1 := utils.CalculateLessDivergentScore(REF1)
		MostDivergScore1 := utils.CalculateMostDivergentScore(REF1, DeltaVector)

		var LessDivergScore2 float64
		var MostDivergScore2 float64
		var REF2 string
		if *R2 != 0 {
			REF2 = string(seqs[*R2-1].Seq)
			LessDivergScore2 = utils.CalculateLessDivergentScore(REF2)
			MostDivergScore2 = utils.CalculateMostDivergentScore(REF2, DeltaVector)
		}

		var wg sync.WaitGroup
		wg.Add(nseqs)
		A1 := make([]float64, nseqs)

		A2 := make([]float64, nseqs)

		for i := 0; i < nseqs; i++ {
			go func(i int) {
				defer wg.Done()
				A1[i] = utils.AsymDist1Seq(REF1, string(seqs[i].Seq), LessDivergScore1, MostDivergScore1, DeltaVector)
				if *R2 != 0 {
					A2[i] = utils.AsymDist1Seq(REF2, string(seqs[i].Seq), LessDivergScore2, MostDivergScore2, DeltaVector)
				}
			}(i)
		}
		wg.Wait()
		// fmt.Println(A2)
		// Open the output file only if max not set

		outputFile := filepath.Join(*outputDir, "output.csv")

		if !*Maxonly && !*Statonly {

			file, err := os.Create(outputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			// fmt.Fprintln(file, "x,y,Label")
			// Write the data to the file
			for i := range A1 {
				A1[i] = utils.RoundNber(A1[i], *precision)
				if *R2 != 0 {
					A2[i] = utils.RoundNber(A2[i], *precision)
					_, err := fmt.Fprintf(file, "%g,%g,%s\n", A1[i], A2[i], seqs[i].ID)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					_, err := fmt.Fprintf(file, "%g,%s\n", A1[i], seqs[i].ID)
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
		if *Statonly { // do minimal reporting as csv values to stdout
			A1Mean, _ := utils.Mean(A1)
			A1Mean = utils.RoundNber(A1Mean, *precision)
			A1SD, _ := utils.SD(A1)
			A1SD = utils.RoundNber(A1SD, *precision)
			A1Max, _ := utils.Max(A1)
			A1Max = utils.RoundNber(A1Max, *precision)
			fmt.Printf("%v,%s,%g,%g,%g_:_seqNber,ID,Mean,SD,Max\n", *R1, seqs[*R1-1].ID, A1Mean, A1SD, A1Max)

			os.Exit(0)
		}
		if *Maxonly { // do minimal reporting as csv values to stdout
			max1_index := 0
			max1_value := A1[0]
			for i, value := range A1 {
				if value > max1_value {
					max1_index = i
					max1_value = value
				}
			}
			fmt.Printf("%s,%v,%s,%d,", seqs[*R1-1].ID, *R1, seqs[max1_index].ID, max1_index+1)
			fmt.Println(max1_value)

		} else {
			fmt.Println("=======================================================================")
			color.HiBlue("Info:")
			fmt.Println("= capasydis - version:\t\t\t", version)
			fmt.Println("= MSA input file:\t\t\t", *inputFilePath)
			fmt.Println("= Data written to:\t\t\t", outputFile)
			fmt.Print("=> Number of sequences:\t\t\t")
			color.Red("%d", nseqs)
			fmt.Print("=> Number of aligned positions:\t\t")
			color.Red("%d", seqLen)
			if *deltaVectorFilePath == "" {
				fmt.Println("=> Delta values:\t\t\t", "default")
			} else {
				fmt.Println("=> Delta values from:\t\t\t", *deltaVectorFilePath)
			}
			// ............................................................
			// commenting on the max values
			fmt.Println("=======================================================================")
			color.HiBlue("Details:")
			fmt.Println("REF1 name:", string(seqs[*R1-1].ID))
			fmt.Println("REF1 number:", *R1)
			fmt.Println("- Most distant sequence to REF1:")
			max1_index := 0
			max1_value := A1[0]
			// fmt.Println("max1_value start:", max1_value)
			for i, value := range A1 {
				if value > max1_value {
					max1_index = i
					max1_value = value
				}
			}

			fmt.Println("\t- name:", seqs[max1_index].ID)
			fmt.Println("\t- sequence nber:", max1_index+1)
			fmt.Println("\t- value:", max1_value)
			// fmt.Println("Max1 seq:", seqs[max1_index].Seq)
			if *R2 != 0 {
				fmt.Println("--------------")
				fmt.Println("REF2 name:", string(seqs[*R2-1].ID))
				fmt.Println("REF2 number:", *R2)
				fmt.Println("- Most distant sequence to REF2:")
				max2_index := 0
				max2_value := A2[0]
				// fmt.Println("max1_value start:", max1_value)
				for i, value := range A2 {
					if value > max2_value {
						max2_index = i
						max2_value = value
					}
				}

				fmt.Println("\t- name:", seqs[max2_index].ID)
				fmt.Println("\t- sequence nber:", max2_index+1)
				fmt.Println("\t- value:", max2_value)
			}
			// fmt.Println("=======================================================================")
			color.HiBlue("Timing:")
			// fmt.Println("=======================================================================")
			fmt.Println("Started at: ", formattedTime)
			endTime := time.Now()
			formattedTime = endTime.Format("2006-01-02 15:04:05") // Use a specific format
			fmt.Println("Finished at:", formattedTime)
			elapsedTime := endTime.Sub(startTime)
			red := "\033[31m"
			reset := "\033[0m" // Reset color
			fmt.Printf("Number of cores (-j): %d \n", *numCores)
			fmt.Printf("Elapsed time: %s%v%s\n", red, elapsedTime, reset)
			fmt.Println("=======================================================================")
		}
	}
}
