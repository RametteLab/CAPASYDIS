test_that("MostDivergentScore", {
  REF1= c("A", "C","G")
  M <-  MostDivergentScore(SEQ=REF1,MatrixMaxDivergence=MatrixMaxDivergence,DeltaMatrix = DeltaMatrix)
  expect_equal(M$Length, 3)
  expect_equal(as.numeric(M$max.delta.vec[1]),0.11)
  expect_equal(round(M$Score,1),  4.3)
  expect_equal(M$MostDivSeq[1], "-")

})
