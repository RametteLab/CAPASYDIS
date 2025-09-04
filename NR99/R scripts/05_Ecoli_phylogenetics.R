# title: NR99 E. coli - phylogenetics
# to be run after the "Analyses_NR99_3D_Ecoli_nFriends.R"
# Aim: compare the asymdist to classical phylogenetic distances for E. coli
  
#-------------------------------------------------------------------
# the first (commented) part below was done on a Linux cluster
  
# # extract all the E.coli to see why there are 2 groups using capasydis
#   
#   ```{sh}
# # MSA after applying capasydis_go, and dedup
# MSA=/path/to/file/dedup_MSA.fasta # also available at:  https://doi.org/10.5281/zenodo.17055348
# cd /home/aramette/projects/DISCOMP/NR99/Ecoli_phylogenetics/2025-08
# grep -c ">"  $MSA #331663
# grep ">"  $MSA | grep "Escherichia;"  | cut -d ";" -f7 | sort | uniq -c  # 948 Escherichia
# 
# # extracting E. coli sequences
# seqkit grep -n -r -p "Escherichia;" $MSA | seqkit seq -g > output_new_MSA_Ecoli_unaligned.fasta
# seqkit grep -n -r -p "Escherichia;" $MSA > output_new_MSA_Ecoli_aligned.fasta
# 
# seqkit stat output_new_MSA_Ecoli_unaligned.fasta
# seqkit stat output_new_MSA_Ecoli_aligned.fasta
# ```
# file                                  num_seqs     sum_len  min_len  avg_len  max_len
# output_new_MSA_Ecoli_unaligned.fasta       948   1,257,566    1,130  1,326.5    1,509
# output_new_MSA_Ecoli_aligned.fasta         948  28,375,536   29,932   29,932   29,932
# 
# # 2. phylogenetic tree computation
# ```{sh}
# conda activate iqtree_2.1.2  
# iqtree -s  output_new_MSA_Ecoli_aligned.fasta -pre iqtree   -m TIM2+F+I
# ```
# Input data: 948 sequences with 29932 nucleotide sites
# Number of constant sites: 28830 (= 96.3183% of all sites)
# Number of invariant (constant or ambiguous constant) sites: 28830 (= 96.3183% of all sites)
# Number of parsimony informative sites: 416
# Number of distinct site patterns: 1521
# 
# SUBSTITUTION PROCESS
# --------------------
#   
#   Model of substitution: TIM2+F+I
# 
# Rate parameter R:
#   
#   A-C: 1.2424
# A-G: 2.0125
# A-T: 1.2424
# C-G: 1.0000
# C-T: 2.3321
# G-T: 1.0000
# 
# State frequencies: (empirical counts from alignment)
# 
# pi(A) = 0.2501
# pi(C) = 0.244
# pi(G) = 0.2712
# pi(T) = 0.2348
# 
# Rate matrix Q:
#   
# A    -1.039    0.2762    0.4974    0.2658
# C    0.2832    -1.029    0.2471     0.499
# G    0.4587    0.2224    -0.895     0.214
# T    0.2832    0.5185    0.2471    -1.049
#-------------------------------------------------------------------
# moving to R local
# data.frame(DeltaMatrix)
# 
# RateMatrixIqTree <- t(data.frame(
#   A =  c( -1.039,    0.2762 ,   0.4974  ,  0.2658),
#  C =   c(0.2832 ,   -1.029 ,   0.2471 ,    0.499),
#  G =   c(0.4587 ,   0.2224 ,   -0.895  ,   0.214),
#  T =   c(0.2832 ,   0.5185 ,   0.2471  ,  -1.049)
# ))
# eRateMatrixIqTree <- exp(RateMatrixIqTree) #https://github.com/orgs/iqtree/discussions/202


library(dplyr)
library(cultevo) # for Hamming distance calculations
library(ape)
library(gridExtra)
library(ggalign)
library(patchwork)
library(ggtree)# BiocManager::install("ggtree", version = "3.10")

# 1) Number of nt differences in MSA vs. asymdist ----
# need to obtain the number of differences between the E. coli sequences
# Import the aligned FASTA file
MSA_Ecoli <- read.dna(file="data/iqtree/output_new_MSA_Ecoli_aligned.fasta", 
                      format = "fasta",as.character =  TRUE)
