# Possible workflow to analyse the NR99 MSA using the provided Golang scripts.


## 1. Get the MSA

**NR99:** SILVA_138.2_SSURef_NR99_tax_silva_full_align_trunc.fasta.gz SILVA release 138 SSU Ref NR 99 138.2 dataset is based on the full SSU Ref 138.2 dataset, in total encompassing ==510,508 sequences.
```
num_seqs sum_len min_len avg_len max_len
510,495 25,524,750,000 50,000 50,000 50,000
```

**RNA**
Number of Sequences:
20389 **Archaea** 3.99%
431166 **Bacteria** 84.46%
58940 **Eukaryota** 11.54%

**NR99:** ==SILVA_138.2_SSURef\_**NR99**\_tax_silva_full_align_trunc==.fasta.gz


## 2. trunCate
```{sh}
DESTDIR=/path/to/SILVA/NR99/output_select_seqs
l $DESTDIR
mkdir -p $DESTDIR
LOG="$DESTDIR/log.txt"


truncAte=/path/to/bin/truncAte/truncAte
$truncAte -v # version:0.3.1
$truncAte -h
$truncAte -f -i /path/to/SILVA/SILVA_138.2_SSURef_NR99_tax_silva_full_align_trunc.fasta -j 60 -o $DESTDIR > $LOG
```
<=======================================================================>
(Select_seqs version:  0.3.1 )
= MSA input file:                        /path/to/SILVA/SILVA_138.2_SSURef_NR99_tax_silva_full_align_trunc.fasta
= Output directory:                      /path/to/SILVA/NR99/output_select_seqs
=> Number of seqs:                       510495
=> Number of aligned positions:          50000
=> Threshold applied to keep a position: 0.9
 <=======================================================================>
Results:
1) After the analysis of the MSA:
    - the first position matching k is: 1144
    - the last matching position is:    41788
(this is before removing columns of only . or -)

2) After removing sequences with N or wobbles, or starting or ending with .
    - Final alignment length:            29932
    - Final number of sequences:         387633
(387633 / 510495 = 75.9 % of the initial number of sequences)
 <=======================================================================>
1 FASTA file  written successfully to:  /path/to/SILVA/NR99/output_select_seqs/output_new_MSA.fasta
Started at:  2025-03-15 15:53:02
Finished at: 2025-03-15 17:09:25
Elapsed time 1h16m22.602304458s

# 3. deduplicateseq
1) deduplicate 
2) remove the colums of . and - (includes degapping)
```{sh}
DESTDIR=/path/to/SILVA/NR99/output_select_seqs
dedup=/path/to/bin/deduplicateseq
$dedup -h
LOG2="$DESTDIR/log2_dedup.txt"
$dedup -i $DESTDIR/output_new_MSA.fasta -j 30 -o $DESTDIR -f > $LOG2
```

 <=======================================================================>
(deduplicateseq version:  0.1.1 )
= MSA input file:                        /path/to/SILVA/NR99/output_select_seqs/output_new_MSA.fasta
= Output directory:                      /path/to/SILVA/NR99/output_select_seqs
=> Number of initial seqs:                       387633
 <=======================================================================>
Results:
=> Number of final seqs:                         331663
after deduplication of the sequences
 <=======================================================================>
Files 1 FASTA file  written successfully to:  /path/to/SILVA/NR99/output_select_seqs/dedup_MSA.fasta
Started at:  2025-03-15 17:38:05
Finished at: 2025-03-15 17:48:21
Elapsed time 10m15.549396749s


**testing to see if well deduplicated**
```{sh}
seqkit rmdup -s $MSA -o test
```
[INFO] 0 duplicated records removed


