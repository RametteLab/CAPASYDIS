# Title: 3D analysis of NR99 analysis using CAPASYDIS  -----
# .........................................................
# run "Analyses_NR99_Preparation.R" before

# R libraries ----
library(dplyr)
library(plotly)
library(grid)
library(geometry) # install.packages("geometry")
library(FNN) #install.packages("FNN")
library(patchwork)
library(forcats)
library(ggalign)

dir.create("plots")

### Figure q99hist ----
HA <-   ggplot(A_Table,aes(x=dist_to_centroid))+
  geom_histogram(fill="red",bins = 50)+ labs(x="")+
  theme_classic()+xlim(0,0.03)+ylab("")+
  annotate("segment", x = q99_A, xend = q99_A, y = 0, yend = Inf,lty=2,lwd=0.5)
  
HB <-  ggplot(B_Table,aes(x=dist_to_centroid))+
  geom_histogram(fill="blue",bins = 50)+ labs(x="")+
  theme_classic()+xlim(0,0.03)+ylab("Number of sequences")+
  annotate("segment", x = q99_B, xend = q99_B, y = 0, yend = Inf,lty=2,lwd=0.5)
  
HE <-  ggplot(E_Table,aes(x=dist_to_centroid))+
    geom_histogram(fill="darkgreen",bins = 50)+ labs(x="Distance to domain centroid")+
    theme_classic()+xlim(0,0.03)+ylab("")+
    annotate("segment", x = q99_E, xend = q99_E, y = 0, yend = Inf,lty=2,lwd=0.5)

Hq99 <- ggalign::align_plots(HA,  HB,    HE,
                     ncol = 1)  
# png("plots/q99hist.png",units = "cm",width =10 ,height = 15,res = 300)
#    print(Hq99)
# dev.off()
  
# some stats .......................
DPhylum %>%
  group_by(domain) %>% 
  summarize(unique_Phylum_count = n_distinct(Phylum))
            
DPhylum %>%
  group_by(domain) %>% filter(Class != "NA") %>% 
  summarize(unique_Class_count = n_distinct(Class)) 

# Plot all 3 domains together ----
## . Dots ---- Figure 4
plot_ly(DPhylum, x = ~x, y = ~y, z = ~z,color = ~domain,
        colors=c('#FF0000','#0000FF','#006400'),
        type = "scatter3d", mode = "markers", marker = list(size = 3)) # %>% 
  ## possible to add hovering but it becomes slow
  # add_trace(data = D, x = ~x, y = ~y, z = ~z,
  #           text = ~label,hoverinfo = 'text',showlegend = F) 

