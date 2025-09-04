#' Sum of square roots of delta values
#'
#' @param DeltaVector
#' The vector of delta values as compared to the aligned reference sequence.
#' @return
#' The score (numeric value) of the vector, as the sum of square root of each delta value. This is the unstandardized calculation. The score is not bound
#' @export
#'
#' @examples
#'
#'sumSqrtScore(rep(0,10000)) # 666716.5  score of a reference (deltas are all 0) that is 10000-nt long
#' sumSqrtScore(rep(0,100)) #  671.4629  score of a reference that is 100-nt long
#' sumSqrtScore(rep(0,10)) #  22.46828  score of a reference that is 10-nt long


sumSqrtScore <- function(DeltaVector=c(0.1,rep(0,10)))# base function to calculate the distance) # delta values to the reference sequence
{
  Nsites <- length(DeltaVector)
  Seq_F.vec <- 1:Nsites + DeltaVector #forward-distance vector
  return(sum(sqrt(Seq_F.vec)))
}

