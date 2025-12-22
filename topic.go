package main

type topic struct {
	name    string
	weight  int
	queries []string
}

var topics []topic = []topic{
	{
		name:   "Animals",
		weight: 10,
		queries: []string{
			"dark animal",
			"dark animal cute",
			"dark animal fierce",
			"dark wild animal",
		},
	},
	{
		name:   "Nature",
		weight: 10,
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
		weight: 5,
		queries: []string{
			"space stars",
			"stars",
		},
	},
}
