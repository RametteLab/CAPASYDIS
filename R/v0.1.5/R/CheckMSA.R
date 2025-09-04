#' Check for absence of unknown or of wobble bases in all the sequences of an alignment
#'
#' @param Alignment the sequences as a data frame of aligned sequences (rows) of single chars
#' @param type the type of bases to check for "N" for unknown only, "W" for wobble bases (not "A","T","G","C"), or "b" for both
#'
#' @return A list with two slots indicating the sequences (rows) where unknowns or wobble bases were found.
#' @export
#'
#' @examples
#' AL <- rbind(
#' S1=c("A","T","G","C","A","T","G","C"),
#' S2=c("A","N","G","C","A","T","G","C"),
#' S3=c("A","T","G","C","A","T","G","C"),
#' S4=c("A","T","G","W","A","T","G","C"),
#' S5=c("A","T","G","W","A","T","G","N")
#' )
#' C <- CheckMSA(Alignment = AL)
#' C$N # 2 5
#' C$W # NULL
#' C <- CheckMSA(Alignment = AL,type="W")
#' C$N # NULL
#' C$W # 2 4 5
#' C <- CheckMSA(Alignment = AL,type="b")
#' C$N # 2 5
#' C$W # 2 4 5
#' C <- CheckMSA(Alignment = AL[c(1,3),],type="b")
#' C$N # 0
#' C$W # 0
#'

CheckMSA = function(Alignment,type="N"){
  Result <- list()
  if(type=="N" || type=="b"){
    # cat("Checking for unknown bases:\n")
    N <- as.numeric(which(apply(Alignment,1,function(x) any(x=="N"))))
    if(length(N)>0){ #cat("  -> Some sequences with Ns were found\n");
      Result$N <- N} else{Result$N <- 0}
  }
  if(type=="W" || type=="b"){
    #cat("Checking for wobble bases (not ATGC):\n")
    W <- apply(Alignment,1,function(x) !all(names(table(x)) %in% c("A","T","G","C")))
    if(any(W)){
      # cat("  -> Some sequences with wobble bases were found\n");
      Result$W <- as.numeric(which(W))} else{Result$W <- 0}
  }
  return(Result)
}
