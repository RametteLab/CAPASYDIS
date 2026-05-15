# used libraries
require(dplyr)
require(ggplot2)
library(plotly)

dir.create("plots/3D_analyses")

nrow(D1) #1048576

Dtemp=D1[1:500000,] #max
plot_ly(Dtemp, 
        x = ~x, y = ~y, z = ~z,
        type = "scatter3d", 
        mode = "markers", 
        marker = list(
          size = 3,
          color = I(Dtemp$color)
      )
)
