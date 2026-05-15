DATE="20260513"
load(paste0(DATE,"_data.RData"))


# used libraries
require(dplyr)
require(ggplot2)

# DATA preparation--------------------------------------------------------------
D1temp <- D1
  colnames(D1temp) <- c("F_x","F_y","F_z","F_label","F_color","F_Nmut" )
D1revtemp <- D1rev
  colnames(D1revtemp) <- c("R_x","R_y","R_z","R_label","R_color","R_Nmut" )


DFR <- cbind(D1temp,D1revtemp)
DFR <- DFR %>% select(F_x,R_x,
                      F_y,R_y,
                      F_z,R_z,
                      F_label,R_label,
                      F_color,R_color,
                      F_Nmut,R_Nmut
)
head(DFR)
Cols_unique_pairs <- unique(D1[c("color", "Nmut")])


# PLOTS ------------------------------------------------------------------------

Plot=function(Object,W=20,H=10,R=300,FileName="test.png"){
  FilePath=paste0("plots/",FileName)
  png(FilePath,units = "cm",width =W ,height = H,res = R)
  print(Object)
  dev.off()
}

Px <- DFR %>% ggplot(aes(x=F_x,y=R_x,color=F_color))+
  geom_point()+
  theme_light()+
  scale_fill_identity(
    name = "Color",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) +
  scale_color_identity(
    name = "Number of mutations",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) 

Plot(Px,FileName="Forward_Reverse_X.png")


Py <- DFR %>% ggplot(aes(x=F_y,y=R_y,color=F_color))+
  geom_point()+
  theme_light()+
  scale_fill_identity(
    name = "Color",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) +
  scale_color_identity(
    name = "Number of mutations",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) 

Plot(Py,FileName="Forward_Reverse_Y.png")


Pz <- DFR %>% ggplot(aes(x=F_z,y=R_z,color=F_color))+
  geom_point()+
  theme_light()+
  scale_fill_identity(
    name = "Color",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) +
  scale_color_identity(
    name = "Number of mutations",
    breaks = Cols_unique_pairs$color,
    labels = Cols_unique_pairs$Nmut,
    guide = "legend"
  ) 

Plot(Pz,FileName="Forward_Reverse_Z.png")

# RV coefficient -------------------------------------------
require(FactoMineR)# install.packages("FactoMineR")
i=1
j=1000 #takes few mins with 10,000 already rv=0.78
X <- DFR[i:j,c("F_x","F_y")]
Y <- DFR[i:j,c("R_x","R_y")]
coeffRV(X,Y)$rv
coeffRV(X,Y)$p.value

# edit(coeffRV)

coeffRV_Modif <- function (X, Y) 
{ # modify to avoid significance testing
    tr <- function(Z) {
      sum(diag(Z))
    }
      Y <- scale(Y, scale = FALSE)
      X <- scale(X, scale = FALSE)
      W1 <- t(X) %*% X
      W2 <- t(Y) %*% Y
      W3 <- t(X) %*% Y
      W4 <- t(Y) %*% X
      rv <- tr(W3 %*% W4)/(tr(W1 %*% W1) * tr(W2 %*% W2))^0.5
    return(rv)
  
}

coeffRV(X,Y)$rv

coeffRV(X,Y)$rv 
coeffRV_Modif(X,Y) #ok validated with j=100, 1000, 10000

X <- DFR[,c("F_x","F_y")]
Y <- DFR[,c("R_x","R_y")]
coeffRV(X,Y)$rv  # Error: cannot allocate vector of size 8192.0 Gb
coeffRV_Modif(X,Y) #0.6023526
**conclusion**: with 0.60, the overall covariance patterns in 2D Forward and 2D reverse are 
somewhat similar (on a 0–1 scale).

#When i tried to run the Pvalue for the whole dataset, message "Error: cannot allocate vector of size 8192.0 Gb"
# issues for the calculation steps betax, betay

# Matrix similarities -------------------------------------------
N=1000
for (N in c(100,500, 1000,5000, 10000)){
 D2_F <- dist(DFR[1:N,] %>% select(F_x,F_y)) #Error: cannot allocate vector of size 4096.0 Gb
  D2_R <- dist(DFR[1:N,] %>% select(R_x,R_y))
  C <- cor(D2_F,D2_R) #0.8382223539
  cat("N=",N,"- Pearson's r=",round(C,3),"\n")
} 
  
