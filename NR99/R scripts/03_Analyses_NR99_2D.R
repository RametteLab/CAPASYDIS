# Title: 2D analysis of NR99 analysis using CAPASYDIS  -----
# run "Analyses_NR99_Preparation.R" before

# used libraries
require(dplyr)
require(ggplot2)
library(geometry) # install.packages("geometry")

dir.create("plots/2D_analyses")

# PlotTest=function(Object,W=20,H=10,R=75){
#   png("plots/test.png",units = "cm",width =W ,height = H,res = R)
#   print(Object)
#   dev.off()
# }

# coordinate calculations ----
#areas of the convex hull
A_coord <- A_Table %>% select(x,y,Phylum,Class,Aboveq99)
A_CHull <- chull(x = A_Table$x,y=A_Table$y)
A_polygon <- A_coord[A_CHull,c("x","y")]
#...................................
B_coord <- B_Table %>% select(x,y,Phylum,Class,Aboveq99)
B_CHull <- chull(x = B_Table$x,y=B_Table$y)
B_polygon <- B_coord[B_CHull,c("x","y")]
#...................................
E_coord <- E_Table %>% select(x,y,Phylum,Class,Aboveq99)
E_CHull <- chull(x = E_Table$x,y=E_Table$y)
E_polygon <- E_coord[E_CHull,c("x","y")]

# [Fig P2D_9panels] ----
P2Dall_A <- DPhylum %>%  select(x,y,domain) %>% ggplot(aes(x=x,y=y))+
   xlim(0,0.10)+ ylim(0,0.16) + theme_minimal()+xlab("")+ylab("")+
   geom_point(data=DPhylum %>% filter(domain !="Archaea"),aes(x=x,y=y),color="grey")+
   geom_point(data=DPhylum %>% filter(domain =="Archaea"),aes(x=x,y=y),color="red")+
   geom_polygon(data=A_polygon,aes(x=x,y=y),colour="red",fill=NA)+
    annotate("text",label="B",x=0.10,y=0.15,size=8,family = c("sans"),fontface="bold" )

P2Dall_B <- DPhylum %>%  select(x,y,domain) %>% ggplot(aes(x=x,y=y))+
    xlim(0,0.10)+ ylim(0,0.16) + theme_minimal()+xlab("")+ylab("")+
    geom_point(data=DPhylum %>% filter(domain!="Bacteria"),aes(x=x,y=y),color="grey")+
    geom_point(data=DPhylum %>% filter(domain=="Bacteria"),aes(x=x,y=y),color="blue")+
    geom_polygon(data=B_polygon,aes(x=x,y=y),colour="blue",fill=NA)+
  annotate("text",label="C",x=0.10,y=0.15,size=8,family = c("sans"),fontface="bold" )


  P2Dall_E <- DPhylum %>%  select(x,y,domain) %>% ggplot(aes(x=x,y=y))+
    xlim(0,0.10)+ ylim(0,0.16) + theme_minimal()+xlab("")+ylab("")+
    geom_point(data=DPhylum %>% filter(domain!="Eukaryota"),aes(x=x,y=y),color="grey")+
    geom_point(data=DPhylum %>% filter(domain=="Eukaryota"),aes(x=x,y=y),color="darkgreen")+
    geom_polygon(data=E_polygon,aes(x=x,y=y),colour="darkgreen",fill=NA)+
    annotate("text",label="D",x=0.10,y=0.15,size=8,family = c("sans"),fontface="bold" )
  
## 2D areas of the convex hull ----
P2D_A <-  A_coord %>%  ggplot(aes(x=x,y=y,fill = factor(Aboveq99)))+ 
    xlim(0.00,0.06)+ ylim(0.09,0.13) + theme_minimal()+xlab("")+ylab("y axis")+
    geom_point(aes(colour=Aboveq99))+ theme(legend.position = "none") +
    scale_color_manual(values =c("red","black"))+
    geom_point(data=data.frame(A_centroid_info$centroid),aes(x=x,y=y),color="white",shape = 10,size=3,inherit.aes = FALSE)+
    geom_polygon(data=A_polygon,aes(x=x,y=y),colour="red",fill=NA)+
    annotate("text",label="E",x=0.06,y=0.129,size=8,family = c("sans"),fontface="bold" )
  
  
P2D_B <-  B_coord %>% ggplot(aes(x=x,y=y,fill = factor(Aboveq99)))+ 
    xlim(0,0.06)+ ylim(0.08,0.14) + xlab("")+ylab("")+theme_minimal()+
  geom_point(aes(colour=Aboveq99))+ theme(legend.position = "none") +
  scale_color_manual(values =c("blue","black"))+
  geom_point(data=data.frame(B_centroid_info$centroid),aes(x=x,y=y),color="white",shape = 10,size=3,inherit.aes = FALSE)+
  geom_polygon(data=B_polygon,aes(x=x,y=y),colour="blue",fill=NA)+
  annotate("text",label="F",x=0.06,y=0.139,size=8,family = c("sans"),fontface="bold" )
  
