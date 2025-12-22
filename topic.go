package main

type topic struct {
	name    string
	weight  int
	queries []string
}

var topics []topic = []topic{
	{
		name:   "Animals",
		weight: 20,
		queries: []string{
			"dark animal",
			"dark animal cute",
			"dark animal fierce",
			"dark wild animal",
		},
	},
	{
		name:   "Nature",
		weight: 20,
		queries: []string{
			"dark nature",
			"dark nature cute",
			"dark nature mountain",
			"dark nature cliffs",
			"dark nature forest",
			"twilight nature",
		},
	},
	{
		name:   "Space",
		weight: 10,
		queries: []string{
			"space stars",
			"stars",
			"deep space",
			"nebula dark",
		},
	},
	{
		name:   "City",
		weight: 5,
		queries: []string{
			"city night",
			"dark city night",
			"city skyline night",
			"rainy city night",
		},
	},
	{
		name:   "Mist",
		weight: 10,
		queries: []string{
			"foggy landscape",
			"misty mountains",
			"foggy forest road",
			"dark fog landscape",
		},
	},
	{
		name:   "Rain",
		weight: 10,
		queries: []string{
			"rain night",
			"storm clouds dark",
			"rainy window night",
			"thunderstorm clouds",
		},
	},
	{
		name:   "Textures",
		weight: 5,
		queries: []string{
			"dark texture",
			"concrete wall dark",
			"black stone texture",
			"dark fabric texture",
		},
	},
	{
		name:   "Ocean",
		weight: 15,
		queries: []string{
			"ocean night",
			"dark ocean waves",
			"stormy sea",
			"night coastline",
		},
	},
}