### taxonomic breakdown of the cleaned MSA
```{sh}
MSA=/path/to/SILVA/NR99/output_select_seqs/dedup_MSA.fasta # also available at: https://doi.org/10.5281/zenodo.17055348
grep -c ">" $MSA # 331663

grep ">"  $MSA | cut --delim "_" -f3 | cut --delim ";" -f1 | sort | uniq -c 

R -e "round(c(6031,291195,34437)*100/331663,1)" 
```
   6031 Archaea      1.8%
 291195 Bacteria    87.8%
  34437 Eukaryota   10.4%

```{sh}
seqkit stat $MSA 
```
  num_seqs       sum_len  min_len  avg_len  max_len
  331,663  9,927,336,916   29,932   29,932   29,932


### start with Ecoli and determine its most distant relative + Archaea
=> taking then E. coli "AB035920.964.2505" as it was rather abundant


### what is the size ranges in the different domains

```{sh}
seqkit grep -n -r -p "__Archaea" $MSA | seqkit seq -g | seqkit stat
seqkit grep -n -r -p "__Eukaryota" $MSA | seqkit seq -g | seqkit stat
seqkit grep -n -r -p "__Bacteria" $MSA | seqkit seq -g | seqkit stat

```
              num_seqs      sum_len  min_len   avg_len  max_len
Archaea          6,031    7,712,784      842   1,278.9    2,787
Eukaryota       34,437   55,283,935    1,087   1,605.4    3,232
Bacteria       291,195  382,530,854    1,036   1,313.7    2,642


## 2.2. buid the axes
position of the Ecoli REF1

```{sh}
cat $MSA | grep ">" | grep -n "AB035920.964.2505" | cut -d ":" -f1
```
5553

HERE for v0.1.8
```{sh}
build_axes=/path/to/bin/build_axes

$build_axes -v
$build_axes -h

output=/path/to/SILVA/NR99/output_build_axesv0.1.8_R1_5553
mkdir -p $output
LOG=$output/log.txt
$build_axes -i $MSA -R1 5553 -stat
```
5553,AB035920.964.2505__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia;count_379,0.0121611485,0.0068819605,0.0983304207_:_seqNber,ID,Mean,SD,Max
```{sh}
$build_axes -i $MSA -R1 5553 -max

cat $MSA | grep ">" | grep -n "AJ879131.1.3425" | cut -d ":" -f1 # 149697
```

AB035920.964.2505__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia;count_379,5553,AJ879131.1.3425__Eukaryota;SAR;Rhizaria;Retaria;Foraminifera;Globothalamea;Rotaliida;Nummulitidae;Heterostegina;Heterostegina;count_1,149697,0.09833042068099496

# building 2 axes with REF1= Ecoli  REF2=Eukaryota;SAR;Rhizaria
 
 ```{sh}
output=/path/to/bin/NR99/output_build_axesv0.1.8_R1_5553_R2_149697
mkdir -p $output
LOG=$output/log.txt
$build_axes -v #version:0.1.8
$build_axes -i $MSA -R1 5553 -R2 149697 -o $output -f > $LOG
more $LOG
```
<=======================================================================>
Info:
= capasydis - version:                   0.1.8
= MSA input file:                        /path/to/SILVA/NR99/output_select_seqs/dedup_MSA.fasta
= Data written to:                       /path/to/bin/NR99/output_build_axesv0.1.8_R1_5553_R2_149697/output.csv
=> Number of sequences:                 331663
=> Number of aligned positions:         29932
=> Delta values:                         default
<=======================================================================>
Details:
REF1 name: AB035920.964.2505__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia;count_379
REF1 number: 5553
- Most distant sequence to REF1:
        - name: AJ879131.1.3425__Eukaryota;SAR;Rhizaria;Retaria;Foraminifera;Globothalamea;Rotaliida;Nummulitidae;Heterostegina;Heterostegina;count_1
        - sequence nber: 149697
        - value: 0.0983304207
