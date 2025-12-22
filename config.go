package main

type config struct {
	// API
	apiAccessKey string
	// Directory
	dirPath string
	dirOld  string
	// Image darknes detection
	darkEnabled bool
	darkSize    int
	darkMedian  float64
	darkP90     float64
}

var c config = config{
	apiAccessKey: "",
	dirPath:      "",
	dirOld:       "",
	darkEnabled:  true,
	darkSize:     512,
	darkMedian:   0.35,
	darkP90:      0.65,
}
