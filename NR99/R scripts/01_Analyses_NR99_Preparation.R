# Title: NR99 v0.1.8 (higher precision) data preparation   -----
# .........................................................
# data was prepared with build_axes version 0.1.8 (e-10)
# followed by cleaning up with SILVA_go

# R libraries ----
library(tidyverse)

theme0 <- function(...) theme( legend.position = "none",
                               panel.background = element_blank(),
                               panel.grid.major = element_blank(),
                               panel.grid.minor = element_blank(),
                               panel.spacing = unit(0,"null"),
                               axis.ticks = element_blank(),
                               axis.text.x = element_blank(),
                               axis.text.y = element_blank(),
                               axis.title.x = element_blank(),
                               axis.title.y = element_blank(),
                               axis.ticks.length = unit(0,"null"),
                               #axis.ticks.margin = unit(0,"null"),
                               panel.border=element_rect(color=NA),...)

# 0) Data preparation  -----
# data colored by domains with 3 axes
D <- read.csv("data/output_build_axesv0.1.8_R1_R2_R3/output_R1_R2_R3_cleaned.csv",
              h=TRUE) # 10.5281/zenodo.17055026

# head(D %>% select(x,y,z))
# 0.03217729 # 
options(digits = 10) # to display the same as produced in the CSV
# The issue is typically with how R displays the numbers, not how it stores them.
# options(digits = n) controls the minimum number of significant digits to print.
# options(scipen = -n) discourages R from using scientific notation for large numbers. A negative value encourages it for very small numbers.

## Uniqueness of values ---- 
# here default is e-10
nrow(D)             #331663
length(unique(D$x)) #331132
length(unique(D$y)) #331113
length(unique(D$z)) #330993


D %>% select(x,y,z) %>% unique() %>% nrow() #331663
D %>% select(x,y) %>% unique() %>% nrow() #331662

D[duplicated(D[,c("x", "y")]), ]
# D[which(D$x==D[267285,]$x),]
# x           y            z
# 120682 5.44151e-05 0.099356024 0.0534556226
# 267285 5.44151e-05 0.099356024 0.0534556225
# label
# 120682 CP027597.1976524.1978045__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia;count_1
# 267285    CP026846.3775350.3776871__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Shigella;count_1


# so the 3D data is really unique
D %>% select(label) %>% unique()%>% nrow()


# truncation to e-8
# D %>% select(x,y,z) %>% round(.,8) %>% unique() %>% nrow() #331632
# D %>% select(x,y) %>% round(.,8) %>% unique() %>% nrow() #331609
# Conclusion: so not that bad even at e-8


# Green == Eukaryota
# Blue == Bacteria
# Red == Archaea

D <- 
  D %>% 
  mutate(
    domain = case_when(
      str_detect(label,"__Eukaryota;") ~ "Eukaryota", 
      str_detect(label,"__Bacteria;") ~ "Bacteria", 
      str_detect(label,"__Archaea;") ~ "Archaea", 
      TRUE ~ NA
    )
)

D <- D %>% mutate(color=case_when(
  domain == "Archaea" ~ "red",
  domain == "Eukaryota" ~ "green",
  domain ==  "Bacteria" ~"blue"

))

## . taxonomy ----
D$Phylum <- sapply(D$label,function(x)strsplit(x,";")[[1]][2])
D$Class <- sapply(D$label,function(x)strsplit(x,";")[[1]][3])
D$Order <- sapply(D$label,function(x)strsplit(x,";")[[1]][7]) 
#D$Family <- sapply(D$label,function(x)strsplit(x,";")[[1]][8]) 
# table(D$Phylum)%>% sort (., decreasing = TRUE) 

# Data without NA for specific taxonomic levels
DPhylum <- D %>% filter(Phylum != "NA")  # 558 NAs removed
DClass  <- D %>% filter(Class != "NA")  # 5031 NAs removed
DOrder  <- D %>% filter(Order != "NA")  # 230198 NAs removed
# write.csv(DPhylum,file = "DPhylum.csv",quote = FALSE,row.names = FALSE)  
# write.csv(DClass,file = "DClass.csv",quote = FALSE,row.names = FALSE)  

## . 3D distance to origin -----
# - do that when all axes are combined => so 3D distance to origin
# d=sqrt(x2+y2+z2)
DPhylum$dist3D <- DPhylum %>% select(x,y,z) %>% 
  apply(.,1,function(x) {sqrt(x[1]^2+x[2]^2+x[3]^2)})
