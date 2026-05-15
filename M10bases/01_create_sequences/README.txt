usage:
`go run 01_create_sequences.go`
	
Creates all combinations of a 10-base-long sequence
The headers reflect the number of mutations 	
	bases {'A', 'T', 'G', 'C'}

e.g. **sequences10bases.fasta** contains:
>S0aaaaa
AAAAAAAAAA
>S1aaaab
TAAAAAAAAA
>S1aaaac
GAAAAAAAAA
>S1aaaad
CAAAAAAAAA
>S1aaaae
ATAAAAAAAA
...
where the start of the labels "Si" corresponds to the number of mutations i for the sequence "S", followed by a unique 5-letter identifier.