for (N in c(15000,20000,30000,50000)){
  D2_F <- dist(DFR[1:N,] %>% select(F_x,F_y)) #Error: cannot allocate vector of size 4096.0 Gb
  D2_R <- dist(DFR[1:N,] %>% select(R_x,R_y))
  C <- cor(D2_F,D2_R) #0.8382223539
  cat("N=",N,"- Pearson's r=",round(C,3),"\n")
}
  # N= 100 -   Pearson's r= 0.872 
  # N= 500 -   Pearson's r= 0.83 
  # N= 1000 -  Pearson's r= 0.838 
  # N= 5000 -  Pearson's r= 0.783 
  # N= 10000 - Pearson's r= 0.753 
  # N= 15000 - Pearson's r= 0.76 
  # N= 20000 - Pearson's r= 0.74 
  # N= 30000 - Pearson's r= 0.703 
  # N= 50000 - Pearson's r= 0.711               

# test slices
nrow(DFR) #1,048,576
Delta=10000
S <- seq(0,100000,Delta)
ResD10KS100K=data.frame(i=rep(NA,length(S)),j=rep(NA,length(S)),r=rep(NA,length(S)))
Compa=1
for (i in S){
  j=i+Delta
  D2_F <- dist(DFR[i:j,c("F_x","F_y")] )
  D2_R <- dist(DFR[i:j,c("R_x","R_y")] )
  C <- cor(D2_F,D2_R) 
  ResD10KS100K[Compa,] <- c(i,j,round(C,4))
  Compa=Compa+1
} 
ResD10KS100K
# i      j      r
# 1  0e+00  10000 0.7534
# 2  1e+04  20000 0.7417
# 3  2e+04  30000 0.7610
# 4  3e+04  40000 0.7129
# 5  4e+04  50000 0.7719
# 6  5e+04  60000 0.7559
# 7  6e+04  70000 0.7798
# 8  7e+04  80000 0.7619
# 9  8e+04  90000 0.7230
# 10 9e+04 100000 0.7203
# 11 1e+05 110000 0.7557


Delta=20000
S <- seq(0,nrow(DFR)-Delta,Delta)
ResD20KAll=data.frame(i=rep(NA,length(S)),j=rep(NA,length(S)),r=rep(NA,length(S)))
Compa=1
for (i in S){
  j=i+Delta
  D2_F <- dist(DFR[i:j,c("F_x","F_y")] )
  D2_R <- dist(DFR[i:j,c("R_x","R_y")] )
  C <- cor(D2_F,D2_R) 
  ResD20KAll[Compa,] <- c(i,j,round(C,4))
  Compa=Compa+1
} 
ResD20KAll
saveRDS(ResD20KAll,"ResD20KAll.RDS")
# ResD20KAll <- readRDS("ResD20KAll.RDS")

# plot the lines based on window scan, and the point estimates 

SlicesFromBeginning <- data.frame(matrix(c(100,0.872,
  500,0.83,
  1000,0.838,
  5000,0.783,
  10000,0.753,
  15000,0.76,
  20000,0.74,
  30000,0.703,
  50000,0.711),9,2,byrow = TRUE))
colnames(SlicesFromBeginning) <- c("i","r")

C1 <- ResD20KAll %>% ggplot(aes(y=r,x=i+10000))+
                        geom_line()+
  theme_light()+ylim(0,1)+ylab("Pearson's r") +xlab("Number of sequences considered") +
  labs(title="Correlation between forward(x,y) and reverse (x,y) distance matrices",
       subtitle = "(window size of 20,000 points)")+
  geom_point(data=SlicesFromBeginning,aes(x=i,y=r),color="red")+
    geom_text(data=data.frame(x=20000,y=0.9),
              aes(x=x,y=y,
              label = "r estimation from the beginning to this point"),
              size=2.5,hjust="left",color="red")
    
Plot(C1,FileName="Matrix_Correlation_FandR_XY.png")












