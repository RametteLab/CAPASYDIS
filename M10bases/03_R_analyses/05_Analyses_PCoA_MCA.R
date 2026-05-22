# The idea is to provide an alternative view to CAPASYDIS,
# using one-hot encoding of the MSA, and submitting that to SVD/MCA dimension reduction
# Note: Computing the MCA or SVD of this small feature matrix is trivial. 
# However, once new sequences are added as additional rows (while keeping the
# same feature columns), the MCA/SVD values based on the enlarged matrix 
# change, because the sample distribution and covariance structure have changed;
# the leading components will therefore shift accordingly. 

options(digits = 10) # 
# used libraries -----------------------------
library(seqinr)        # version 4.2-44
library(ggplot2)       # version 4.0.3
library(dplyr)         # version 1.2.1
library(vegan)         # version 2.7-3
library(FactoMineR)    # version 2.14
library(ape)           # version 5.8.1
library(ggdendro)      # version 0.2.0

options(digits = 10)

Plot=function(Object,W=20,H=10,R=300,FileName="test.png"){
  FilePath=paste0("plots/",FileName)
  png(FilePath,units = "cm",width =W ,height = H,res = R)
  print(Object)
  dev.off()
}

#getting the colors and Nmut data
DATE="20260513"
load(paste0(DATE,"_data.RData"))
D1_100 <- D1[1:100,]
Cols_unique_pairs <- unique(D1_100[c("color", "Nmut")])

# Getting the mldist values -----------------------------
# Visualization of the MLDIST data (ML data based on classical phylogenetics)
MLdf <- read.table("03_phylogenetics/iqtree/iqtree.mldist", header = FALSE, skip = 1, stringsAsFactors = FALSE)
# Extract the sequence names
MLseq_names <- MLdf[, 1]
# Create a distance matrix
MLdist_matrix <- as.matrix(as.dist(MLdf[, -1]))
# Assign sequence names to the matrix
rownames(MLdist_matrix) <- colnames(MLdist_matrix)<- MLseq_names
# MLdist_matrix[1:3,1:3]

# Hierarchical clustering -------------------------
hc <-hclust(as.dist(MLdist_matrix) )
hc_short <- hc
hc_short$height <- log(hc$height+1) #to shorten the visual ploting of branch length
dendro_data_raw <- dendro_data(hc_short, type = "rectangle")
tree_labels <- label(dendro_data_raw)
tree_colored_labels <- tree_labels %>%
  left_join(D1_100[, c("label", "color")], by = c("label" = "label"))

Phylogram <- ggplot() +
  # Draw the tree branches
  geom_segment(data = segment(dendro_data_raw), 
               aes(x = x, y = y, xend = xend, yend = yend), 
               color = "grey40") +
  
  # Draw the text labels using the matched colors from your data frame
  geom_text(data = tree_colored_labels, 
            aes(x = x, y = y, label = label, color = color), 
            hjust = -0.1,   # Negative value pushes text rightward, away from the tips
            vjust = 0.5,    # Centers text vertically on the tip line
            size = 1.5) +   #
  
  # Tell ggplot to use the literal color strings from your data frame
  scale_color_identity() +
  
  # Flip coordinates so it reads left-to-right horizontally
  coord_flip() +
  scale_x_continuous(expand = expansion(mult = c(0.02, 0.02))) +
  scale_y_reverse(expand = expansion(mult = c(0.01, 0.08))) +
  
  # Clear the canvas completely
  theme_void() +
  theme(
    panel.grid.major = element_blank(),
    panel.grid.minor = element_blank()
  )

Plot(Phylogram,FileName="Phylogram.png")

# unrooted representation
phylo_tree <- as.phylo(hc_short)
tip_metadata <- data.frame(label = phylo_tree$tip.label) %>%
  left_join(D1_100[, c("label", "color","Nmut")], by = "label")

png("plots/unrooted_tree.png")
plot(phylo_tree, 
     type = "unrooted",                     # Forces the unrooted/radial layout
     tip.color = tip_metadata$color,  # Applies your sequence-specific colors
     cex = 0.5,                       # Font size of your text labels
     # lab4ut = "radial",               # Rotates text to follow the angle of the branches
     no.margin = TRUE,                # Removes base R layout margins
     label.offset = 0.02,             # Pushes labels slightly away from branch tips
     edge.color = "grey40",           # Color of the tree branches
     edge.width = 1.2                # Thickness of the tree branches
)
tiplabels(pch = 21,                    # Shape 21 allows both a fill and a border
          bg = tip_metadata$color,     # Fill color from your vector
          col = "white",               # Border color around the dot (makes them pop)
          cex = 1.2,                   # Size of the dot
          lwd = 0.5)                   # Border line width

