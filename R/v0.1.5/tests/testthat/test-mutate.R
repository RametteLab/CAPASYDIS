test_that("mutate test", {
  Seqi <-  paste(rep(c("A","T","G","C"),25))
  Seqm <- mutate(Seqi,Nmut=10)
  expect_equal(sum(Seqi!=Seqm),10)
  expect_equal(sum(Seqi==Seqm),90)
})