P2D_E <-  E_coord %>% ggplot(aes(x=x,y=y,fill = factor(Aboveq99)))+ 
   xlab("")+ylab("")+theme_minimal()+
  geom_point(aes(colour=Aboveq99))+ theme(legend.position = "none") +
  scale_color_manual(values =c("darkgreen","black"))+
  geom_point(data=data.frame(E_centroid_info$centroid),aes(x=x,y=y),color="white",shape = 10,size=3,inherit.aes = FALSE)+
  geom_polygon(data=E_polygon,aes(x=x,y=y),colour="darkgreen",fill=NA)+
  annotate("text",label="G",x=0.1,y=0.125,size=8,family = c("sans"),fontface="bold" )

## 2D contour plots: domains ----
P2DC_A <-   ggplot(A_coord %>% filter(x<0.025,y<0.10), aes(x=x, y=y)) +
  stat_density_2d(geom = "polygon", contour = TRUE,
                  aes(fill = after_stat(level)), colour = "red",
                  bins = 50) +xlab("")+ylab("")+
  scale_fill_distiller(palette = "Reds", direction = 1) +
  theme_classic()+ xlim(0.015,0.021)+ ylim(0.092,0.099)+
  geom_point(data=A_centroid_info$centroid,aes(x=x,y=y),color="white",shape = 10,size=4)+
  theme(legend.position= c(0.8, 0.8),legend.key.size = unit(0.2, 'cm'))+
  guides(fill=guide_legend(title="Density"),color=guide_legend(override.aes=list(fill=NA)))+
  annotate("text",label="H",x=0.021,y=0.0987,size=8,family = c("sans"),fontface="bold" )

 # PlotTest(Object=P2DC_A,W=20,H=10,R=75)
 
## . Bacteria
P2DC_B <-ggplot(B_coord %>% filter(x<0.03,y<0.11), aes(x=x, y=y)) +
  stat_density_2d(geom = "polygon", contour = TRUE,
                  aes(fill = after_stat(level)), colour = "blue",
                  bins = 50) +xlab("")+ylab("")+xlab("x axis")+
  scale_fill_distiller(palette = "Blues", direction = 1) + theme_classic()+
  geom_point(data=B_centroid_info$centroid,aes(x=x,y=y),color="white",shape = 10,size=5)+
theme(legend.position= c(0.15, 0.2),legend.key.size = unit(0.2, 'cm'))+
  guides(fill=guide_legend(title="Density"),color=guide_legend(override.aes=list(fill=NA)))+
  annotate("text",label="I",x=0.018,y=0.105,size=8,family = c("sans"),fontface="bold" )


## . Eukaryota
P2DC_E <-ggplot(E_coord %>% filter(x<0.04,y>0.08), aes(x=x, y=y)) +
  stat_density_2d(geom = "polygon", contour = TRUE,
                  aes(fill = after_stat(level)), colour = "darkgreen",
                  bins = 50) +ylab("")+xlab("")+
  scale_fill_distiller(palette = "Greens", direction = 1) +  theme_classic()+
  geom_point(data=E_centroid_info$centroid,aes(x=x,y=y),color="white",shape = 10,size=5)+
theme(legend.position= c(0.8, 0.2),legend.key.size = unit(0.2, 'cm'))+
  guides(fill=guide_legend(title="Density"),color=guide_legend(override.aes=list(fill=NA)))+
  annotate("text",label="J",x=0.036,y=0.094,size=8,family = c("sans"),fontface="bold" )



 ## assembling the panels ----
# PlotTest(Object=P2DC_E,W=20,H=10,R=75)
P2D_9panels <- ggalign::align_plots(P2Dall_A,  P2Dall_B,    P2Dall_E,
                                 P2D_A,        P2D_B,       P2D_E,
                                 P2DC_A,       P2DC_B,      P2DC_E,
                                 ncol = 3)
png("plots/2D_analyses/P2D_9panels.png",units = "cm",width =30 ,height = 20,res = 300)
  print(P2D_9panels)
dev.off()


## . Polygons Areas -----
# functions ---------------------------------------
Calculate2DPolygonArea <- function(Polygon=A_polygon){
  require(sf) #install.packages("sf") 
  require(dplyr)
  # The first and last points must be identical to close the polygon for sf
  poly_coords_closed <- rbind(Polygon, Polygon[1,]) %>% as.matrix()
  # Create an sfg (simple feature geometry) polygon
  sfg_polygon <- st_polygon(list(poly_coords_closed))  
  # Create an sf (simple feature) object (optional, but common for data)
  sf_polygon <- st_sfc(sfg_polygon)
  # Calculate the area; Note: st_area() returns a units object
  polygon_area_sf <- st_area(sf_polygon)
  return(polygon_area_sf) # "Area of the polygon", 
}

A_area <- Calculate2DPolygonArea(Polygon=A_polygon) # 0.00022
B_area <- Calculate2DPolygonArea(Polygon=B_polygon) # 0.00124
E_area <- Calculate2DPolygonArea(Polygon=E_polygon) # 0.00401


# B_area/A_area # 5.6 x
# E_area/A_area # 18.1 x