# length(D$dist3D)==nrow(D)  # really unique values

DClass$dist3D <- DClass %>% select(x,y,z) %>% 
  apply(.,1,function(x) {sqrt(x[1]^2+x[2]^2+x[3]^2)})
# length(D$dist3D)==nrow(D)  # really unique values

DOrder$dist3D <- DOrder %>% select(x,y,z) %>% 
  apply(.,1,function(x) {sqrt(x[1]^2+x[2]^2+x[3]^2)})
# length(D$dist3D)==nrow(D)  # really unique values


## . split by domain ----
APhylum <- DPhylum %>% filter(domain =="Archaea") %>% as_tibble()
EPhylum <- DPhylum %>% filter(domain =="Eukaryota")%>% as_tibble()
BPhylum <- DPhylum %>% filter(domain =="Bacteria")%>% as_tibble()


## . Convex hulls -----
# Calculate convex hull for each group
A_3Dcoord <- APhylum %>% select(x,y,z) %>% as.data.frame()
A_3Dhull_indices <- geometry::convhulln(A_3Dcoord, options = "Pp") # "Pp" option returns points on the hull
A_3Dhull_points <- as.data.frame(A_3Dcoord[A_3Dhull_indices, ])
A_3Dcoord$domain <- "Archaea"
A_3Dhull_points$domain <- "Archaea"
##
B_3Dcoord <- BPhylum %>% select(x,y,z) %>% as.data.frame()
B_3Dhull_indices <- geometry::convhulln(B_3Dcoord, options = "Pp") # "Pp" option returns points on the hull
B_3Dhull_points <- as.data.frame(B_3Dcoord[B_3Dhull_indices, ])
B_3Dcoord$domain <- "Bacteria"
B_3Dhull_points$domain <- "Bacteria"
##
E_3Dcoord <- EPhylum %>% select(x,y,z) %>% as.data.frame()
E_3Dhull_indices <- geometry::convhulln(E_3Dcoord, options = "Pp") # "Pp" option returns points on the hull
E_3Dhull_points <- as.data.frame(E_3Dcoord[E_3Dhull_indices, ])
E_3Dcoord$domain <- "Eukaryota"
E_3Dhull_points$domain <- "Eukaryota"

All_3Dcoord <- rbind(A_3Dcoord,B_3Dcoord,E_3Dcoord)
All_3Dhull_points <- rbind(A_3Dhull_points,B_3Dhull_points,E_3Dhull_points)

# All_3Dcoord %>% nrow        # 331105      
# All_3Dhull_points %>% nrow  #678   


#.....................................................................
# Variance per domain ----
#.....................................................................
A_variance3D <- vector("list",0)
A_variance3D[["axes"]] <- apply(A_3Dcoord[,1:3],2,var)
A_variance3D[["total_variance"]] <- sum(apply(A_3Dcoord[,1:3],2,var))

B_variance3D <- vector("list",0)
B_variance3D[["axes"]] <- apply(B_3Dcoord[,1:3],2,var)
B_variance3D[["total_variance"]] <- sum(apply(B_3Dcoord[,1:3],2,var))

E_variance3D <- vector("list",0)
E_variance3D[["axes"]] <- apply(E_3Dcoord[,1:3],2,var)
E_variance3D[["total_variance"]] <- sum(apply(E_3Dcoord[,1:3],2,var))

A_variance3D
B_variance3D
E_variance3D

#.....................................................................
# Animals ----
#...........
ANIMALs <- EPhylum %>% filter(grepl("Animalia;BCP",label))
# dim(ANIMALs) #12375     9
# BCP Clade (Bilateria, Coelomata, Protostomia):
# _Eukaryota;Amorphea;Obazoa;Opisthokonta;Holozoa;Choanozoa;Metazoa;Animalia;BCP
#   This clade within Animalia represents a significant evolutionary branch characterized by the presence of bilateral symmetry, a true body cavity (coelom), and the development of the mouth before the anus during embryonic development

# grep(pattern = "Animalia",E$label,value = TRUE) %>% length #12375
# ANIMAL <- grep(pattern = "Animalia",E$label,value = TRUE,perl = TRUE) %>% strsplit(.,"Animalia") 
# lapply(ANIMAL,function(x) x[[2]]) %>% unlist %>% table


#.....................................................................
# Centroids ----
#.....................................................................
# as the median in each dimension

