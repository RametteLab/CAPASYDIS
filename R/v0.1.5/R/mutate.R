#' Mutate the bases in a DNA sequence (only SNVs)
#'
#' @param S Sequence of bases, as a vector of bases.
#' @param Nmut number of mutations (SNVs) in the sequence
#'
#' @return A vector of length(S) containing Nmut mutated bases
#' @export
#'
#' @examples
#' Seqi <-  paste(rep(c("A","T","G","C"),25))
#' Seqm <- mutate(Seqi,Nmut=10)
#' table(Seqi==Seqm)

mutate=function(S,Nmut=3){
  if(!is.character(S) || length(S) <= 1){stop("a vector of bases should be provided. Exiting")}
  n <- length(S)
  POSr <- sample(n,Nmut,replace = FALSE) #define random positions
  Smut <- S
  Pool=c("A","T","G","C")
  for(i in 1:Nmut){
    Base <- S[POSr[i]]
    Smut[POSr[i]] <- sample(setdiff(Pool,Base),1)
  }
  return(Smut)
}




