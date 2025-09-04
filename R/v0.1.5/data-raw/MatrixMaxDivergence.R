MatrixMaxDivergence <- c(
  A=max(DeltaMatrix[grep("A",names(DeltaMatrix))]),
  G=max(DeltaMatrix[grep("G",names(DeltaMatrix))]),
  C=max(DeltaMatrix[grep("C",names(DeltaMatrix))]),
  T=max(DeltaMatrix[grep("T",names(DeltaMatrix))]),
  "-"=max(DeltaMatrix[grep("-",names(DeltaMatrix))]))

usethis::use_data(MatrixMaxDivergence, overwrite = TRUE)


