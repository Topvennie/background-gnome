package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func main() {
	fSave := flag.Bool("save", false, "Save the current image")
	flag.Parse()

	if *fSave {
		if err := save(); err != nil {
			fmt.Printf("Error saving image %v\n", err)
		}

		return
	}

	if err := update(); err != nil {
		fmt.Printf("Error updating background %v\n", err)
	}
}

func save() error {
	if c.dirSave == "" {
		return errors.New("no save directory configured")
	}

	entries, err := os.ReadDir(c.dirPath)
	if err != nil {
		return fmt.Errorf("reading path %s directory %v", c.dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		f, err := os.ReadFile(c.dirPath + entry.Name())
		if err != nil {
			return fmt.Errorf("read file %s %v", entry.Name(), err)
		}

		if err := os.WriteFile(c.dirSave+entry.Name(), f, os.ModePerm); err != nil {
			return fmt.Errorf("write file %s %v", entry.Name(), err)
		}
	}

	return nil
}

func update() error {
	if len(topics) == 0 {
		return errors.New("no topics configured")
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

	fmt.Printf("Chosen topic %s\n", chosenTopic.name)

	if len(chosenTopic.queries) == 0 {
		return errors.New("topic has no queries")
	}

	// Get image
	var data []byte

	for {
		chosenQuery := chosenTopic.queries[rand.IntN(len(chosenTopic.queries))]
		fmt.Printf("Chosen query %s\n", chosenQuery)

		var err error

		fmt.Println("Getting image")
		data, err = getImage(chosenQuery)
		if err != nil {
			return fmt.Errorf("getting image %v", err)
		}

		fmt.Println("Checking darkness")
		dark, err := isDark(data)
		if err != nil {
			return fmt.Errorf("checking image darkness %v", err)
		}

		if dark {
			break
		}

		fmt.Println("Image was not dark enough, getting another one...")
	}

	// Move / delete old image
	entries, err := os.ReadDir(c.dirPath)
	if err != nil {
		return fmt.Errorf("reading path %s directory %v", c.dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if c.dirOld == "" {
			if err := os.Remove(c.dirPath + entry.Name()); err != nil {
				return fmt.Errorf("deleting old file %s %v", entry.Name(), err)
			}
		} else {
			if err := os.Rename(c.dirPath+entry.Name(), c.dirOld+entry.Name()); err != nil {
				return fmt.Errorf("moving old file %s to %s %v", entry.Name(), c.dirOld+entry.Name(), err)
			}
		}
	}

	// Save image
	fileName := fmt.Sprintf("%s_%s.jpg", time.Now().Format("02_01_06_15_04_05"), strings.ToLower(strings.ReplaceAll(chosenTopic.name, " ", "_")))

	if err := os.WriteFile(c.dirPath+fileName, data, os.ModePerm); err != nil {
		return fmt.Errorf("writing image to disk %v", err)
	}

	if err := setBackground(c.dirPath + fileName); err != nil {
		return fmt.Errorf("setting background %v", err)
	}

	fmt.Println("Background updated")

	return nil
}
