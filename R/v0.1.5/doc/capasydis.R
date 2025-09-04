## ----include = FALSE----------------------------------------------------------
knitr::opts_chunk$set(
  collapse = TRUE,
  comment = "#>"
)

## ----setup--------------------------------------------------------------------
library(capasydis)

## -----------------------------------------------------------------------------
Alignment_df <-
 data.frame(rbind(
   S0=c("A","C","G","T","A","T","G","C","A","T"),
   S1=c("A","A","G","T","A","T","G","C","A","T"),
   S2=c("A","C","G","T","C","T","G","C","A","T"),
   S3=c("A","C","G","T","A","T","G","C","T","T"))
 )
A <- asymDist1D(Alignment = Alignment_df,DeltaMatrix=DeltaMatrix,PosRef=1)
colnames(Alignment_df) <- paste0("P",1:10)
Alignment_df

## -----------------------------------------------------------------------------
A


## -----------------------------------------------------------------------------
A$asymdist[2] == (A$Score[2]-A$minScore)/A$RangeScore

## ----echo=FALSE---------------------------------------------------------------
plot(A$asymdist,ylab="asymdist1D",pch=16,xlim=c(1,5),col=c("blue",1,1,1),main="1D example")
text(A$asymdist,rownames(Alignment_df),pos = 4,col=c("blue",1,1,1))


## -----------------------------------------------------------------------------
A$asymdist[3:4]
diff(A$asymdist[3:4])

## -----------------------------------------------------------------------------
B <- asymDist2D(Alignment = Alignment_df,DeltaMatrix = DeltaMatrix,
                   PosRef1 = 1,PosRef2 = 4)

B

## -----------------------------------------------------------------------------
plot(x=B$disttoRefPos1,y=B$disttoRefPos2,xlab="asymDist to ref1 (S0)",ylab="asymDist to ref2 (S3)",pch=16,xlim=c(0,.12),ylim=c(0,.12),main="2D example")
text(x=B$disttoRefPos1,y=B$disttoRefPos2,labels = B$Names,pos = 4,col=c("blue",1,1,"red"))

## -----------------------------------------------------------------------------
Alignment_df2m <-
 data.frame(rbind(
   S0=c("A","C","G","T","A","T","G","C","A","T","A","C","G","T","A","T","G","C","A","T"),
   S1=c("A","A","G","T","A","T","G","C","A","T","A","A","G","T","A","T","G","C","A","T"),
   S2=c("A","C","G","T","C","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.1=c("T","C","G","T","C","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.2=c("A","G","G","T","C","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.3=c("A","C","C","T","C","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.4=c("A","C","G","A","C","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.5=c("A","C","G","T","T","T","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.6=c("A","C","G","T","C","G","G","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.7=c("A","C","G","T","C","T","C","C","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.8=c("A","C","G","T","C","T","G","A","A","T","A","C","G","T","C","T","G","C","A","T"),
   S2.9=c("A","C","G","T","C","T","G","C","T","T","A","C","G","T","C","T","G","C","A","T"),
   S2.10=c("A","C","G","T","C","T","G","C","A","G","A","C","G","T","C","T","G","C","A","T"),
   S3=c("A","C","G","T","A","T","G","C","T","T","A","C","G","T","A","T","G","C","T","T"))
   
 )
Alignment_df2m <- rbind(Alignment_df2m,S0_10m =mutate(as.character(Alignment_df2m["S0",]),Nmut=10))
nrow(Alignment_df2m)

T_df2m <- asymDist2D(Alignment = Alignment_df2m,DeltaMatrix = DeltaMatrix,
                   PosRef1 = 1,PosRef2 = nrow(Alignment_df2m)) # the last one is the last row , 

plot(x=T_df2m$disttoRefPos1,y=T_df2m$disttoRefPos2,xlab="asymDist to ref1",ylab="asymDist to ref2 ",pch=16,main="2D mutants",xlim=c(0,0.3),col=c("blue","black","black",rep("lightpink",10),"black","red"))
text(x=T_df2m$disttoRefPos1,y=T_df2m$disttoRefPos2,labels = T_df2m$Names,pos = 4,col=c("blue","black","black",rep("lightpink",10),"black","red"))

## -----------------------------------------------------------------------------
Alignment_df2m_T <- data.frame(t(Alignment_df2m))
Alignment_ape <- ape::as.alignment(Alignment_df2m_T)

# genetic distance matrix 
DIST <- ape::dist.dna(x=ape::as.DNAbin(Alignment_ape), model = "raw", variance = FALSE,
              gamma = FALSE, pairwise.deletion = FALSE,
              base.freq = NULL, as.matrix = FALSE)

plot(hclust(DIST),hang = -1)

## -----------------------------------------------------------------------------
DIST_df <- as.matrix(DIST)
plot(x=DIST_df[,1],y=T_df2m$disttoRefPos1,pch=16,ylab="asymDist to Ref1",xlab="Phylogenetic distance to Ref1",col=c("blue","black","black",rep("lightpink",10),"black","red"),main="asymdist vs. phylogenetic distances")

## -----------------------------------------------------------------------------
table(DIST_df[,1])

## -----------------------------------------------------------------------------
table(T_df2m$disttoRefPos1)
all(unique(T_df2m$disttoRefPos1) == T_df2m$disttoRefPos1)

