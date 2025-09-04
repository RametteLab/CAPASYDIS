#' Calculate the asymmetric distance for a dataframe of aligned sequences based on one reference sequence
#'
#' @param Alignment the alignment provided as a data frame of sequences (rows) by aligned positions (columns).
#' @param DeltaMatrix the substitution matrix used to define the distance between pairs of bases being compared.
#' @param PosRef the index (row number) of the sequence to be used as the reference.
#' @param Ncores the number of cores chosen for parallelization
#' @return A list of class "list", and "asymdist1D" containing the **minScore**, **maxScore**, range of the scores **RangeScore**, and the asymmetric distances (**asymdist**) for each sequence using the PosRef sequence as reference sequence.
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
#' asymDist1D(Alignment=Alignment_df,DeltaMatrix=DeltaMatrix,PosRef=1)

asymDist1D <- function(Alignment,
                            DeltaMatrix=capasydis::DeltaMatrix,
                            PosRef=1,
                            Ncores=1)
{
  REFSeq <- as.character(Alignment[PosRef,])
  Nseqs=nrow(Alignment)
  Npos=ncol(Alignment)
  Results <- data.frame(Names=rownames(Alignment),Score=rep(NA,Nseqs),asymdist=rep(NA,Nseqs))
  #------------
  # if(!is.integer(Ncores)){stop("Ncores should be an integer. Calculations aborted.")}
   if(Ncores<1){stop("Ncores should be positive. Calculations aborted.")}
  if(Ncores==1){
    REFSeq <- as.character(Alignment[PosRef,])
    DeltaVecDf_x <- apply(Alignment,1,FUN=function(x) capasydis::CalcDelta2seqs(REF=REFSeq,x,DeltaVector=DeltaMatrix))
    DeltaVecDf <- matrix(DeltaVecDf_x,ncol = Nseqs,nrow=Npos,byrow = FALSE) # need to swap rows and columns !
  }
  #------------
  if(Ncores>1){
    cl <- parallel::makeCluster(Ncores)
    junk <- parallel::clusterEvalQ(cl,
                                   {library(capasydis)}
    )
    REFSeq <- as.character(Alignment[PosRef,])
    parallel::clusterExport(cl, list("REFSeq"), envir = environment())
    DeltaVecDf_x <- parallel::parRapply(cl,Alignment,FUN=function(x) capasydis::CalcDelta2seqs(REF=REFSeq,x,DeltaVector=DeltaMatrix))
    parallel::stopCluster(cl)
    DeltaVecDf <- matrix(DeltaVecDf_x,ncol = Nseqs,nrow=Npos,byrow = FALSE) # need to swap rows and columns !
  }
  #--------------------------
  Results <- vector("list")
  Results[["Score"]] <- apply(DeltaVecDf,2,capasydis::sumSqrtScore)

  #standardization
    # minScore the minimum score to be used to standardize the distance (optional). If missing, this is inferred based on the chosen reference sequence.
    # maxScore the maximum score to be used to standardize the distance (optional). If missing, this is inferred using both the chosen reference sequence and the DeltaMatrix values.
    minScore <- Results[["Score"]][PosRef]
    MatrixMaxDivergence <- c(
      A=max(DeltaMatrix[grep("A",names(DeltaMatrix))]),
      G=max(DeltaMatrix[grep("G",names(DeltaMatrix))]),
      C=max(DeltaMatrix[grep("C",names(DeltaMatrix))]),
      T=max(DeltaMatrix[grep("T",names(DeltaMatrix))]),
      "-"=max(DeltaMatrix[grep("-",names(DeltaMatrix))]))
    MostDivergentSequenceDelta <- vapply(REFSeq,FUN=function(x)MatrixMaxDivergence[x],double(1))
    maxScore <- capasydis::sumSqrtScore(MostDivergentSequenceDelta)

  RangeScore <- maxScore-minScore
  Results[["minScore"]] <- minScore[1]
  Results[["maxScore"]] <- maxScore[1]
  Results[["RangeScore"]] <- RangeScore[1]
  Results[["asymdist"]] <- vapply(Results$Score,FUN=function(x){(x-minScore)/RangeScore},double(1))
  class(Results) <- c("list","asymdist1D") # Class definition
  return(Results)
}
