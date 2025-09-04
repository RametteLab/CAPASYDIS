#' Calculate the asymmetric distance for a dataframe of aligned sequences based on two reference sequences
#'
#' @param Alignment the alignment provided as a data frame of sequences (rows) by aligned positions (columns).
#' @param DeltaMatrix the substitution matrix used to define the distance between pairs of bases being compared.
#' @param PosRef1 the index (row number) of the first sequence to be used as the reference.
#' @param PosRef2 the index (row number) of the second sequence to be used as the reference.
#' @param Ncores the number of cores chosen for parallelization
#' @return A list of class "asymdist2D" containing the names of sequences, the distance of each sequence to sequence given in position 1 (disttoRefPos1), the distance of each sequence to sequence given in position 2 (disttoRefPos2).
#' @export
#'
#' @examples
#' Alignment_df <-
#'  data.frame(rbind(
#'    Reference=c("A","C","G","T","A","T","G","C","A","T"),
#'    S1=c("A","A","G","T","A","T","G","C","A","T"),
#'    S2=c("A","C","G","T","C","T","G","C","A","T"),
#'    S3=c("A","C","G","T","A","T","G","C","T","T"))
#'  )
#'
#'asymDist2D(Alignment=Alignment_df,DeltaMatrix=capasydis::DeltaMatrix,PosRef1=1,PosRef2=4,Ncores=2)

asymDist2D <- function(Alignment,
                            DeltaMatrix=capasydis::DeltaMatrix,
                            PosRef1=1,
                            PosRef2=4,
                            Ncores=2)
{
  T1 <- asymDist1D(Alignment=Alignment,DeltaMatrix=DeltaMatrix,PosRef=PosRef1,Ncores=Ncores)
  T2 <- asymDist1D(Alignment=Alignment,DeltaMatrix=DeltaMatrix,PosRef=PosRef2,Ncores=Ncores)
  T_df <- list(Names=rownames(Alignment),disttoRefPos1=T1$asymdist,disttoRefPos2=T2$asymdist)
  class(T_df) <- "asymdist2D" # Class definition
  return(T_df)
}
