# Title: 3D analysis of NR99 analysis using CAPASYDIS  -----
# focusing on E. coli and relatives
# .........................................................
# run "Analyses_NR99_Preparation.R" before

# R libraries ----
library(dplyr)
library(plotly)
library(geometry) # install.packages("geometry")
library(gridExtra)
library(ggalign)
library(patchwork) # inset_element()



# [Figure 3D Entero_XY] ----
## . Dots ---- 
plot_ly(Entero_3Dcoord, x = ~x, y = ~y, z = ~z,color = ~Genus,
        colors=c('#FF0000','#0000FF','#006400'),
        type = "scatter3d", mode = "markers", marker = list(size = 3)) 

# Family: 2D dots and marginals   ----
DF <- DEnterobacteriaceaeTop15 
GenusNamesTop15 <-  DF %>% select(Genus) %>% unique() %>% pull()

# library(ggpubr ) # for the legend
# 
# p0 <- ggplot(DF,aes(x=x,y=y,colour=factor(Genus))) + geom_point() +
#   theme_bw()  +labs(colour="Enterobacteriaceae genera")
# 
# leg0 <- get_legend(p0)
# as_ggplot(leg0)
# 
# p1 <- ggplot(DF,aes(x=x,y=y,colour=factor(Genus))) + geom_point() +
#   theme_bw() + 
#   theme(legend.position="none",plot.margin=unit(c(0,0,0,0),"points"))
# 
# p2 <- ggplot(DF,aes(x=x, colour=factor(Genus),fill=factor(Genus))) + 
#   geom_density(alpha=0.5) + 
#   # scale_x_continuous(breaks=NULL,expand=c(0.02,0)) +
#   # scale_y_continuous(breaks=NULL,expand=c(0.00,0)) +
#   theme_bw() +
#   theme0(plot.margin = unit(c(1,0,-0.48,2.2),"lines")) 
# 
# p3 <- ggplot(DF,aes(x=y, colour=factor(Genus),fill=factor(Genus))) + 
#   geom_density(alpha=0.5) + 
#   coord_flip()  + 
#   # scale_x_continuous(labels = NULL,breaks=NULL,expand=c(0.02,0)) +
#   # scale_y_continuous(labels = NULL,breaks=NULL,expand=c(0.00,0)) +
#   theme_bw() +
#   theme0(plot.margin = unit(c(0,1,1.2,-0.48),"lines"))
# 
# grid.arrange(arrangeGrob(p2,ncol=2,widths=c(3,1)),
#              arrangeGrob(p1,p3,ncol=2,widths=c(3,1)),
#              heights=c(1,3))

# Family Enterobacteriaceae: each axis with inset ----
## .[Figure Dens_x] ----
Dens_x <- ggplot(DF,aes(x=x, fill=factor(Genus))) + 
  geom_density(alpha=0.5) + labs(fill="Genera")+
  theme_classic()
# Dens_x
Dens_x_zoom <- Dens_x +xlim(0,0.0045)

Ex <- Dens_x_zoom  +theme(legend.position = "none") +labs(x="x axis",y="")
  #  inset_element(Dens_x, left = 0.25, bottom = 0.5, right = 0.99, top = 0.99)+
  # theme(legend.position="none")+labs(x="x axis",y="")

#Extracting colors for each genus
# Ex 
gEx <- ggplot_build(Ex)
Info_genus <- data.frame(
        Genus=levels(factor(DF$Genus)),
        Color=unique(gEx$data[[1]]["fill"])
)

## .[Figure Dens_y] ----
Dens_y <- ggplot(DF,aes(x=y, fill=factor(Genus))) + 
  geom_density(alpha=0.5) + labs(fill="Genera")+
  theme_classic()
# Dens_y

Dens_y_zoom <- Dens_y +xlim(0.0985,0.1015)

Ey <- Dens_y_zoom + labs(x="y axis",y="density")
  # inset_element(Dens_y, left = 0.25, bottom = 0.5, right = 0.99, top = 0.99)+
  # theme(legend.position="none")+

## .[Figure Dens_z] ----
Dens_z <- ggplot(DF,aes(x=z, fill=factor(Genus))) + 
  geom_density(alpha=0.5) + labs(fill="Genera")+
  theme_classic()
# Dens_z

Dens_z_zoom <- Dens_z +xlim(0.0525,0.056)

Ez <-Dens_z_zoom  + theme(legend.position="none")+labs(x="z axis",y="")
  # inset_element(Dens_z, left = 0.25, bottom = 0.5, right = 0.99, top = 0.99)+
  # theme(legend.position="none")+labs(x="",y="")

P_EnteroMarginals <- ggalign::align_plots(Ex,  Ey,    Ez,
                                          ncol = 1)

# PlotTest(P_EnteroMarginals,W = 18,H = 20,R=300)

