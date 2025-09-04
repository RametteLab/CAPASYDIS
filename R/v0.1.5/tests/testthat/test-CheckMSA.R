test_that("CheckMSA() test", {
  AL <- rbind(
    S1=c("A","T","G","C","A","T","G","C"),
    S2=c("A","N","G","C","A","T","G","C"),
    S3=c("A","T","G","C","A","T","G","C"),
    S4=c("A","T","G","W","A","T","G","C"),
    S5=c("A","T","G","W","A","T","G","N")
  )
  expect_equal(C <- CheckMSA(Alignment = AL)$N, c(2,5))
  expect_equal(C <- CheckMSA(Alignment = AL)$W, NULL)
  expect_equal(C <- CheckMSA(Alignment = AL,type="W")$N, NULL)
  expect_equal(C <- CheckMSA(Alignment = AL,type="W")$W, c(2,4,5))
  expect_equal(C <- CheckMSA(Alignment = AL,type="b")$N, c(2,5))
  expect_equal(C <- CheckMSA(Alignment = AL,type="b")$W, c(2,4,5))
  expect_equal(C <- CheckMSA(Alignment = AL[c(1,3),],type="b")$N, 0)
  expect_equal(C <- CheckMSA(Alignment = AL[c(1,3),],type="b")$W, 0)
})


