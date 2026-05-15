# Title: Simulation of 10 bases, fully mutated sequences 
# using CAPASYDIS v0.1.8 (higher precision) data preparation   
# .........................................................
# data was prepared with build_axes version 0.1.8 (e-10)
# using 10-bases completely mutated, and reversed too.

# following the scripts:  https://github.com/RametteLab/CAPASYDIS/blob/main/NR99/R%20scripts/01_Analyses_NR99_Preparation.R

# R libraries ----
library(tidyverse)

# 0) Data preparation  -----
# data colored by domains with 3 axes
D <- read.csv("output_build_axesv0.1.8_R1_R2_R3/output_R1_R2_R3_with_color.csv",
              h=TRUE) # 


options(digits = 10) # to display the same as produced in the CSV
# The issue is typically with how R displays the numbers, not how it stores them.
# options(digits = n) controls the minimum number of significant digits to print.
# options(scipen = -n) discourages R from using scientific notation for large numbers. A negative value encourages it for very small numbers.

## Uniqueness of values ---- 
# here default is e-10
Nrow=nrow(D)             
Nrow            #1048576
length(unique(D$x)) #1048508  
  round(100*length(unique(D$x))/Nrow,2) #99.99
length(unique(D$y)) #1048476  
  round(100*length(unique(D$y))/Nrow,2) #99.99
length(unique(D$z)) #1048469
  round(100*length(unique(D$z))/Nrow,2) #99.99

#unique coordinate
D %>% select(x,y,z) %>% unique() %>% nrow() #1048576   100%
D %>% select(x,y) %>% unique() %>% nrow()   #1048576   100%

D[duplicated(D[,c("x", "y","z")]), ] #%>% arrange(desc(x),desc(y),desc(z)) %>% head()
# all unique!!

# adding the number of mutations
D1 <- D %>% 
  mutate(
    prefix = str_sub(label, 1, 2),
    Nmut  = str_remove(prefix, "^S")
  )
D1 <- D1 %>%
  mutate(
    prefix = str_sub(label, 1, 2),
    Nmut  = case_when(
      prefix == "S_" ~ 10L,
      TRUE           ~ as.numeric(str_remove(prefix, "^S"))
    )
)

table(D1$Nmut)
round(100*table(D1$Nmut)/nrow(D1),2)
#   0      1      2      3      4      5      6      7      8      9     10 # Nmut
#   1     30    405   3240  17010  61236 153090 262440 295245 196830  59049 # nber of cases
# 0.00  0.00   0.04   0.31   1.62   5.84  14.60  25.03  28.16  18.77   5.63 # %

D1 <- D1 %>% select(-prefix)
head(D1)


############################################################### 3' -> 5' sequence reading
Drev <- read.csv("output_build_axesv0.1.8_R1_R2_R3_Rev/output_R1_R2_R3_with_color.csv",
                 h=TRUE) # 


## Uniqueness of values ---- 
# here default is e-10
NrowRev=nrow(Drev)             
NrowRev            #1048576
length(unique(Drev$x)) #1048508  
round(100*length(unique(Drev$x))/NrowRev,2) #99.99
length(unique(Drev$y)) #1048476  
round(100*length(unique(Drev$y))/NrowRev,2) #99.99
length(unique(Drev$z)) #1048469
round(100*length(unique(Drev$z))/NrowRev,2) #99.99

#unique coordinate
Drev %>% select(x,y,z) %>% unique() %>% nrow() #1048576   100%
Drev %>% select(x,y) %>% unique() %>% nrow()   #1048576   100%

Drev[duplicated(D[,c("x", "y","z")]), ] #%>% arrange(desc(x),desc(y),desc(z)) %>% head()
# all unique!!

# adding the number of mutations
D1rev <- Drev %>% 
  mutate(
    prefix = str_sub(label, 1, 2),
    Nmut  = str_remove(prefix, "^S")
)
D1rev <- D1rev %>%
  mutate(
    prefix = str_sub(label, 1, 2),
    Nmut  = case_when(
      prefix == "S_" ~ 10L,
      TRUE           ~ as.numeric(str_remove(prefix, "^S"))
    )
)

round(100*table(D1$Nmut)/nrow(D1),2)
#   0      1      2      3      4      5      6      7      8      9     10 # Nmut
#   1     30    405   3240  17010  61236 153090 262440 295245 196830  59049 # nber of cases
# 0.00  0.00   0.04   0.31   1.62   5.84  14.60  25.03  28.16  18.77   5.63 # %

D1rev <- D1rev %>% select(-prefix)
head(D1rev)

rm(list = ls()[ ! ls() %in% c("D1", "D1rev") ]) # keeping D1 and D1rev
DATE <- format(Sys.time(), "%Y%m%d")
save.image(paste0(DATE,"_data.RData"))