MSA_Ecoli <- apply(MSA_Ecoli,2,toupper)
# print(MSA_Ecoli)
MSA_Ecoli_labels <- as.character(attr(MSA_Ecoli,"dimnames")[[1]] )
attr(MSA_Ecoli,"dimnames")[[1]] <- paste0("S",1:length(MSA_Ecoli_labels))
colnames(MSA_Ecoli) <- paste0("P",1:ncol(MSA_Ecoli))

# small test: 
# table(MSA_Ecoli[1,]==MSA_Ecoli[2,])
# MSA_Ecoli_test <- MSA_Ecoli[1:2,which (MSA_Ecoli[1,]!=MSA_Ecoli[2,])] # alignment of differences
## P66 P70 P525 P527 P534 P3125 P3136 P23658 P23665
## S1 "C" "-" "G"  "T"  "C"  "T"   "T"   "A"    "G"   
## S2 "G" "A" "T"  "C"  "T"  "A"   "A"   "G"    "T"
# cultevo::hammingdists(MSA_Ecoli_test) == ncol(MSA_Ecoli_test) ## 9 # as expected

MSA_Ecoli <- read.dna(file="data/iqtree/output_new_MSA_Ecoli_aligned.fasta", 
                      format = "fasta",as.character = TRUE)
# View(MSA_Ecoli)
# count of sequence differences
# MSA_raw_dist <- cultevo::hammingdists(MSA_Ecoli) # rather slow
# saveRDS(MSA_raw_dist,"data/MSA_raw_dist.RDS")
MSA_raw_dist <- readRDS("data/MSA_raw_dist.RDS")
MSA_raw_dist_matrix <- as.matrix(MSA_raw_dist) # full matrix for easier viewing
# dim(MSA_raw_dist_matrix) #948 948


# Extract the values from the lower triangle (or upper triangle) of the matrix
# Put the values into a data frame for ggplot2
dist_df <- data.frame(Ndifferences = MSA_raw_dist)
H1 <- ggplot(dist_df, aes(x = Ndifferences)) + scale_x_continuous()+
  geom_histogram(binwidth = 0.5, fill = "skyblue", color = "black") +
  labs(title = "",
       x = "Number of nt differences in MSA",
       y = "Frequency") +
  theme_minimal()
H2 <- ggplot(dist_df, aes(x = Ndifferences)) + scale_x_continuous(limits = c(0,60))+
  geom_histogram(binwidth = 0.5, fill = "skyblue", color = "black") +
  theme_classic()+ xlab("")+ylab("")

png("plots/EcolinFriends/Hist_Nber_MSAdiff.png",width = 20,height = 30,units="cm",res=300)
    H1 + inset_element(H2, left = 0.25, bottom = 0.4, right = 0.95, top = 0.9)
dev.off()

# table((dist_df))
summary(dist_df) # how differnet are the sequences
  # Ndifferences      
  # Min.   :  1.00000  
  # 1st Qu.:  9.00000  
  # Median : 15.00000  
  # Mean   : 22.06438  
  # 3rd Qu.: 22.00000  
  # Max.   :385.00000  

# Checking: Are the sequence order in the MSA and in the asymdist comparable?
# table(MSA_Ecoli_labels == Ecoli2DLabel$label)
# yes

# for the sequences with one SNP differences, are the nt differences the same pairs?
length(which(dist_df==1)) #655 OK

N <- nrow(MSA_raw_dist_matrix)
MSA_raw_dist_matrix1diff <- MSA_raw_dist_matrix
for(i in 1:N){ # to remove upper diag
  for (j in i:N){
    MSA_raw_dist_matrix1diff[i,j] <- NA
  }
}

colnames(MSA_raw_dist_matrix1diff) <- paste0("S",1:length(MSA_Ecoli_labels))
rownames(MSA_raw_dist_matrix1diff) <- paste0("S",1:length(MSA_Ecoli_labels))
# MSA_raw_dist_matrix1diff[1:3,1:3]

ListOfPairs1SNP <- which(MSA_raw_dist_matrix1diff == 1, arr.ind = TRUE)
ListOfPairs1SNP <- data.frame(ListOfPairs1SNP)
ListOfPairs1SNP$Pairs <- NA
rownames(ListOfPairs1SNP) <- 1:nrow(ListOfPairs1SNP)