GetCentroidInfo=function(data=A){
  require(FNN)
  R <- vector("list",3);  names(R) <- c("centroid","closest_point","min_dist")
  Centroid <- matrix(apply(data %>% select(x,y,z) ,2,median), nrow = 1,
                     dimnames = list(NULL, c("x","y","z")))
  R[[1]] <- Centroid
  # closest point to the centroid
  # Find the k=1 nearest neighbor
  # knn.dist returns distances to the k nearest neighbors
  # knn.index returns the indices of the k nearest neighbors
  closest_nn <- get.knnx(data = as.matrix(data%>% select(x,y,z)), query = Centroid, k = 1)
  closest_point_index <- closest_nn$nn.index[1, 1] # Index of the closest point
  R[[2]] <- data[closest_point_index, ]
  R[[3]] <- closest_nn$nn.dist[1, 1]     # Distance to the closest point
  
  return(R)
}

A_centroid_info <- GetCentroidInfo(data=APhylum)
B_centroid_info <- GetCentroidInfo(data=BPhylum)
E_centroid_info <- GetCentroidInfo(data=EPhylum)

options(pillar.sigfig = 8)
A_centroid_info
B_centroid_info
E_centroid_info

#.....................................................................
## . Euclidean distance to the centroid ----
#.....................................................................
A_dist_to_centroid <- apply(APhylum %>% select(x,y,z), 1, function(x) {
  sqrt(sum((x - A_centroid_info$centroid)^2))
})

B_dist_to_centroid <- apply(BPhylum  %>% select(x,y,z), 1, function(x) {
  sqrt(sum((x - B_centroid_info$centroid)^2))
})

E_dist_to_centroid <- apply(EPhylum %>% select(x,y,z), 1, function(x) {
  sqrt(sum((x - E_centroid_info$centroid)^2))
})
#.....................................................................
## . Average Distance to Centroid ----
#.....................................................................
AvDistCentroid <- as.data.frame(matrix(c("Eukaryota",c(mean(E_dist_to_centroid),sd(E_dist_to_centroid)),
                                         "Archaea",c(mean(A_dist_to_centroid),sd(A_dist_to_centroid)),
                                         "Bacteria",c(mean(B_dist_to_centroid),sd(B_dist_to_centroid))),3,3,byrow = TRUE))
colnames(AvDistCentroid) <- c("domain","avg_distance_to_centroid", "sd")
AvDistCentroid[,2] <- round(as.numeric(AvDistCentroid[,2]),5)
AvDistCentroid[,3] <- round(as.numeric(AvDistCentroid[,3]),6)
AvDistCentroid


# between distances
Centroids <- rbind(
  A_centroid_info$centroid,
  B_centroid_info$centroid,
  E_centroid_info$centroid
)
BetweenDist_Centroids <- dist(Centroids) # Eucl distances

AllCentroids_dist <- data.frame(rbind(AvDistCentroid,c("Between domains",round(mean(BetweenDist_Centroids),5),
                                                       sd=round(sd(BetweenDist_Centroids),6))))

AllCentroids_dist$domain <- factor(AllCentroids_dist$domain,levels=c("Archaea","Bacteria","Eukaryota","Between domains"))
AllCentroids_dist$avg_distance_to_centroid <- as.numeric(AllCentroids_dist$avg_distance_to_centroid)
AllCentroids_dist$sd <- as.numeric(AllCentroids_dist$sd)
AllCentroids_dist$Type <- factor(c("within-domain","within-domain","within-domain","between-domain"),levels=c("within-domain","between-domain"))

AllCentroids_dist

## . all Distances to Centroid ----
## 
Raw_centroid_dists <- c(A_dist_to_centroid,B_dist_to_centroid,E_dist_to_centroid,BetweenDist_Centroids)
Raw_centroid_dists_tble <- data.frame(Type=c(rep("Archaea",length(A_dist_to_centroid)),
                                             rep("Bacteria",length(B_dist_to_centroid)),
                                             rep("Eukaryota",length(E_dist_to_centroid)),
                                             rep("Between domains",length(BetweenDist_Centroids))
),distance=as.numeric(Raw_centroid_dists))
Raw_centroid_dists_tble$Type <- factor(Raw_centroid_dists_tble$Type,levels=c("Archaea","Bacteria","Eukaryota","Between domains"))

