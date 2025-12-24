package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func main() {
	if len(topics) == 0 {
		fmt.Println("No topics configured")
		return
	}

	// Choose a topic
	totalWeight := 0
	for _, t := range topics {
		totalWeight += t.weight
	}

	randomWeight := rand.IntN(totalWeight)
	currentWeight := 0

	chosenTopic := topics[0]
	for _, t := range topics {
		currentWeight += t.weight
		if currentWeight > randomWeight {
			chosenTopic = t
			break
		}
	}

	fmt.Printf("Chosen topic: %s\n", chosenTopic.name)

	if len(chosenTopic.queries) == 0 {
		fmt.Println("Topic has no queries")
		return
	}

	// Get image
	var data []byte

	for {
		chosenQuery := chosenTopic.queries[rand.IntN(len(chosenTopic.queries))]
		fmt.Printf("Chosen query: %s\n", chosenQuery)

		var err error

		fmt.Println("Getting image")
		data, err = getImage(chosenQuery)
		if err != nil {
			fmt.Printf("Error getting image %v\n", err)
			return
		}

		fmt.Println("Checking darkness")
		dark, err := isDark(data)
		if err != nil {
			fmt.Printf("Error checking image darkness %v\n", err)
			return
		}

		if dark {
			break
		}

		fmt.Println("Image was not dark enough, getting another one...")
	}

	// Move / delete old image
	entries, err := os.ReadDir(c.dirPath)
	if err != nil {
		fmt.Printf("Error reading path %s directory %v\n", c.dirPath, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if c.dirOld == "" {
			if err := os.Remove(c.dirPath + entry.Name()); err != nil {
				fmt.Printf("Error deleting old file %s %v\n", entry.Name(), err)
				return
			}
		} else {
			if err := os.Rename(c.dirPath+entry.Name(), c.dirOld+entry.Name()); err != nil {
				fmt.Printf("Error moving old file %s to %s %v\n", entry.Name(), c.dirOld+entry.Name(), err)
				return
			}
		}
	}

	// Save image
	fileName := fmt.Sprintf("%s_%s.jpg", time.Now().Format("02_01_06_15_04_05"), strings.ToLower(strings.ReplaceAll(chosenTopic.name, " ", "_")))

	if err := os.WriteFile(c.dirPath+fileName, data, os.ModePerm); err != nil {
		fmt.Printf("Error writing image to disk %v\n", err)
		return
	}

	if err := setBackground(c.dirPath + fileName); err != nil {
		fmt.Printf("Error setting background %v\n", err)
		return
	}

	fmt.Println("Background updated")
}
