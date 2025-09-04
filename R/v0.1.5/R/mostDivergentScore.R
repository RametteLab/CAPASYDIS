#' Calculate the most divergent score and sequence for one sequence of characters
#'
#' @param SEQ the sequence of characters of interest (e.g. SEQ=REF1)
#' @param MatrixMaxDivergence the vector of maximum divergence
#' @param DeltaMatrix the delta matrix
#'
#' @return  A list is produced with  `max.delta.vec`, as the maximum diverging delta for each character in the provided sequence; `Score` the associated score of this most diverging sequence; and `MostDivSeq`, the most divergent sequence of characters, based on DeltaMatrix and MatrixMaxDivergence values.
#' @export
#'
#' @examples
#' SEQ1= c("A", "C","G","T","A","T","G","C","A","T")
#' MostDivergentScore(SEQ=SEQ1,MatrixMaxDivergence=MatrixMaxDivergence,DeltaMatrix = DeltaMatrix)
#'# $Length
#'# [1] 10
#'#
#'# $max.delta.vec
#'# A    C    G    T    A    T    G    C    A    T
#'# 0.11 0.14 0.13 0.12 0.11 0.12 0.13 0.14 0.11 0.12
#'#
#'# $Score
#'# [1] 22.77254
#'#
#'# $MostDivSeq
#'# [1] "-" "-" "-" "-" "-" "-" "-" "-" "-" "-"

MostDivergentScore=function(SEQ,MatrixMaxDivergence=MatrixMaxDivergence,DeltaMatrix=DeltaMatrix)
{
  Bases_1=vapply(as.character(names(DeltaMatrix)),FUN=function(x) strsplit(x,split = "")[[1]][1],character(1))
  DeltaMatrix_df <- data.frame(
    Bases_1=Bases_1,
    Bases_2=vapply(as.character(names(DeltaMatrix)),FUN=function(x) strsplit(x,split = "")[[1]][2],character(1)),
    Value=as.numeric(DeltaMatrix)
  )
  N=length(SEQ)
  SEQ=as.character(SEQ)
  # the delta values for each position
  delta.vec <- vapply(as.character(SEQ),FUN=function(x)MatrixMaxDivergence[x],double(1))
  Score <- sumSqrtScore(delta.vec)

  # The most dissimilar sequence
  MostDivSeq <- rep("N",N)
  for(k in 1:N){
    Dtemp <- DeltaMatrix_df[Bases_1==SEQ[k],]
    MostDivSeq[k] <- as.character(Dtemp[order(Dtemp$Value,decreasing = TRUE),][1,2])

  }
  return( list(
    Length=N,
    max.delta.vec = delta.vec,# maximum diverging sequence
    Score = Score, #raw asymdist
    MostDivSeq =MostDivSeq # most divergent sequence in bases
  )
  )
}