## defining outliers vs. centroid ----
### . Archaea
q99_A <- Raw_centroid_dists_tble %>% filter(Type=="Archaea") %>% select(distance) %>% 
  summarise(q99=quantile(distance,probs = 0.99)) %>% pull

A_Table <- data.frame(APhylum %>% select(x,y,z,label,Phylum,Class),dist_to_centroid=A_dist_to_centroid)
A_Table$Aboveq99 <- A_Table$dist_to_centroid>= q99_A

### . Bacteria
q99_B <- Raw_centroid_dists_tble %>% filter(Type=="Bacteria") %>% select(distance) %>% 
    summarise(q99=quantile(distance,probs = 0.99)) %>% pull
  
  B_Table <- data.frame(BPhylum %>% select(x,y,z,label,Phylum,Class),
                                 dist_to_centroid=B_dist_to_centroid)
  B_Table$Aboveq99 <- B_Table$dist_to_centroid>= q99_B
  
### . Eukaryota
  q99_E <- Raw_centroid_dists_tble %>% filter(Type=="Eukaryota") %>% select(distance) %>% 
    summarise(q99=quantile(distance,probs = 0.99)) %>% pull

  E_Table <- data.frame(EPhylum %>% select(x,y,z,label,Phylum,Class),
                                 dist_to_centroid=E_dist_to_centroid)
   
  E_Table$Aboveq99 <- E_Table$dist_to_centroid>= q99_E
  

  
# outliers (q99) ----
#comparison of the q99 values
  format(c(q99_A, q99_B,q99_E), scientific = TRUE, digits = 3)
  
  A_Table %>% select(Aboveq99) %>% table #61
  A_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table
  A_Table_outliers <- round((A_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table)/(A_Table %>% filter(Aboveq99) %>% nrow)*100,1)
  # Aenigmarchaeota, deep =>                  “Deep Sea Euryarchaeotic Group (DSEG)”
  A_Table_outliers[apply(A_Table_outliers,1,sum) %>% sort(decreasing = TRUE) %>% names(),] #to sort the table by importance
  
  B_Table %>% filter(Aboveq99) %>% nrow  # 2912
  B_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table
  B_Table_raw <- B_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table
  
  B_Table_outliers <- round((B_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table)/(B_Table %>% filter(Aboveq99) %>% nrow)*100,1)
  B_outliers_names <- (apply(B_Table_outliers,1,sum) %>% sort(decreasing = TRUE) %>% names())
  
  apply(B_Table_outliers,1,sum) %>% sort(decreasing = TRUE)          
  B_Table_outliers[B_outliers_names[1:3],] #to sort the table by importance; top 3
  B_Table_raw[B_outliers_names[1:3],]
  
  
  E_Table %>% filter(Aboveq99) %>% nrow  # 344
  E_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table  
  E_Table_outliers <- round((E_Table %>% filter(Aboveq99) %>% select(Phylum,Class) %>% table)/(E_Table %>% filter(Aboveq99) %>% nrow)*100,1)
  E_Table_outliers[apply(E_Table_outliers,1,sum) %>% sort(decreasing = TRUE) %>% names(),] #to sort the table by importance
  
