package main

import (
	"encoding/json" //for converting the json to go format
	"fmt"
	"net/http"
	"os" //to get the api key in env
	"time"
)

func getAPIURL() string {
	apiKey := os.Getenv("API_KEY")
	return "https://api.cricapi.com/v1/currentMatches?apikey=" + apiKey + "&offset=0"
}

func fetchMatches() error {
	resp, err := http.Get(getAPIURL())
	if err != nil {
		return fmt.Errorf("failed to fetch matches: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "success" {
		return fmt.Errorf("API error: %s", apiResp.Status)
	}

	var matches []Match
	for _, md := range apiResp.Data {
		if len(md.Teams) < 2 {
			continue
		}
		match := Match{
			ID:        md.ID,
			Title:     md.Name,
			Status:    md.Status,
			TeamA:     md.Teams[0],
			TeamB:     md.Teams[1],
			Venue:     md.Venue,
			Score:     md.Score,
			UpdatedAt: time.Now(),
		}
		matches = append(matches, match)
	}

	cache.Set(matches)
	fmt.Printf("Cache updated: %d matches fetched\n", len(matches))
	return nil
}

func StartFetcher() {
	go func() { // we used goroutine because we dont want the fetcher func to block our server requestes...it should be running concurrently
		if err := fetchMatches(); err != nil {
			fmt.Println("Initial fetch error:", err)
		}
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := fetchMatches(); err != nil {
				fmt.Println("Fetch error:", err)
			}
		}
	}()
}