--------------
REF2 name: AJ879131.1.3425__Eukaryota;SAR;Rhizaria;Retaria;Foraminifera;Globothalamea;Rotaliida;Nummulitidae;Heterostegina;Heterostegina;count_1
REF2 number: 149697
- Most distant sequence to REF2:
        - name: AB302407.1.2962__Archaea;Thermoproteota;Thermoproteia;Thermoproteales;Thermoproteaceae;Pyrobaculum;Pyrobaculum;count_1
        - sequence nber: 128973
        - value: 0.1298365572
Timing:
Started at:  2025-08-17 15:06:43
Finished at: 2025-08-17 15:07:42
Number of cores (-j): 95 
Elapsed time: 59.124874951s
<=======================================================================>

### coloring E. coli
Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;
Enterobacteriaceae: 
  - Klebsiella, 
  - Enterobacter, 
  - Citrobacter, 
  - Salmonella, 
  - Escherichia coli, 
  - Shigella, 
  - (Proteus, Serratia and other species)

```{sh}
# WD=/path/to/bin/NR99/output_build_axes_R1_5553_R2_149697
# cd $WD
# code Patterns.tsv
# cp Patterns.tsv $WD
# colorCSVTaxonomy=/path/to/bin/colorCSVTaxonomy/colorCSVTaxonomy
# $colorCSVTaxonomy -h
# $colorCSVTaxonomy -v # version:0.1.2
# $colorCSVTaxonomy -i $WD/output.csv -o $WD/output_with_color.csv -p $WD/Patterns.tsv
```
;Escherichia-Shigella; red
;Klebsiella; blue 
;Salmonella; green
;Vibrio; yellow
;Haemophilus; pink
;Serratia; magenta
;Enterobacter; brown
;Citrobacter; lightgreen
=> then go to D3.js

# edit the header

```{sh}
# code $WD/output_with_color.csv
```
x,y,label,color

## coloring the domains
```{sh}
# WD=/path/to/bin/NR99/output_build_axes_R1_5553_R2_149697
# cd $WD
# code Patterns_domain.tsv
# colorCSVTaxonomy=/path/to/bin/colorCSVTaxonomy/
# $colorCSVTaxonomy -h
# $colorCSVTaxonomy -v # version:0.1.2
# $colorCSVTaxonomy -i $WD/output.csv -o $WD/output_domains.csv -p $WD/Patterns_domain.tsv
```
8 seconds


## retrieve AB302407.1.2962

```{sh}
# seqkit grep -r -p "AB302407\.1\.2962"  /path/to/SILVA/NR99/output_select_seqs/dedup_MSA.fasta | seqkit seq -g
```

## building 3D coordinates
REF1= Ecoli    (index=5553)
REF2=Eukaryota;SAR;Rhizaria  (index=149697)
REF3=Archaea (index=128973)
  obtained when computing "output_build_axes_R1_5553_R2_149697"
      Most distant sequence to REF2:
            - name: AB302407.1.2962__Archaea;Thermoproteota;Thermoproteia;Thermoproteales;Thermoproteaceae;Pyrobaculum;Pyrobaculum;count_1
            - sequence nber: 128973
 ```{sh}
output=/path/to/bin/NR99/output_build_axesv0.1.8_R1_5553_R3_128973
mkdir -p $output
LOG=$output/log.txt
$build_axes -v #version:0.1.8
$build_axes -i $MSA -R1 5553 -R2 128973 -o $output -f > $LOG
more $LOG
head $output/output.csv
```
0.032177294,0.0668218237,HA782847.3.1866__Eukaryota;Amorphea;Obazoa;Opisthokonta;Holozoa;Choanozoa;Metazoa;Animalia;BCP;count_1
0.010367972900000001,0.0578364071,AB001521.1.1560__Bacteria;Pseudomonadota;Gammaproteobacteria;Coxiellales;Coxiellaceae;Coxiella;Ornithodoros;count_1
0.027880804000000002,0.06341756550000001,AY929353.1.1788__Eukaryota;Archaeplastida;Chloroplastida;Charophyta;Phragmoplastophyta;Streptophyta;Embryophyta;
0.0280990271,0.06384614200000001,AY929368.1.1768__Eukaryota;Archaeplastida;Chloroplastida;Charophyta;Phragmoplastophyta;Streptophyta;Embryophyta;Tracheophyta;
0.123456789|123456789 scale

