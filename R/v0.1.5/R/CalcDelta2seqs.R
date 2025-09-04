#' Calculate the delta vector between two vectors of characters
#'
#' @param REF the reference vector of characters (e.g. bases)
#' @param SEQ the query vector of characters (e.g. bases)
#' @param DeltaVector the name vector of delta value: e.g. DeltaMatrix
#'
#' @return for each pair in REF and SEQ, the delta value is calculated.
#' @export
#'
#' @examples
#' CalcDelta2seqs(REF=c("A","A"),SEQ=c("A","A"),DeltaVector = capasydis::DeltaMatrix) # 0 0
#' CalcDelta2seqs(REF=c("A","A"),SEQ=c("G","C"),DeltaVector = capasydis::DeltaMatrix) # 0.01 0.06
#'
#'
CalcDelta2seqs <- function(REF,SEQ,DeltaVector)
{
  V <- vapply(paste0(REF,SEQ),FUN = function(x){  #x="TC"
    if(x %in% c("AA","GG","CC","TT","--","..")){return (0)}
    else{return(DeltaVector[x])} #
  },double(1)
  )
  return(as.numeric(V))
}

