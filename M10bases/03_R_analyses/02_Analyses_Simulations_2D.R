# used libraries
require(dplyr)
require(ggplot2)
# library(geometry) # install.packages("geometry")
library(patchwork) # install.packages("patchwork")

dir.create("plots/2D_analyses")


#only 2D because already unique coordinates

Plot=function(Object,W=20,H=10,R=300,FileName="test"){
  FilePath=paste0("plots/",FileName)
  png(FilePath,units = "cm",width =W ,height = H,res = R)
  print(Object)
  dev.off()
}
##########################################
# 5' data ---------
##########################################
head(D1)
nrow(D1) #1048576
summary(D1 %>%  select(x,y))
  # x                   y            
  # Min.   :0.0000000   Min.   :0.0000000  
  # 1st Qu.:0.2691991   1st Qu.:0.2895374  
  # Median :0.3415536   Median :0.3566115  
  # Mean   :0.3422929   Mean   :0.3558861  
  # 3rd Qu.:0.4146800   3rd Qu.:0.4229330  
  # Max.   :0.7293353   Max.   :0.6691824 

# -----------------------------------
# Nmut distribution

df <- D1 #[1:10000,]
Nmax <- max(df$Nmut, na.rm = TRUE)

NmutDensities_X <- ggplot(df, aes(x = x,y=Nmut)) +
  scale_fill_identity() +
  # Density curves: distribution of x per color
  geom_density(aes(y = after_stat(..density..) ,
                   color = color),
               linewidth = 0.6) +
  scale_color_identity() +
  labs(x = "x", y = "Frequency") +
  theme_light()  + xlab("")+
  theme(
    panel.grid.major.y = element_blank(),#Minor horizontal grid lines removed.
    panel.grid.minor.y = element_blank(),#Minor horizontal grid lines removed.
    axis.text.x  = element_blank(),  # remove x tick labels
    axis.ticks.x = element_blank(),   # remove x tick marks
    plot.margin = margin(t = 2, r = 5, b = -1, l = 5, unit = "pt")
  )

Plot(NmutDensities_X,FileName="Forward_NmutDensities_X.png")

# -----------------------------------
NmutDensities_Y <- ggplot(df, aes(x = y,y=Nmut)) +
  # Bars: Nmut by x, filled by color
  # geom_col(aes(y = Nmut, fill = color), color = NA) +
  scale_fill_identity() +
  
  # Density curves: distribution of x per color
  geom_density(aes(y = after_stat(..density..) ,
                   color = color),
               linewidth = 0.6) +
  scale_color_identity() +
  labs(x = "x", y = "Frequency") +
  theme_light()  + xlab("")+
  theme(
    panel.grid.major.x = element_blank(),#Minor horizontal grid lines removed.
    panel.grid.minor.x = element_blank(),#Minor horizontal grid lines removed.
    axis.text.y  = element_blank(),  # remove y tick labels
    axis.ticks.y = element_blank(),   # remove y tick marks
    plot.margin = margin(t = 2, r = 5, b = -1, l = 5, unit = "pt")
  )+
  coord_flip()

Plot(NmutDensities_Y,FileName="Forward_NmutDensities_Y.png")
# -----------------------------------
# all dots
Cols_unique_pairs <- unique(df[c("color", "Nmut")])

Fullplot <- df%>%  
  ggplot(aes(x=x,y=y,color=color,label=label))+
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
  ) +
  theme(
    plot.margin = margin(t = 0, r = 5, b = 2, l = 5, unit = "pt")
  )+
  theme(
    legend.position        = "inside",
    legend.position.inside = c(0.35, 0.25),   # x, y in [0,1] (top‑right corner)
    legend.justification   = c("right", "top"),
    legend.text       = element_text(size = 7),          # smaller text
    legend.title      = element_text(size = 8),          # smaller title
    legend.key.height = unit(0.25, "lines"),             # shorter keys
    legend.key.width  = unit(0.5,  "lines"),             # narrower keys
    legend.spacing.y  = unit(0.1,  "lines"),             # less vertical space
    legend.spacing.x  = unit(0.1,  "lines")              # less horizontal space
  )+
  theme(
    legend.box.background = element_rect(
      color = "black",   # border color
      linewidth  = 0.5,       # border line width
      fill  = NA         # or a color like "white" if you want a filled box
    ),
    legend.box.margin = margin(4, 4, 4, 4)  # a bit of padding inside the box
  )