# Ecoli and Friends analyses ----

  B_Table$Genus <- sapply(B_Table$label,function(x)strsplit(x,";")[[1]][7]) 
  B_Table$Family <- sapply(B_Table$label,function(x)strsplit(x,";")[[1]][5]) 
  
  # grep("Escherichia",x = B_Table$label,value = TRUE)
  # B_Table %>% filter(grepl("CP002212",label))
  #CP002212.4561157.4562698__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia
  
  # Data without NA for specific taxonomic levels
  B_Table %>% filter(Genus == "NA") %>% count  # 160520 NAs 
  B_Genus  <- B_Table %>% filter(Genus != "NA")  # 
  
  # B_Table %>% filter(Family == "NA") %>% count()
  # B_Family  <- B_Table %>% filter(Family != "NA")  # 17142 NAs removed
  
  # table(B_Genus$Genus)%>% sort (., decreasing = TRUE) 
  
  B_Genus %>% filter(Genus=="Escherichia") %>% count() #948
  
  # DFamily <- DFamily # remove NAs in Genus
  
  
  # Enterobacteriaceae genera:
  ValidGenusNames=c( 
    "Apirhabdus",
    "Buttiauxella",
    "Cedecea",
    "Citrobacter",
    "Cronobacter",
    "Dryocola",
    "Enterobacillus",
    "Enterobacter",
    "Escherichia",
    "Franconibacter",
    "Gibbsiella",
    "Intestinirhabdus",
    "Izhakiella",
    "Klebsiella",
    "Kluyvera",
    "Kosakonia",
    "Leclercia",
    "Lelliottia",
    "Limnobaculum",
    "Mangrovibacter",
    "Metakosakonia",
    "Phytobacter",
    "Pluralibacter",
    "Pseudescherichia",
    "Pseudocitrobacter",
    "Rosenbergiella",
    "Saccharobacter",
    "Salmonella",
    "Scandinavium",
    "Shigella",
    "Shimwellia",
    "Siccibacter",
    "Silvania",
    "Tenebrionibacter",
    "Tenebrionicola",
    "Trabulsiella",
    "Yokenella",
    "Annandia",#candidatus
    "Arocatia",
    "Aschnera",
    "Benitsuchiphilus",
    "Blochmannia",
    "Curculioniphilus",
    "Cuticobacterium",
    "Doolittlea",
    "Gillettellia",
    "Gullanella",
    "Hamiltonella",
    "Hartigia",
    "Hoaglandella",
    "Ischnodemia",
    "Ishikawaella",
    "Kleidoceria",
    "Kotejella",
    "Macropleicola",
    "Mikella",
    "Moranella",
    "Phlomobacter",
    "Profftia",
    "Purcelliella",
    "Regiella",
    "Riesia",
    "Rohrkolberia",
    "Rosenkranzia",
    "Schneideria",
    "Stammera",
    "Stammerula",
    "Tachikawaea",
    "Westeberhardia",
    "Aquamonas",
    "Atlantibacter",
    "Superficieibacter")
  
  
  
ValidGenusNames %>% length() #72
  
DEnterobacteriaceae <- B_Genus %>% #keeping only valid genus names
    filter(Genus %in% ValidGenusNames)
  
# DEnterobacteriaceae %>% select(Genus) %>% unique() %>% arrange(Genus)
  
#selecting top Genera with more than 10 sequences (they represent the top 15 names)
DEnterobacteriaceaeTop15 <- DEnterobacteriaceae %>%
    group_by(Genus) %>%
    filter(n() >= 10) %>%
    ungroup() 
  
DEnterobacteriaceaeTop15 %>% select(Genus) %>% unique()
  
## . Convex hulls -----
# Calculate convex hull for each group
Entero_3Dcoord <- DEnterobacteriaceaeTop15 %>% select(x,y,z,Genus) %>% as.data.frame()
Entero_3Dhull_indices <- geometry::convhulln(Entero_3Dcoord%>% select(x,y,z), options = "Pp") # "Pp" option returns points on the hull
Entero_3Dhull_points <- as.data.frame(Entero_3Dcoord[Entero_3Dhull_indices, ])
  
# Entero_3Dcoord %>% nrow        #4343      
# Entero_3Dhull_points %>% nrow  #300

# only Escherichia ----
DEnterobacteriaceae %>% nrow() #4408
B_Eco <- DEnterobacteriaceae %>% 
  mutate(isEcoli=case_when(
    Genus=="Escherichia" ~ TRUE,
    .default = FALSE
  ) ) #948
# colnames(B_Eco)

Eco_coord <- B_Eco %>% filter(isEcoli) %>% select(x,y,Aboveq99)
Eco_CHull <- chull(x = Eco_coord$x,y=Eco_coord$y)
Eco_polygon <- Eco_coord[Eco_CHull,c("x","y")]


Eco_only <- B_Eco %>%  filter(isEcoli)
Eco_only %>% count #948

Ecoli3D <- Eco_only  %>%  select(x,y,z)
Ecoli2D <- Eco_only  %>%  select(x,y)
Ecoli2DLabel <- Eco_only  %>%  select(x,y,label)
Ecoli1D <- Eco_only  %>%  select(x,Aboveq99)
Ecoli1D$seq_names_short <- as.character(sapply(Eco_only$label,function(x) strsplit(x,"__")[[1]][1]))
Ecoli2D$seq_names_short <- as.character(sapply(Eco_only$label,function(x) strsplit(x,"__")[[1]][1]))
colnames(Ecoli1D)[1] <- "asymDist1D"
colnames(Ecoli2D)[1] <- "asymDist2D"

# save.image("data/01_NR99.RData")
load("data/01_NR99.RData")