for (i in 1: nrow(ListOfPairs1SNP)){ #
  X <- MSA_Ecoli[as.numeric(ListOfPairs1SNP[i,1:2]),]
  ListOfPairs1SNP$Pairs[i] <- paste0(X[,which (X[1,]!=X[2,])],sep="",collapse = "") # alignment of differences
}

Table1SNP <- ListOfPairs1SNP$Pairs %>% table()
Table1SNPpairs <- names(Table1SNP)
Table1SNPPcent <- round(Table1SNP*100/655,1)

T1SNP <- data.frame(Table1SNPpairs=Table1SNPpairs,
                    Count=as.numeric(Table1SNP),
                    Pcent=as.numeric(Table1SNPPcent)
)

# T1SNP %>% arrange(desc(Count))
#         N     %            N     %
#   TC   134  20.5 (+) CT    60   9.2
#   AG   126  19.2 (+) GA    92  14.0
#   -G    59   9.0 (+) G-    15   2.3
#   AC    42   6.4 (+) CA    21   3.2
             
# ..................................


# Difference in nt number vs CAPASYDIS distances between points ----
# also captured by the  seen in CAPADYDIS?

# unfolding the matrices to compare the values
DF_nberDiff_asymdist <- data.frame(
  MSA_raw_dist = MSA_raw_dist ,
  Ecoli1D_dist = dist(Ecoli1D$asymDist1D) ,
  Ecoli2D_dist = dist(Ecoli2D[,1:2]), 
  Ecoli3D_dist = dist(Ecoli3D) 
)
# n= nrow(Ecoli1D)
# n*(n-1) /2 =448878
# nrow(DF_nberDiff_asymdist)==n*(n-1) /2 

# 1SNP diff: Compare matrices 1 seq. differences with CAPASYDIS ----
# verification that we indeed unique coordinates for each sequence using CAPASYDIS
  # Ecoli1D$asymDist1D %>% length()
  # Ecoli1D$asymDist1D %>%unique() %>% length() # 946 at e-10
  # Ecoli2D[,1:2] %>% unique() %>% nrow() # 948 at e-10
  # Ecoli3D %>% unique() %>% nrow() # 948 at e-10

# Are the distances also unique? But why should they be unique? What matters is that the "points" are unique in the seqverses!
DF_nberDiff_asymdist_1SNP <- DF_nberDiff_asymdist %>% filter(MSA_raw_dist==1)
  DF_nberDiff_asymdist_1SNP %>% nrow #655
  DF_nberDiff_asymdist_1SNP %>% select(Ecoli1D_dist) %>% unique() %>% nrow() #441 67.3%
  DF_nberDiff_asymdist_1SNP %>% select(Ecoli1D_dist,Ecoli2D_dist) %>% unique() %>% nrow() #516 78.8%
  DF_nberDiff_asymdist_1SNP %>% select(Ecoli1D_dist,Ecoli2D_dist,Ecoli3D_dist) %>% unique() %>% nrow() #585 89.3%


# Compare matrices CAPASYDIS 1D, 2D, 3D, vs. nber of seq. differences ----
## Asymdist 1D vs. MSA_Ndist ----
# Ecoli1D %>% head()
# Ecoli1D %>% select(Aboveq99) %>% table()

cor.test(x=DF_nberDiff_asymdist$Ecoli1D_dist,y=DF_nberDiff_asymdist$MSA_raw_dist) # 0.6434485  P<0.001
cor.test(x=DF_nberDiff_asymdist$Ecoli2D_dist,y=DF_nberDiff_asymdist$MSA_raw_dist) # 0.6395944  P<0.001
cor.test(x=DF_nberDiff_asymdist$Ecoli3D_dist,y=DF_nberDiff_asymdist$MSA_raw_dist) # 0.6358163  P<0.001
                         
E1D_1 <- ggplot(DF_nberDiff_asymdist,aes(x=MSA_raw_dist,y=Ecoli1D_dist)) +
  theme_bw() + geom_point(size=1)+ geom_smooth()+scale_x_continuous()+scale_y_continuous()+
  labs(x="Number of nt differences between sequence pairs",y="distances between 1D coordinates")