dir.create("plots/EcolinFriends/")
png("plots/EcolinFriends/P_EnteroMarginals.png",units = "cm",width =18 ,height = 20,res = 300)
     print(P_EnteroMarginals)
dev.off()


## . density on each axis
PlotDensity1Genus=function(Name="Escherichia",Yscale=5000,colour="green"){
  Ex <- ggplot(DF %>% filter(Genus==Name), aes(x=x)) + 
    geom_density(alpha=0.5,fill=colour) + 
    theme_classic()+ylim(0,Yscale)+labs(title=Name)+theme(plot.title = element_text(face = "italic"))
  
  Ey <- ggplot(DF %>% filter(Genus==Name),
               aes(x=y)) + labs(y="")+
    geom_density(alpha=0.5,fill=colour) + ylim(0,Yscale)+
    theme_classic()
  
  
  Ez <- ggplot(DF %>% filter(Genus==Name),
               aes(x=z)) + labs(y="")+ylim(0,Yscale)+
    geom_density(alpha=0.5,fill=colour) + 
    theme_classic()
  
  P <- Ex+Ey+Ez+plot_layout(nrow = 1)  
  plot(P)
  
}

# PlotDensity1Genus(Name="Escherichia",Yscale=c(0,6000),colour="lightblue")
# PlotDensity1Genus(Name="Pluralibacter",Yscale=c(0,12000),colour="lightblue")

# Genera: Density on each single axis


Info_genus$Yscale <- c(2000,2100,4000,
                       2100,2000,6000,
                       6000,2200,4000,
                       6000,2000,4000,
                       2100,2200,6200
)

i= 1; a1 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2]) 
i= 2; a2 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 3; a3 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 4; a4 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 5; a5 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 6; a6 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 7; a7 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 8; a8 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 9; a9 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 10; a10 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 11; a11 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 12; a12 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 13; a13 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 14; a14 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])
i= 15; a15 <- PlotDensity1Genus(Name=Info_genus[i,1],Yscale=Info_genus[i,3],colour=Info_genus[i,2])


All_a <- ggalign::align_plots(
             a6,a1,a2,
             a3,a4,a5,
             a7,a8,a9,
             a10,a11,a12,
             a13,a14,a15, 
            ncol = 3)
png("plots/EcolinFriends/Enterobacteriaceae_genera_by_axis.png",units = "cm",width =50 ,height = 30,res = 300)
  print(All_a)
dev.off()


#Most extreme outlier in 3D: 
# AOQL01000076.469.2061__Bacteria;Pseudomonadota;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia-Shigella;Escherichia
B_Eco %>% filter(isEcoli,Aboveq99)


PEco_A <- B_Eco %>%  ggplot(aes(x=x,y=y))+
  xlim(0,0.0150)+ ylim(0.095,0.11) + theme_minimal()+xlab("")+ylab("")+
  geom_point(data=B_Eco %>% filter(!isEcoli),aes(x=x,y=y),color="grey")+
  geom_point(data=B_Eco %>% filter(isEcoli),aes(x=x,y=y),color="blue")+
  geom_polygon(data=Eco_polygon,aes(x=x,y=y),colour="blue",fill=NA)+
  # annotate("text",label="A",x=0.0005,y=0.110,size=7,family = c("sans"),fontface="bold" )+
  geom_point(data=B_Eco %>% filter(isEcoli,Aboveq99),aes(x=x,y=y),color="black",pch=15,size=2)+
  annotate("rect",xmin=0,xmax=0.0025,ymin=0.0985,ymax=0.1015, fill=NA,color="black",lty=2)

PEco_B <-ggplot(Eco_coord, aes(x=x, y=y)) +
  stat_density_2d(geom = "polygon", contour = TRUE,
                  aes(fill = after_stat(level)), colour = "blue",
                  bins = 50,contour_var = "density") +xlab("")+ylab("")+
  scale_fill_distiller(palette = "Blues", direction = 1) +
  # theme_classic()+xlim(-0.0005,0.0025)+ ylim(0.0990,0.1015) +
  theme_classic()+xlim(-0.0001,0.0011)+ ylim(0.0990,0.1001) +
  # annotate("text",label="B",x=0,y=0.1001,size=7,family = c("sans"),fontface="bold" )+
  theme(legend.position= c(0.85, 0.3),legend.key.size = unit(0.2, 'cm'))

## assembling the panels ----
# PlotTest(Object=EcoC_E,W=20,H=10,R=75)
Eco_2panels <- ggalign::align_plots(PEco_A,  PEco_B,  
                                    ncol = 1)
# png("plots/EcolinFriends/Eco_2panels.png",units = "cm",width =12 ,height = 20,res = 300)
#   print(Eco_2panels)
# dev.off()





# phylogenetics of E. coli ----
# see file 05_Ecoli_phylogenetics.R
  