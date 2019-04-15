graph <- function(filename) {
  data <- read.csv('issues.dat', header = FALSE)
  dates <- lapply(data$V1, as.Date)
  sn <- cumsum(data$V2)
  
  plot(dates, sn, xaxt='n', main = "Cumulative bugs fixed over time", xlab = "", ylab = "Total number of issues resolved")
  axis.Date(1, at=seq(dates[[1]], dates[[301]], by="months"), format="%m-%Y")
}

export <- function(filename, width = 640, height = 480) {
  png(filename=filename, width = width, height = height)
  graph()
  dev.off()
  graph()
}

export('bugs.png')