test_that("asymDist2D_para test", {
  Alignment_df <- data.frame(rbind(
      Reference=c("A","C","G","T","A","T","G","C","A","T"),
      S1=c("A","A","G","T","A","T","G","C","A","T"),
      S2=c("A","C","G","T","C","T","G","C","A","T"),
      S3=c("A","C","G","T","A","T","G","C","T","T"))
    )
  A <- asymDist2D(Alignment=Alignment_df,
                   DeltaMatrix=DeltaMatrix,
                   PosRef1=1,
                   PosRef2=4,
                   Ncores=2
  )
  expect_equal(A$disttoRefPos1[1],0)
  expect_equal(A$disttoRefPos2[4],0)
  expect_equal(round(A$disttoRefPos1[2],5),0.06921)
  expect_equal(round(A$disttoRefPos2[1],8),0.04348834)
})