## Combine the R1, R2, R3 coordinates

```{sh}
FileR1R2=/path/to/bin/NR99/output_build_axesv0.1.8_R1_5553_R2_149697/output.csv  # R1 R2
FileR1R3=/path/to/bin/NR99/output_build_axesv0.1.8_R1_5553_R3_128973/output.csv  # R1 R3

## [OK] sanity check to see if the same entries are present in both
# cat $FileR1R2 | cut -f 3 -d ","  > tmp1
# cat $FileR1R3 | cut -f 3 -d ","  > tmp2
# diff -s tmp1 tmp2 # Files tmp1 and tmp2 are identical
# rm -f tmp1; rm -f tmp2 

# aggregating the info
mkdir -p output_build_axesv0.1.8_R1_R2_R3
merge3D=/path/to/bin/merge3D/bin
$merge3D -h
$merge3D -i $FileR1R2 -j $FileR1R3 -o output_build_axesv0.1.8_R1_R2_R3/output_R1_R2_R3.csv > output_build_axesv0.1.8_R1_R2_R3/log.txt
```
2025/08/17 15:14:11 Starting CSV merge and validation...
2025/08/17 15:14:12 Reached end of both files after processing 331663 rows in each file.
2025/08/17 15:14:12 Processing complete. Wrote 331663 records to 'output_build_axesv0.1.8_R1_R2_R3/output_R1_R2_R3.csv'.
2025/08/17 15:14:12 Script finished successfully.
Timing:
Started at:  2025-03-29 18:51:24
Finished at: 2025-03-29 18:51:25
Elapsed time: 677.509251 micros


## checking uniqueness with build_axes -r  for 1e10 (default)

```{sh}
CSVfile=/path/to/bin/NR99/output_build_axesv0.1.8_R1_R2_R3/output_R1_R2_R3.csv
wc -l $CSCfile                                            # 331664
cut --delim=',' -f1  $CSVfile | sort | uniq -c | wc -l    #  331133
cut --delim=',' -f2  $CSVfile | sort | uniq -c | wc -l    #  331114
cut --delim=',' -f3  $CSVfile | sort | uniq -c | wc -l    #  330994
head  $CSVfile
```
x,y,z,label
0.032177294,0.0893223508,0.0668218237,HA782847.3.1866__Eukaryota;Amorphea;Obazoa;Opisthokonta;Holozoa;Choanozoa;Metazoa;Animalia;BCP;count_1
0.010367972900000001,0.10293579950000001,0.0578364071,AB001521.1.1560__Bacteria;Pseudomonadota;Gammaproteobacteria;Coxiellales;Coxiellaceae;Coxiella;
0.027880804000000002,0.0875629124,0.06341756550000001,AY929353.1.1788__Eukaryota;Archaeplastida;Chloroplastida;Charophyta;Phragmoplastophyta;Streptophyta;
0.0280990271,0.0874346374,0.06384614200000001,AY929368.1.1768__Eukaryota;Archaeplastida;Chloroplastida;Charophyta;Phragmoplastophyta;Streptophyta;Embryophyta;
0.123456789|123456789 scale **************


## Cleaning the taxonomy
```{sh}
WD=/path/to/bin/NR99/output_build_axesv0.1.8_R1_R2_R3
more $WD/output_R1_R2_R3.csv

cd /path/to/bin/SILVA_go
conda activate /opt/conda_envs/golang_1.24.5
go run main.go -h
go run main.go -v
go run main.go -i $WD/output_R1_R2_R3.csv -o $WD/output_R1_R2_R3_cleaned.csv -p BadNames.txt -F 7
```