# E1D_2 <- ggplot(DF_nberDiff_asymdist,aes(x=MSA_raw_dist)) +
#   geom_histogram(bins = 500)+
#   # geom_density(alpha=0.5,colour="black") +
#   theme_bw() +scale_x_continuous()+
#   theme0(plot.margin = unit(c(1,0,-0.48,2.2),"lines"))
# 
# E1D_3 <- ggplot(DF_nberDiff_asymdist,aes(x=Ecoli1D_dist) )+
#   # geom_density(alpha=0.5,colour="black") +
#   geom_histogram(bins = 100000)+
#   coord_flip()  +scale_x_continuous()+
#   # scale_x_continuous(labels = NULL,breaks=NULL,expand=c(0.02,0)) +
#   theme_bw() +
#   theme0(plot.margin = unit(c(0,1,1.2,-0.2),"lines"))

# png("plots/EcolinFriends/MargPlot_Ndiff_Asymdist1D.png",width = 30,height = 20,units="cm",res=300)
# E1D_1
# # grid.arrange(arrangeGrob(E1D_2,ncol=2,widths=c(3,1)),
# #              arrangeGrob(E1D_1,E1D_3,ncol=2,widths=c(3,1)),
# #              heights=c(1,3))
# dev.off()

# Create legend
# E1D_text <- c("Figure MargPlot_Ndiff_Asymdist1D. Using method = 'gam' and formula = 'y ~ s(x, bs = \"cs\")",
# "")
# writeLines(E1D_text, con = "plots/EcolinFriends/Plot_Ndiff_Asymdist1-3D_Legend.txt")


E2D_1 <- ggplot(DF_nberDiff_asymdist,aes(x=MSA_raw_dist,y=Ecoli2D_dist)) +
  theme_bw() + geom_point(size=1)+ geom_smooth()+scale_x_continuous()+scale_y_continuous()+
  labs(x="Number of nt differences between sequence pairs",y="distances between 2D coordinates")

# png("plots/EcolinFriends/MargPlot_Ndiff_Asymdist2D.png",width = 30,height = 20,units="cm",res=300)
#   E2D_1
# dev.off()

E3D_1 <- ggplot(DF_nberDiff_asymdist,aes(x=MSA_raw_dist,y=Ecoli3D_dist)) +
  theme_bw() + geom_point(size=1)+ geom_smooth()+scale_x_continuous()+scale_y_continuous()+
  labs(x="Number of nt differences between sequence pairs",y="distances between 3D coordinates")

# png("plots/EcolinFriends/MargPlot_Ndiff_Asymdist3D.png",width = 30,height = 20,units="cm",res=300)
# E3D_1
# dev.off()


png("plots/EcolinFriends/Plot_Ndiff_Asymdist1-3D.png",width = 20,height = 30,units="cm",res=300)
  E1D_1 + E2D_1 + E3D_1+plot_layout(ncol=1)
dev.off()


# 2) Comparison with "classical" phylogenetic analysis? ----
# ===> do this for the sequence with 1 SNP only!!!

library(ape)
# Read the mldist file
mldist_df <- read.table("data/iqtree/iqtree.mldist", header = FALSE, skip = 1, 
                 stringsAsFactors = FALSE)
# mldist_df[1:3,1:2]

# Extract the sequence names
mldist_df_seq_names <- mldist_df[, 1]
mldist_df_seq_names_short <- as.character(sapply(mldist_df_seq_names,function(x) strsplit(x,"__")[[1]][1]))

# Create a distance matrix
mldist_dist_matrix <- as.matrix(as.dist(mldist_df[, -1])) # the first columns are the labels
# mldist_dist_matrix[1:2,1:2]
# mldist_dist_matrix %>% dim  # 948 948

options(digits = 10) # to display the same as produced in the CSV
# using the info above about where the 1SNP is, extract the corresponding ML distances

ListOfPairs1SNP$mldist <- NA
for (i in 1: nrow(ListOfPairs1SNP)){ 
  ListOfPairs1SNP$mldist[i] <-mldist_dist_matrix[as.numeric(ListOfPairs1SNP[i,1]),as.numeric(ListOfPairs1SNP[i,2])]
}

# now do the same for the Asymdist
# Ecoli3D # has the x,y,z coordinates
# using ListOfPairs1SNP get the lines, => calculate their distance