Plot(Fullplot,FileName="Forward_all.png")

#combined dots (bottom) + density curves (top)
design <- "
AB
CD
"
FullCombined <- NmutDensities_X + guide_area() +
  Fullplot + NmutDensities_Y +
  plot_layout(
    design  = design,
    heights = c(1, 4),   # top vs bottom
    widths  = c(4, 1),    # left (main) vs right (y‑margin)
    guides  = "collect"
  ) &
  theme(legend.position = "right")
Plot(FullCombined,FileName="Forward_all_dots_densities.png")

# -----------------------------------
# top 100 (zooming-in)
N=100
df <- D1[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p100_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p100_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p100 <- p100_top / p100_main + plot_layout(heights = c(1, 4))
Plot(p100,FileName="Forward_top100.png")

# -----------------------------------
# top 1000
N=1000
df <- D1[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p1000_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p1000_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p1000 <- p1000_top / p1000_main + plot_layout(heights = c(1, 4))
Plot(p1000,FileName="Forward_top1000.png")

# -----------------------------------
# top 10000
N=10000
df <- D1[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p10000_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p10000_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p10000 <- p10000_top / p10000_main + plot_layout(heights = c(1, 4))
Plot(p10000,FileName="Forward_top10000.png")

##########################################
# 3' data ---------
##########################################
head(D1rev)
nrow(D1rev) #1048576
summary(D1rev %>%  select(x,y))
# x                   y            
# Min.   :0.0000000   Min.   :0.0000000  
# 1st Qu.:0.2691991   1st Qu.:0.2895374  
# Median :0.3415536   Median :0.3566115  
# Mean   :0.3422929   Mean   :0.3558861  
# 3rd Qu.:0.4146800   3rd Qu.:0.4229330  
# Max.   :0.7293353   Max.   :0.6691824 

# -----------------------------------
# Nmut distribution
# PlotNmut <-   

# png("plots/Reverse_all_Nmut.png",res = 300,width = 30,height = 10, units = "cm")
df <- D1rev #[1:10000,]
Nmax <- max(df$Nmut, na.rm = TRUE)

NmutDensities_X <- ggplot(df, aes(x = x,y=Nmut)) +
  # Bars: Nmut by x, filled by color
  # geom_col(aes(y = Nmut, fill = color), color = NA) +
  scale_fill_identity() +
  
  # Density curves: distribution of x per color
  geom_density(aes(y = after_stat(..density..) ,
                   color = color),
               linewidth = 0.6) +
  scale_color_identity() +
  labs(x = "x", y = "Frequency") +
  theme_light()  + xlab("")+
  theme(
    panel.grid.major.y = element_blank(),#Minor horizontal grid lines removed.
    panel.grid.minor.y = element_blank(),#Minor horizontal grid lines removed.
    axis.text.x  = element_blank(),  # remove x tick labels
    axis.ticks.x = element_blank(),   # remove x tick marks
    plot.margin = margin(t = 2, r = 5, b = -1, l = 5, unit = "pt")
  )

Plot(NmutDensities_X,FileName="Reverse_NmutDensities_X.png")

# -----------------------------------
NmutDensities_Y <- ggplot(df, aes(x = y,y=Nmut)) +
  # Bars: Nmut by x, filled by color
  # geom_col(aes(y = Nmut, fill = color), color = NA) +
  scale_fill_identity() +
  
  # Density curves: distribution of x per color
  geom_density(aes(y = after_stat(..density..) ,
                   color = color),
               linewidth = 0.6) +
  scale_color_identity() +
  labs(x = "x", y = "Frequency") +
  theme_light()  + xlab("")+
  theme(
    panel.grid.major.x = element_blank(),#Minor horizontal grid lines removed.
    panel.grid.minor.x = element_blank(),#Minor horizontal grid lines removed.
    axis.text.y  = element_blank(),  # remove y tick labels
    axis.ticks.y = element_blank(),   # remove y tick marks
    plot.margin = margin(t = 2, r = 5, b = -1, l = 5, unit = "pt")
  )+
  coord_flip()

Plot(NmutDensities_Y,FileName="Reverse_NmutDensities_Y.png")
# -----------------------------------
# all dots
Cols_unique_pairs <- unique(df[c("color", "Nmut")])

Fullplot <- df%>%  
  ggplot(aes(x=x,y=y,color=color,label=label))+
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
  ) +
  theme(
    plot.margin = margin(t = 0, r = 5, b = 2, l = 5, unit = "pt")
  )+
  theme(
    legend.position        = "inside",
    legend.position.inside = c(0.35, 0.25),   # x, y in [0,1] (top‑right corner)
    legend.justification   = c("right", "top"),
    legend.text       = element_text(size = 7),          # smaller text
    legend.title      = element_text(size = 8),          # smaller title
    legend.key.height = unit(0.25, "lines"),             # shorter keys
    legend.key.width  = unit(0.5,  "lines"),             # narrower keys
    legend.spacing.y  = unit(0.1,  "lines"),             # less vertical space
    legend.spacing.x  = unit(0.1,  "lines")              # less horizontal space
  )+
  theme(
    legend.box.background = element_rect(
      color = "black",   # border color
      linewidth  = 0.5,       # border line width
      fill  = NA         # or a color like "white" if you want a filled box
    ),
    legend.box.margin = margin(4, 4, 4, 4)  # a bit of padding inside the box
  )

Plot(Fullplot,FileName="Reverse_all.png")

#combined dots (bottom) + density curves (top)
design <- "
AB
CD
"
FullCombined <- NmutDensities_X + guide_area() +
  Fullplot + NmutDensities_Y +
  plot_layout(
    design  = design,
    heights = c(1, 4),   # top vs bottom
    widths  = c(4, 1),    # left (main) vs right (y‑margin)
    guides  = "collect"
  ) &
  theme(legend.position = "right")
Plot(FullCombined,FileName="Reverse_all_dots_densities.png")

# -----------------------------------
# top 100 (zooming-in)
N=100
df <- D1rev[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p100_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p100_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p100 <- p100_top / p100_main + plot_layout(heights = c(1, 4))
Plot(p100,FileName="Reverse_top100.png")

# -----------------------------------
# top 1000
N=1000
df <- D1rev[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p1000_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p1000_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p1000 <- p1000_top / p1000_main + plot_layout(heights = c(1, 4))
Plot(p1000,FileName="Reverse_top1000.png")

# -----------------------------------
# top 10000
N=10000
df <- D1rev[1:N, ]
xrange <- range(df$x, na.rm = TRUE)
yrange <- range(df$y, na.rm = TRUE)

p10000_top <- 
  ggplot(df, aes(x = x, y = Nmut, fill = color)) +
  geom_col(width = 0.002) +
  scale_fill_identity() +
  coord_cartesian(xlim = xrange) +
  theme_minimal() +
  theme(
    axis.title.x = element_blank(),
    axis.text.x = element_blank(),
    axis.ticks.x = element_blank(),
    panel.grid.minor = element_blank(),
    legend.position = "none"
  ) +
  ylab("Nmut")

p10000_main <- ggplot(df, aes(x = x, y = y, color = color, label = label)) +
  geom_point() +
  scale_color_identity() +
  coord_cartesian(xlim = xrange, ylim = yrange) +
  theme_minimal()

p10000 <- p10000_top / p10000_main + plot_layout(heights = c(1, 4))
Plot(p10000,FileName="Reverse_top10000.png")