legend_info <- tip_metadata %>%
  select(Nmut, color) %>%
  distinct() %>%
  arrange(Nmut) # Keeps it alphabetical

legend(x = "bottomleft",                       # Placement position ("topleft", "bottomright", etc.)
       legend = legend_info$Nmut,             # Text labels for the legend entries
       pt.bg = legend_info$color,              # Color fills for the dots
       col = "white",                          # Border color for the dots (matches tiplabels)
       pch = 21,                               # Matches the circle shape used in tiplabels
       pt.cex = 1.3,                           # Scale size of the legend dots
       cex = 0.8,                              # Font size of the legend text
       bty = "n",                              # "n" removes the harsh black border box around the legend
       title = "Number of mutations",              # Title text for the legend box
       title.font = 2)

dev.off()


# PCoA (Metric MDS) ----------------------
# Tries to maximize the linear, exact distances between points in a low-dimensional space.
# using wcmdscale from vegan library to use  Lingoes correction for negative eigenvalues
pcoa_res <- vegan::wcmdscale(MLdist_matrix, eig = TRUE, add="lingoes")
# Extract coordinates for plotting
pcoa_coords <- as.data.frame(pcoa_res$points)
colnames(pcoa_coords) <- c("PCoA1", "PCoA2","PCoA3")
# Calculate variance explained by the axes
var_explained <- round(pcoa_res$eig / sum(pcoa_res$eig) * 100, 1)
cat("Variance explained:\n Axis 1 =", var_explained[1], 
    "%\n Axis 2 =",var_explained[2],
    "%\n Axis 3 =", var_explained[3], "%\n")
  # Variance explained:
  # Axis 1 = 2.1 %
  # Axis 2 = 2 %
  # Axis 3 = 1.5 %

# plot
# Combine your coordinates and colors into the same data frame
pcoa_coords$color <- D1_100$color
pcoa_coords$Nmut <- D1_100$Nmut
pcoa_coords$label <- rownames(pcoa_coords)
# Build the ggplot
PCOA_plot <- ggplot(pcoa_coords, aes(x = PCoA1, y = PCoA2,color=color)) +
  geom_point(aes(color = color), size = 2, show.legend = TRUE) +
  geom_text(aes(label = label, color = color), 
            hjust = -0.2, vjust = 0.5, size = 2, show.legend = FALSE) +
  labs(
    x = paste0("PCoA1 (", var_explained[1], "%)"),
    y = paste0("PCoA2 (", var_explained[2], "%)"),
    title = "" 
  ) +
  labs(
    color = "Number of mutations"
  )+
  theme_light()+
theme(panel.grid.major = element_blank(), # Removes major grid lines
              panel.grid.minor = element_blank()  # Removes minor grid lines
  )+
  scale_fill_identity(
    name = "color",
    breaks = pcoa_coords$color,
    labels = pcoa_coords$Nmut,
    guide = "legend"
  ) +
  scale_color_identity(
    name = "Number of mutations",
    breaks = pcoa_coords$color,
    labels = pcoa_coords$Nmut,
    guide = "legend"
) 
if(!dir.exists("plots/alternatives")) {dir.create("plots/alternatives")}
Plot(PCOA_plot,FileName="alternatives/PCoA.png")


# MCA (Multiple Correspondence Analysis) ------------------------------
# the most important step with binary data is deciding whether or not to center and scale the data before running the SVD.
# MCA is simply a standard SVD executed on a specially transformed version of your categorical data
# input for MCA is 01 data, so we will convert the sequences in the MSA using one-hot encoding 
# of the sequences (while keeping the length of the sequences [MSA] constant).

## Read FASTA MSA 
fasta_file <- "03_phylogenetics/MSAtop100.fasta"  # path to the MSA
aln <- seqinr::read.alignment(file = fasta_file, format = "fasta")  #

## aln$seq is a character vector of aligned sequences
seqs <- unlist(aln$seq)
seqs <- toupper(seqs)
L <- unique(nchar(seqs)) # nber of positions in the MSA (here 10)