## . Convex hulls ----
plot_ly(All_3Dhull_points, x = ~x, y = ~y, z = ~z,color = ~domain,
        colors=c('#FF0000','#0000FF','#006400'),
        type = "scatter3d", mode = "markers", marker = list(size = 3)) %>%
  add_trace(data = A_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('pink',nrow(A_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3)  %>%
  add_trace(data = B_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('lightblue',nrow(B_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3)  %>%
  add_trace(data = E_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('darkgreen',nrow(E_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3) 


# 2) Plot each domain separately ----
## . Archaea ----

plot_ly(DPhylum %>% filter(domain=="Archaea"), x = ~x, y = ~y, z = ~z,color = ~domain,
        colors=c('#FF0000'),
        type = "scatter3d", mode = "markers", marker = list(size = 3))

plot_ly(A_3Dhull_points, x = ~x, y = ~y, z = ~z,color=~domain,
        colors='#FF0000',
        type = "scatter3d", mode = "markers", marker = list(size = 1)) %>%
  add_trace(data = A_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('pink',nrow(A_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3) 


## . Bacteria ----
plot_ly(DPhylum %>% filter(domain=="Bacteria"), x = ~x, y = ~y, z = ~z,color = ~domain,
        colors=c('#0000FF'),
        type = "scatter3d", mode = "markers", marker = list(size = 3))

plot_ly(B_3Dhull_points, x = ~x, y = ~y, z = ~z,color=~domain,
        colors='#0000FF',
        type = "scatter3d", mode = "markers", marker = list(size = 1)) %>%
  add_trace(data = B_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('lightblue',nrow(B_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3) 

## . Eukaryota ----
plot_ly(DPhylum %>% filter(domain=="Eukaryota"), x = ~x, y = ~y, z = ~z,color = ~domain,
        colors=c('#006400'),
        type = "scatter3d", mode = "markers", marker = list(size = 3))

plot_ly(E_3Dhull_points, x = ~x, y = ~y, z = ~z,color=~domain,
        colors='#006400',
        type = "scatter3d", mode = "markers", marker = list(size = 1)) %>%
  add_trace(data = E_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('darkgreen',nrow(E_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3) 


#.....................................................................  
# ## 3) Are asymdist distances associated with taxonomic information? ---- 
#.....................................................................  
# -see if the phyla overlap or not in 2d, 3d...
# phylum (and class) by distance => plot x= phyla by y= ranked distances (boxplot)
# - do that per axis

###   [Figure Px] .. X axis ----
Px <- DPhylum %>%
  arrange(desc(domain), Phylum) %>% 
  # Now reorder the factor levels of P based on 
  # their new order of appearance
  mutate(Phylum = fct_inorder(Phylum)) %>%
  ggplot(aes(x=x,y=Phylum, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - X axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Px.png",units = "cm",width =10 ,height = 15,res = 300)
  print(Px)
dev.off()
# DPhylum %>%
#   arrange(desc(domain), Phylum) %>%
#   # Now reorder the factor levels of P based on their new order of appearance
#   mutate(Phylum = fct_inorder(Phylum)) %>% 
#   ggplot(aes(x=x,y=Phylum, fill=domain))+
#   geom_boxplot()+xlab("asymdist - X axis") + xlim(0,0.13)+
#   theme_classic() + theme(axis.text.y = element_text(size = 5))

### [Figure Py].. Y axis ----
Py <- DPhylum %>%
  arrange(desc(domain), Phylum) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Phylum = fct_inorder(Phylum)) %>% 
  ggplot(aes(x=y,y=Phylum, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - Y axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Py.png",units = "cm",width =10 ,height = 15,res = 300)
print(Py)
dev.off()

# DPhylum %>%
#   arrange(desc(domain), Phylum) %>%
#   # Now reorder the factor levels of P based on their new order of appearance
#   mutate(Phylum = fct_inorder(Phylum)) %>% 
#   ggplot(aes(x=y,y=Phylum, fill=domain))+
#   geom_boxplot()+xlab("asymdist - Y axis") + xlim(0,0.13)+
#   theme_classic() + theme(axis.text.y = element_text(size = 5))


### [Figure Pz].. Z axis ----
Pz <- DPhylum %>%
  arrange(desc(domain), Phylum) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Phylum = fct_inorder(Phylum)) %>% 
  ggplot(aes(x=z,y=Phylum, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - Z axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Pz.png",units = "cm",width =10 ,height = 15,res = 300)
print(Pz)
dev.off()
# DPhylum %>%
#   arrange(desc(domain), Phylum) %>%
#   # Now reorder the factor levels of P based on their new order of appearance
#   mutate(Phylum = fct_inorder(Phylum)) %>% 
#   ggplot(aes(x=z,y=Phylum, fill=domain))+
#   geom_boxplot()+xlab("asymdist - Z axis") + xlim(0,0.13)+
#   theme_classic() + theme(axis.text.y = element_text(size = 5))

### [Figure P3d].. 3D distances and taxo ----
P3d <- DPhylum %>%
  arrange(desc(domain), Phylum) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Phylum = fct_inorder(Phylum)) %>% 
  ggplot(aes(x=dist3D,y=Phylum, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/P3d.png",units = "cm",width =10 ,height = 15,res = 300)
print(P3d)
dev.off()

# DPhylum %>%
#   arrange(desc(domain), Phylum) %>%
#   # Now reorder the factor levels of P based on their new order of appearance
#   mutate(Phylum = fct_inorder(Phylum)) %>% 
#   ggplot(aes(x=dist3D,y=Phylum, fill=domain))+
#   geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
#   theme_classic() + theme(axis.text.y = element_text(size = 5))

### [Figure Pdens].. with density lines ----
P1 <- DPhylum  %>%  # geom_point
  arrange(desc(domain), Phylum) %>%
  mutate(Phylum = fct_inorder(Phylum)) %>% 
  ggplot(aes(x=dist3D,y=Phylum, group=domain,colour=domain,fill=domain))+
  xlab("asymdist - 3D distances to origin") + 
  geom_point(size=1) +
  theme_classic() + theme(axis.text.y = element_text(size = 3))

P2 <- DPhylum %>%  #densities
  arrange(desc(domain), Phylum) %>%
  mutate(Phylum = fct_inorder(Phylum)) %>% 
  ggplot(aes(x=dist3D,colour=domain,fill=domain)) + 
  geom_density(alpha=0.5) + 
  theme_bw() +
  theme0(plot.margin = unit(c(1,0,-0.48,2.2),"lines")) 
  
Pdens <-   P2 + P1 + plot_layout(ncol = 1)
  png("plots/Pdens.png",units = "cm",width =10 ,height = 15,res = 300)
print(Pdens)
dev.off()

## . "Class" level ----
### [Figure Cx].. X axis ----
  # D %>% filter(Class == "NA") %>% count
Cx <-   DClass %>%
  arrange(desc(domain), Class) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=x,y=Class, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - X axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Cx.png",units = "cm",width =10 ,height = 15,res = 300)
print(Cx)
dev.off()
### [Figure Cy].. Y axis ----
Cy <-   DClass %>%
  arrange(desc(domain), Class) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=y,y=Class, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - Y axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Cy.png",units = "cm",width =10 ,height = 15,res = 300)
print(Cy)
dev.off()

### [Figure Cz].. Z axis ----
Cz <-   DClass%>%
  arrange(desc(domain), Class) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=z,y=Class, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - Z axis") + xlim(0,0.13)+
  theme_classic() + theme(axis.text.y = element_text(size = 5))
png("plots/Cz.png",units = "cm",width =10 ,height = 15,res = 300)
print(Cz)
dev.off()

### .. 3D distances and taxo ----
C3d <-   DClass %>%
  arrange(desc(domain), Class) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=dist3D,y=Class, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
  theme_classic() + theme(axis.text.y = element_text(size = 3))
png("plots/C3d.png",units = "cm",width =10 ,height = 15,res = 300)
print(C3d )
dev.off()
  # DClass %>%
  # arrange(desc(domain), Class) %>%
  # # Now reorder the factor levels of P based on their new order of appearance
  # mutate(Class = fct_inorder(Class)) %>% 
  # ggplot(aes(x=dist3D,y=Class, fill=domain))+
  # geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
  # theme_classic() + theme(axis.text.y = element_text(size = 2))

### .. 3D points with density lines ----
Cdens1 <- DClass  %>%  # geom_point
  arrange(desc(domain), Class) %>%
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=dist3D,y=Class, group=domain,colour=domain,fill=domain))+
  xlab("asymdist - 3D distances to origin") + 
  geom_point(size=0.5) +
  theme_classic() + theme(axis.text.y = element_text(size = 3))

Cdens2 <- DClass %>%  #densities
  arrange(desc(domain), Class) %>%
  mutate(Class = fct_inorder(Class)) %>% 
  ggplot(aes(x=dist3D,colour=domain,fill=domain)) + 
  geom_density(alpha=0.5) + 
  theme_bw() +
  theme0(plot.margin = unit(c(1,0,-0.48,2.2),"lines")) 

Cdens <- Cdens2 + Cdens1 + plot_layout(ncol = 1)
png("plots/Cdens.png",units = "cm",width =10 ,height = 15,res = 300)
print(Cdens)
dev.off()

## . "Order" level ----
### .. [Figure O3d].. 3D distances and taxo ----
# D %>% filter(Order == "NA") %>% count

O3d <- DOrder %>%
  arrange(desc(domain), Order) %>%
  # Now reorder the factor levels of P based on their new order of appearance
  mutate(Order = fct_inorder(Order)) %>% 
  ggplot(aes(x=dist3D,y=Order, fill=domain,group=domain))+
  geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
  theme_classic() + theme(axis.text.y = element_text(size = 3))
png("plots/O3d.png",units = "cm",width =10 ,height = 15,res = 300)
print(O3d)
dev.off()

# DOrder %>%
#   arrange(desc(domain), Order) %>%
#   # Now reorder the factor levels of P based on their new order of appearance
#   mutate(Order = fct_inorder(Order)) %>% 
#   ggplot(aes(x=dist3D,y=Order, fill=domain))+
#   geom_boxplot()+xlab("asymdist - 3D distances to origin") + 
#   theme_classic() + theme(axis.text.y = element_text(size = 2))

### . [Figure Odens] .. 3D points with density lines ----
Odens1 <- DOrder  %>%  # geom_point
  arrange(desc(domain), Order) %>%
  mutate(Order = fct_inorder(Order)) %>%
  ggplot(aes(x=dist3D,y=Order, group=domain,colour=domain,fill=domain))+
  xlab("asymdist - 3D distances to origin") +
  geom_point(size=0.5) +
  theme_classic() + theme(axis.text.y = element_text(size = 3))

Odens2 <- DOrder %>%  #densities
  arrange(desc(domain), Order) %>%
  mutate(Order = fct_inorder(Order)) %>%
  ggplot(aes(x=dist3D,colour=domain,fill=domain)) +
  geom_density(alpha=0.5) +
  theme_bw() +
  theme0(plot.margin = unit(c(1,0,-0.48,2.2),"lines"))
# 
Odens <- Odens2 + Odens1 + plot_layout(ncol = 1)
png("plots/Odens.png",units = "cm",width =10 ,height = 15,res = 300)
print(Odens)
dev.off()


#.....................................................................
## . 18S of Animalia ----
#.....................................................................
plot_ly(E_3Dhull_points, x = ~x, y = ~y, z = ~z,
        type = "scatter3d", mode = "markers", marker = list(size = 1,color='#006400',showscale=FALSE)) %>%
  add_trace(data = E_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('darkgreen',nrow(E_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3,showlegend = F)  %>%
  add_trace(data = data.frame(E_centroid_info$centroid), x = ~x, y = ~y, z = ~z,
            type = "scatter3d", mode = "markers", marker = list(size = 5,color="yellow"),showlegend = F) %>%
  add_trace(data = ANIMALs, x = ~x, y = ~y, z = ~z,
            type = "scatter3d", mode = "markers", marker = list(size = 2,color="red"),showlegend = F) 


# nrow(ANIMALs)/nrow(E_3Dcoord) # 36.02% of the E data


#.....................................................................
# 5) Centroids ----
#.....................................................................
# as the median in each dimension

## . 3D Plots with centroids----    
plot_ly(A_3Dhull_points, x = ~x, y = ~y, z = ~z,
        type = "scatter3d", mode = "markers", marker = list(size = 1,color='#FF0000',showscale=FALSE)) %>%
  add_trace(data = A_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('pink',nrow(A_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3,showlegend = F)  %>%
  add_trace(data = data.frame(A_centroid_info$centroid), x = ~x, y = ~y, z = ~z,
                        type = "scatter3d", mode = "markers", marker = list(size = 5),showlegend = F)  
#--
plot_ly(B_3Dhull_points, x = ~x, y = ~y, z = ~z,
        type = "scatter3d", mode = "markers", marker = list(size = 1,color='#0000FF',showscale=FALSE)) %>%
  add_trace(data = B_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('lightblue',nrow(B_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3,showlegend = F)  %>%
  add_trace(data = data.frame(B_centroid_info$centroid), x = ~x, y = ~y, z = ~z,
            type = "scatter3d", mode = "markers", marker = list(size = 5),showlegend = F)  
#--
plot_ly(E_3Dhull_points, x = ~x, y = ~y, z = ~z,
        type = "scatter3d", mode = "markers", marker = list(size = 1,color='#006400',showscale=FALSE)) %>%
  add_trace(data = E_3Dcoord, x = ~x, y = ~y, z = ~z,
            facecolor = rep('darkgreen',nrow(E_3Dcoord)*3) ,
            type = "mesh3d", opacity = 0.3,showlegend = F)  %>%
  add_trace(data = data.frame(E_centroid_info$centroid), x = ~x, y = ~y, z = ~z,
            type = "scatter3d", mode = "markers", marker = list(size = 5,color="yellow"),showlegend = F) 

## . between and within distances to centroids ----    


#.....................................................................
## . [Figure centroid_hist_boxplots.png] ----
#.....................................................................

## [Figure AvDistCentr barplot] 
# within-distances
# AvDistCentroid %>% ggplot(aes(x=domain, y=avg_distance_to_centroid))+
#     geom_bar(stat="identity",fill=c("darkgreen","red","blue"))+ theme_classic()

B1 <- AllCentroids_dist %>% ggplot(aes(x=factor(domain), y=avg_distance_to_centroid))+
  geom_errorbar(aes(ymin=avg_distance_to_centroid-0.001, ymax=avg_distance_to_centroid+sd),
                   width=.2)+                     # Width of the error bars
  geom_bar(stat="identity",fill=c("darkgreen","red","blue","darkgrey"))+ theme_classic()+  
  xlab("")+ylab("Distance to centroids")+ #ylim(0,0.15)+
  theme(axis.text.x = element_text(face="bold"),
        panel.grid.major = element_line(colour = "lightgrey"),
        panel.grid.major.x = element_blank(),
        panel.grid.minor.x = element_blank()
        )+
  scale_y_continuous(limits=c(0, 0.08),breaks=seq(0,0.08,0.01))

  # geom_text(aes(x=1:4, y=rep(-0.0025,4), 
  #               label=c(length(A_dist_to_centroid),
  #                       length(B_dist_to_centroid),
  #                       length(E_dist_to_centroid),
  #                       length(BetweenDist_Centroids)
  #                       )
  #           ), 
  #           col='black', size=3,fontface = "bold")+
      

png("plots/AvDistCentr_barplot.png",units = "cm",width =10 ,height = 15,res = 300)
  print(B1)
dev.off()
#.....................................................................
## . [Figure centroid_hist_boxplots.png] ----
#.....................................................................
R1 <- Raw_centroid_dists_tble %>% #
  ggplot(aes(x=Type, y=distance,fill=Type))+
  theme_classic()+ylab("Distance to centroid") +xlab("") +ylab("") +
  theme(legend.position="none",
        axis.text.x = element_text(face="bold"),
        panel.grid.major = element_line(colour = "lightgrey"),
        panel.grid.major.x = element_blank(),
        panel.grid.minor.x = element_blank(),
        axis.text.y=element_blank(),
        )+
  geom_jitter(aes(color=Type), width = 0.4,size=0.4)+
  scale_color_manual(values =c("red","blue","darkgreen","darkgrey"))+
  geom_boxplot(fatten=2,alpha=.01,colour="black")+ 
  scale_y_continuous(limits=c(0, 0.08),breaks=seq(0,0.08,0.01))

R2 <- Raw_centroid_dists_tble %>% filter(Type!="Between domains") %>% 
  ggplot(aes(x=distance,fill=Type))+
  geom_density(alpha=0.6)+ coord_flip()  + # xlim(0,0.10)+
  theme_classic()+ theme(legend.position="none")+
  scale_fill_manual(values =c("red","blue","darkgreen","darkgrey"))+
  ylab("density") +xlab("")+
  theme(axis.title.y=element_blank(),
        axis.text.y=element_blank(),
        axis.ticks.y=element_blank(),
        plot.margin = unit(c(1,1,1,-1.5),"lines")
    )+
  scale_x_continuous(limits=c(0, 0.08),breaks=seq(0,0.08,0.01))+
  geom_vline(xintercept=mean(BetweenDist_Centroids),linetype = "dotted")

R12 <- ggalign::align_plots(B1,R1,R2,Hq99, widths = c(2,2,1,2)) 

ggsave(plot = R12,filename = "plots/centroid_hist_boxplots_q99.png",
       dpi = 300,width = 15,height = 6)


#.....................................................................
## . [Figure AvDistCentr - Histograms] Histograms  ----
#.....................................................................
E_hist1 <- ggplot(data.frame(x=E_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "darkgreen") + theme_classic() + labs(x="distance to centroid")
E_hist2 <- ggplot(data.frame(x=E_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "darkgreen") + geom_density() + theme_bw() + xlim(0,0.01)  +
  labs(x="",y="")+ theme_classic()
png("plots/AvDistCentr_E_Hist.png",units = "cm",width =10 ,height = 15,res = 300)
  E_hist1 + inset_element(E_hist2, left = 0.25, bottom = 0.4, right = 0.9, top = 0.9)
dev.off()

A_hist1 <- ggplot(data.frame(x=A_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "red") + theme_classic() + labs(x="distance to centroid")
A_hist2 <- ggplot(data.frame(x=A_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "red") + geom_density() + theme_bw() + xlim(0,0.01)  +
  labs(x="",y="")+ theme_classic()
png("plots/AvDistCentr_A_Hist.png",units = "cm",width =10 ,height = 15,res = 300)
  A_hist1 + inset_element(A_hist2, left = 0.25, bottom = 0.4, right = 0.9, top = 0.9)
dev.off()


B_hist1 <- ggplot(data.frame(x=B_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "blue") + theme_classic() + labs(x="distance to centroid")
B_hist2 <- ggplot(data.frame(x=B_dist_to_centroid), aes(x=x, y=..density..)) + 
  geom_histogram(bins = 100,bg = "blue") + geom_density() + theme_bw() + xlim(0,0.01)  +
  labs(x="",y="density")+ theme_classic()
png("plots/AvDistCentr_B_Hist.png",units = "cm",width =10 ,height = 15,res = 300)
  B_hist1 + inset_element(B_hist2, left = 0.25, bottom = 0.4, right = 0.9, top = 0.9)
dev.off()


#.....................................................................
## . Why are there 2 peaks in the histogram for Bact? ---- 
#  Do they correspond to specific partitions of the lineages?
#.....................................................................
B_dist_to_centroid_names <- data.frame(cbind(label=BPhylum$label,dist_to_centroid=as.numeric(B_dist_to_centroid)))
# B_dist_to_centroid_names %>% head()
B_dist_to_centroid_names$dist_to_centroid <- as.numeric(B_dist_to_centroid_names$dist_to_centroid)

# B_hist2 + geom_vline( xintercept=0.0033)
Peak1 <-   B_dist_to_centroid_names %>% filter(B_dist_to_centroid<0.0033)
Peak2 <-   B_dist_to_centroid_names %>% filter(B_dist_to_centroid>=0.0033)
Peak1_mean <- mean(Peak1$dist_to_centroid)
Peak2_mean <- mean(Peak2$dist_to_centroid)

# aligning insets vertically
png("plots/DistCentr_aligned.png",units = "cm",width =10 ,height = 15,res = 300)
  A_hist2 + B_hist2 + E_hist2 + plot_layout(ncol = 1)+
  labs(x="distance to centroid")
dev.off()
  
## [Figure AvDistCentrBact] ----
png("plots/AvDistCentr_A.png",units = "cm",width =10 ,height = 15,res = 300)
  A_hist1 + inset_element(A_hist2, left = 0.25, bottom = 0.4, right = 0.9, top = 0.9)
dev.off()


AvDistCentr_B <- B_hist2 + geom_vline( xintercept=c(Peak1_mean,Peak2_mean),linewidth = 0.5,linetype = 2)+
  geom_vline( xintercept=0.0033,linewidth = 0.5,linetype = 1)+
  # geom_text(
    # aes(x, y, label = c(paste0("<- mean 1:",round(Peak1_mean,5)),paste0("<- mean 2:",round(Peak2_mean,5)))), 
    # data = data.frame(x=c(Peak1_mean,Peak2_mean), y=c(260,260)), 
    # hjust = -0.05, vjust = 0.5, size = 3  ) +
  ylim(0,300)+labs(y="density",x="Distance to centroid")+
  geom_text(
    aes(x, y, label = c("<- Group 1 ->","<- Group 2 ->")), 
    data = data.frame(x=c(0.0008,0.0035), y=c(80,80)), 
    hjust = -0.05, vjust = 0.5, size = 3,colour="white",fontface="bold"
  )

# png("plots/AvDistCentr_B.png",units = "cm",width =10 ,height = 15,res = 300)
# AvDistCentr_B
# dev.off()

## Taxo composition of each group and whether they differ ----
B_dist_to_centroid_names  <-  B_dist_to_centroid_names %>% 
  mutate(Peak=case_when(
    dist_to_centroid <= 0.0033 ~ "G1",
    dist_to_centroid > 0.0033 ~ "G2",
))
B_dist_to_centroid_names$Peak %>% table(.)  # how many sequences in each group
B_dist_to_centroid_names$Peak %>% table(.)/nrow(B_dist_to_centroid_names)


# B_dist_to_centroid_names %>% head()
B_dist_to_centroid_names$Phylum <- sapply(B_dist_to_centroid_names$label,function(x)strsplit(x,";")[[1]][2])
# B_dist_to_centroid_names %>% dim() # 290722      4
# B_dist_to_centroid_names %>% head

# by phyla
B_Table_Peaks_Phyla <- as.data.frame.matrix(t(table(B_dist_to_centroid_names$Peak,B_dist_to_centroid_names$Phylum)))
B_Table_Peaks_Phyla$total <- apply(B_Table_Peaks_Phyla,1,sum)
B_Table_Peaks_Phyla_filtered <- B_Table_Peaks_Phyla %>% filter(total>=100) %>% arrange(desc(total)) 

T_Phy <-as.matrix( B_Table_Peaks_Phyla_filtered[,1:2])
B_Perc_phy <- round(proportions(T_Phy,1)*100,1)
B_Table_Peaks_Percent_phyla <- data.frame(P1=B_Table_Peaks_Phyla_filtered[,1],P1pc=B_Perc_phy[,1],
                                          P2=B_Table_Peaks_Phyla_filtered[,2],P2pc=B_Perc_phy[,2],
                                          Total= B_Table_Peaks_Phyla_filtered[,3]
)
# B_Table_Peaks_Percent_phyla %>% head
B_Table_Peaks_Percent_phyla  %>% arrange(desc(P1pc),P2pc)

PHist1 <- B_Table_Peaks_Percent_phyla %>% 
  ggplot(aes(P1pc))+geom_histogram(bins=50,color="blue",fill="lightblue")+ 
  ylim(0,10)+theme_classic() + ylab("density")+xlim(0,100)+
  annotate("text", x = 25, y = 10,label="Phylum",fontface="bold")+
  xlab("Assignment (%) to Group 1")

PHist2 <- B_Table_Peaks_Percent_phyla %>% 
  ggplot(aes(P2pc))+geom_histogram(bins=50,color="blue",fill="lightblue")+ 
  ylim(0,10)+theme_classic() + ylab("")+
  xlab("Assignment (%) to Group 2")
PHist <- PHist1 + PHist2 + plot_layout(ncol=2)

# png("plots/PHist.png",units = "cm",width =15 ,height = 10,res = 300)
#   print(PHist)
# dev.off()
# B_Table_Peaks_Percent_phyla$phyla <- rownames(B_Table_Peaks_Percent_phyla)
# B_Table_Peaks_Percent_phyla  %>% head

B_Table_Peaks_Percent_phyla %>% filter(P1pc>=90 | P2pc>=90) %>% arrange(desc(P1pc)) %>% nrow #22 
# B_Table_Peaks_Percent_phyla %>% nrow
# 22/51 = 43.1% of the phyla were classified with >90% of the sequences in group 1 (21 cases) or group 2(1 case "Pseudomonadota")
B_Table_Peaks_Percent_phyla %>% filter(P1pc>=90 | P2pc>=90) %>% arrange(desc(P1pc))

#-------------------- Classes
B_dist_to_centroid_names$Class <- sapply(B_dist_to_centroid_names$label,function(x)strsplit(x,";")[[1]][3])
# table(B_dist_to_centroid_names$Phylum)
# table(B_dist_to_centroid_names$Class)
# B_dist_to_centroid_names %>% head

  B_Table_Peaks <- as.data.frame.matrix(t(table(B_dist_to_centroid_names$Peak,B_dist_to_centroid_names$Class)))
  B_Table_Peaks$total <- apply(B_Table_Peaks,1,sum)
  B_Table_Peaks_filtered <- B_Table_Peaks %>% filter(total>=100) %>% arrange(desc(total)) 
  
  T <-as.matrix( B_Table_Peaks_filtered[,1:2])
  B_Perc <- round(proportions(T,1)*100,1)
  # B_Perc %>% head
  # B_Table_Peaks_filtered %>% head
  B_Table_Peaks_Percent <- data.frame(P1=B_Table_Peaks_filtered[,1],P1pc=B_Perc[,1],
                                 P2=B_Table_Peaks_filtered[,2],P2pc=B_Perc[,2],
                                 Total= B_Table_Peaks_filtered[,3]
                                 )
  # B_Table_Peaks_Percent %>% head
  B_Table_Peaks_Percent  %>% arrange(desc(P1pc),P2pc)
  
  CHist1 <- B_Table_Peaks_Percent %>% ggplot(aes(P1pc))+
    geom_histogram(bins=50,color="darkgreen",fill="lightgreen")+  
    ylim(0,15)+theme_classic() + xlab("Assignment (%) to Group 1")+ylab("density")+
  annotate("text", x = 25, y = 15,label="Class",fontface="bold")
  
  CHist2 <- B_Table_Peaks_Percent %>% ggplot(aes(P2pc))+ylab("")+
    geom_histogram(bins=50,color="darkgreen",fill="lightgreen")+  ylim(0,15)+theme_classic() + xlab("Assignment (%) to Group 2")
  CHist <- CHist1 + CHist2 + plot_layout(ncol=2)
  
# png("plots/CHist.png",units = "cm",width =15 ,height = 10,res = 300)
#   print(CHist)
# dev.off()
  
  B_Table_Peaks_Percent %>% filter(P1pc>=90 | P2pc>=90) %>% arrange(desc(P1pc)) %>% nrow #42
  B_Table_Peaks_Percent %>% filter(P1pc>=90 ) %>% nrow #36 
  B_Table_Peaks_Percent %>% filter(P2pc>=90 ) %>% nrow #5 
  B_Table_Peaks_Percent %>% nrow #104
  # 42/104 = 40.4% of the classes were classified with >90% of the sequences in group 1 (36 cases) or group 2 (5 cases)
  


# assembling the figure
png("plots/AvDistCenter_Bact_P_C.png",units = "cm",width =15 ,height = 20,res = 300)
  AvDistCentr_B + PHist + CHist + plot_layout(ncol=1)
dev.off()


