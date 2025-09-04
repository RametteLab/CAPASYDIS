test_that("CalcDelta2seqs", {
  expect_equal(CalcDelta2seqs(REF=c("A","A"),SEQ=c("A","A"),DeltaVector = capasydis::DeltaMatrix), c(0,0))
  expect_equal(CalcDelta2seqs(REF=c("A","A"),SEQ=c("G","C"),DeltaVector = capasydis::DeltaMatrix), c(0.01,0.06))
})