## One-hot encoding function (per sequence) 
alphabet=c("A", "C", "G", "T")
one_hot_encode <- function(seq, alphabet=alphabet) {
  bases <- strsplit(seq, "")[[1]] 
  mat <- matrix(0, nrow = length(bases), ncol = length(alphabet),
                dimnames = list(NULL, alphabet))
  for (i in seq_along(bases)) {
    b <- bases[i]
    if (b %in% alphabet) {
      mat[i, b] <- 1
    } else {
      ## handle gaps/ambiguous bases: here we leave the row as all zeros
    }
  }
  as.vector(t(mat))
}

## Build the sequence × feature matrix  
p <- length(alphabet) * L

X <- t(vapply(seqs, one_hot_encode,
              FUN.VALUE = numeric(p),
              alphabet = alphabet))
rownames(X) <- aln$nam
X <- as.data.frame(X)
# here we have to remove the empty columns each time!!!
# columns where every single row is a 1 or 0  will also crash an MCA because 
# they have zero variance. 
X1 <-  X %>%
  select(where(function(col) {
    num_col <- as.numeric(as.character(col))
    # Keep only if the sum is greater than 0 AND not equal to the total number of rows
    sum(num_col) > 0 & sum(num_col) != length(num_col)
}))
ncol(X) #40
ncol(X1) #14

#converting the columns to factors
X2 <- X1%>%
  mutate(across(everything(), as.factor))

# apply MCA
mca_res<- MCA(X2,
               ncp=2,
               graph = FALSE)

mca_eig <- mca_res$eig

MCA_coords <-data.frame(mca_res$ind$coord )
MCA_coords[,1:2] <- MCA_coords[,1:2]%>%
  mutate(across(everything(), as.numeric))
colnames(MCA_coords) <- c("MCA1","MCA2")
MCA_coords$color <- D1_100$color
MCA_coords$Nmut <- D1_100$Nmut
MCA_coords$label <- rownames(MCA_coords)

head(MCA_coords)
range(MCA_coords$MCA1)
range(MCA_coords$MCA2)
MCA_plot <- ggplot(MCA_coords, aes(x=MCA1,y=MCA2,color=color)) +
    geom_point(aes(color = color), size = 2, show.legend = TRUE) +
    geom_text_repel(
      aes(label = label, color = color),
      size = 2,                  # Font size of the text
      max.overlaps = Inf,        # Forces ggplot to show ALL labels, no matter how crowded
      box.padding = 0.1,         # Adds breathing room around each text box
      point.padding = 0.01,       # Adds distance between the point and its text label
      segment.color = "grey60",  # Color of the tiny connector lines
      segment.size = 0.3,        # Thickness of the connector lines
      show.legend = FALSE
    )  +
    xlim(-0.8,0.8)+ylim(-0.6,0.8)+
  labs(
    x = paste0("MCA Dim 1 (",round(mca_eig[1,2],2),"%)"),
    y = paste0("MCA Dim 2 (",round(mca_eig[2,2],2),"%)"),
    title = "" 
  ) +
    theme_light()+
  labs(
    color = "Number of mutations"
  )+
  theme_light()+
  theme(panel.grid.major = element_blank(), # Removes major grid lines
        panel.grid.minor = element_blank()  # Removes minor grid lines
  )+
  scale_fill_identity(
    name = "color",
    breaks = MCA_coords$color,
    labels = MCA_coords$Nmut,
    guide = "legend"
  ) +
  scale_color_identity(
    name = "Number of mutations",
    breaks = MCA_coords$color,
    labels = MCA_coords$Nmut,
    guide = "legend"
) 

Plot(MCA_plot,FileName="alternatives/MCA.png")

# uniqueness of values
length(unique(MCA_coords[,1])) #64  
length(unique(MCA_coords[,2])) #98  

#unique coordinate
MCA_coords %>% select(MCA1,MCA2) %>% unique() %>% nrow() #100

MCA_coords %>% filter(MCA1>0.5 & MCA1<0.55) %>% select(MCA1) %>% range()

MCA_plot +   xlim(0.4,0.6)+ylim(-0.225,-0.2)

# all possible digits
table(MCA_coords[,1]) %>% length()  #21 values mathematically diff, 
table(MCA_coords[,2])%>% length()  #51

#at 10^-10
table(round(MCA_coords[,1],10)) %>% length()  #9 
table(round(MCA_coords[,2],10))%>% length()  #9