ListOfPairs1SNP$dist3D <- NA
ListOfPairs1SNP$dist2D <- NA
ListOfPairs1SNP$dist1D <- NA

for (i in 1: nrow(ListOfPairs1SNP)){ 
  ListOfPairs1SNP$dist3D[i] <- dist(Ecoli3D[as.numeric(ListOfPairs1SNP[i,1:2]),1:3])
  ListOfPairs1SNP$dist2D[i] <- dist(Ecoli3D[as.numeric(ListOfPairs1SNP[i,1:2]),1:2])
  ListOfPairs1SNP$dist1D[i] <- dist(Ecoli3D[as.numeric(ListOfPairs1SNP[i,1:2]),1])
}
ListOfPairs1SNP %>% head()
ListOfPairs1SNP %>% nrow() #655
ListOfPairs1SNP %>% filter(mldist>0.00001) %>% arrange(Pairs) %>% nrow() #530
ListOfPairs1SNP %>% filter(mldist>0.0001) %>% arrange(Pairs) %>% head()

ML1 <- ListOfPairs1SNP %>% ggplot(aes(x=mldist,y=dist1D))+geom_point()+theme_minimal()+xlab("")+ylab("distances between asymdist 1D ")
ML2 <- ListOfPairs1SNP %>% ggplot(aes(x=mldist,y=dist2D))+geom_point()+theme_minimal()+xlab("")+ylab("distances between asymdist 2D ")
ML3 <- ListOfPairs1SNP %>% ggplot(aes(x=mldist,y=dist3D))+geom_point()+theme_minimal()+ylab("distances between asymdist 3D ")

png("plots/EcolinFriends/mlstdist_Asymdist.png",width = 20,height = 20,units="cm",res=300)
  ML1+ML2+ML3+plot_layout(ncol=1)
dev.off()

ListOfPairs1SNP$mldist %>% unique() %>% length()# how many unique values there are
ListOfPairs1SNP$mldist %>% length()


# trees ----
mldist_dist_matrix %>% head()
# Assign sequence names to the matrix
attr(mldist_dist_matrix, "Labels") <- mldist_df_seq_names_short
hc <-hclust(as.dist(mldist_df[, -1])) 

png("plots/EcolinFriends/Ecoli_phylogram.png",width = 20,height = 30,units="cm",res=300)
plot(as.phylo(hc), type = "phylogram", show.tip.label = TRUE,cex = 0.1, 
     edge.color = "black", edge.width = 1, edge.lty = 1,  
     tip.color = "black")
dev.off()

# Unrooted
png("plots/EcolinFriends/Ecoli_unrooted.png",width = 30,height = 20,units="cm",res=300)
plot(as.phylo(hc), type = "unrooted", cex = 0.2,
     no.margin = TRUE)
dev.off()


png("plots/EcolinFriends/Ecoli_mlstdist.png",width = 20,height = 20,units="cm",res=300)
hist(mldist_dist_matrix,col="darkgrey",xlab="ML distances")    
dev.off()

# How do the distances compare?
## ref1 in 3D CAPASYDIS axes
which (mldist_df_seq_names_short=="AB035920.964.2505")# 28

# build a tree with those 1SNP Ecoli seqs

#finding the rows and columns to be highlighted
Ids_SNPSeqs <- ListOfPairs1SNP %>% select(row,col) %>% pull() %>% unique() # length() #200 sequences

ColorsSeqs <- rep("black",)

mldist_dist_matrix_1SNPseqs <- mldist_dist_matrix[Ids_SNPSeqs,Ids_SNPSeqs]
mldist_dist_matrix_1SNPseqs %>% dim() #200 200



hc1 <-hclust(as.dist(mldist_dist_matrix_1SNPseqs)) 
hc1$labels <- attr(mldist_dist_matrix,"Labels")[Ids_SNPSeqs]
unrooted <- ggtree(hc1,layout ="ape")+ 
   geom_tiplab(size=0.5)
Dendro <- ggtree(hc1,layout ="rec")+ 
  geom_tiplab(size=1.5)


Eco_4panels <- ggalign::align_plots(PEco_A,  Dendro, PEco_B,unrooted ,
                                    ncol = 2)
png("plots/EcolinFriends/Eco_4panels.png",units = "cm",width =15 ,height = 20,res = 300)
print(Eco_4panels)
dev.off()

