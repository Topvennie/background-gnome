package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	accessKey = ""
	path      = ""
	old       = ""
)

func main() {
	// Choose a topic
	totalWeight := 0
	for _, t := range topics {
		totalWeight += t.weight
	}

	randomWeight := rand.IntN(totalWeight)
	currentWeight := 0

	var chosenTopic *topic
	for _, t := range topics {
		currentWeight += t.weight
		if currentWeight > randomWeight {
			chosenTopic = &t
			break
		}
	}

	if chosenTopic == nil {
		fmt.Printf("No topic found for weight %d\n", randomWeight)
		return
	}

	fmt.Printf("Chosen topic: %s\n", chosenTopic.name)

	if len(chosenTopic.queries) == 0 {
		fmt.Println("Topic has no queries")
		return
	}

	chosenQuery := chosenTopic.queries[rand.IntN(len(chosenTopic.queries))]

	fmt.Printf("Chosen query: %s\n", chosenQuery)

	// Get image
	fmt.Println("Getting image")
	data, err := getImage(chosenQuery)
	if err != nil {
		fmt.Printf("Error getting image %v", err)
		return
	}

	// Move / delete old image
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading path %s directory %v\n", path, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if old == "" {
			if err := os.Remove(path + entry.Name()); err != nil {
				fmt.Printf("Error deleting old file %s %v\n", entry.Name(), err)
				return
			}
		} else {
			if err := os.Rename(path+entry.Name(), old+entry.Name()); err != nil {
				fmt.Printf("Error moving old file %s to %s %v\n", entry.Name(), old+entry.Name(), err)
				return
			}
		}
	}

	// Save image
	fileName := fmt.Sprintf("%s_%s.png", strings.ToLower(strings.ReplaceAll(chosenTopic.name, " ", "_")), time.Now().Format("02_01_06_15_04_05"))

	if err := os.WriteFile(path+fileName, data, os.ModePerm); err != nil {
		fmt.Printf("Error writing image to disk %v\n", err)
		return
	}

	if err := setBackground(path + fileName); err != nil {
		fmt.Printf("Error setting background %v\n", err)
		return
	}

	fmt.Println("Background updated")
}

type apiImageResp struct {
	URLs struct {
		Raw string `json:"raw"`
	} `json:"urls"`
}

func getImage(query string) ([]byte, error) {
	params := url.Values{}
	params.Add("orientation", "landscape")
	params.Add("query", strings.ReplaceAll(query, " ", "-"))

	url := fmt.Sprintf("https://api.unsplash.com/photos/random?%s", params.Encode())

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create new http request %w", err)
	}
	req.Header.Set("Accept-Version", "v1")
	req.Header.Set("Authorization", "Client-ID "+accessKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong resp status code %s", resp.Status)
	}

	var random apiImageResp
	if err := json.NewDecoder(resp.Body).Decode(&random); err != nil {
		return nil, fmt.Errorf("decode resp %w", err)
	}

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s&fm=png", random.URLs.Raw), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create new http request %w", err)
	}
	req.Header.Set("Accept-Version", "v1")
	req.Header.Set("Authorization", "Client-ID "+accessKey)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong resp status code %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all bytes %w", err)
	}

	return data, nil
}

func setBackground(path string) error {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	colorSchemeArg := "picture-uri"
	if strings.TrimSpace(string(output)) == "'prefer-dark'" {
		colorSchemeArg = "picture-uri-dark"
	}

	cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", colorSchemeArg, "file://"+path)
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}
