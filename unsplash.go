package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
	req.Header.Set("Authorization", "Client-ID "+c.apiAccessKey)

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

	highestRes, err := highestMonitorResolution()
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s&fm=jpg&q=100&cs=srgb&w=%d&h=%d&fit=crop", random.URLs.Raw, highestRes.Width, highestRes.Height), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create new http request %w", err)
	}
	req.Header.Set("Accept-Version", "v1")
	req.Header.Set("Authorization", "Client-ID "+c.apiAccessKey)

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
