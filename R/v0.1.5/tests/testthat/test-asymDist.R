test_that("asymDist test", {
  Alignment_df <-
    data.frame(rbind(
      Reference=c("A","C","G","T","A","T","G","C","A","T"),
      S1=c("A","A","G","T","A","T","G","C","A","T"),
      S2=c("A","C","G","T","C","T","G","C","A","T"),
      S3=c("A","C","G","T","A","T","G","C","T","T"))
    )
  # colnames(Alignment_df) <- paste0("P",1:10)

    A <- asymDist1D(Alignment=Alignment_df,PosRef=1,Ncores=1)
    A2 <- asymDist1D(Alignment=Alignment_df,PosRef=1,Ncores=2)
  expect_equal(all(unlist(A)==unlist(A2)),TRUE)
  expect_equal(as.numeric(round(A$Score[1],8)),22.468278)
  expect_equal(as.numeric(round(A$minScore,8)),22.468278)
  expect_equal(as.numeric(round(A$maxScore,8)),22.7725394)
  expect_equal(as.numeric(round(A$RangeScore,8)),0.30426118)
  expect_equal(as.numeric(round(A$asymdist[1],8)),0)
  expect_equal(as.numeric(round(A$asymdist[2],8)),0.06920517)
  expect_equal(as.numeric(round(A$asymdist[3],8)),0.04396354)
  expect_equal(as.numeric(round(A$asymdist[4],8)),0.04372505)
})

