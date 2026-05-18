package main

import "time"

type Score struct {
	Runs    int     `json:"runs"`
	Wickets int     `json:"wickets"`
	Overs   float64 `json:"overs"`
	Innings string  `json:"innings"`
}

type Match struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	TeamA     string    `json:"teamA"`
	TeamB     string    `json:"teamB"`
	Venue     string    `json:"venue"`
	Score     []Score   `json:"score"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// The API response and MatchData structs are important because unlike in node js and express js you cannot extract the data directly from the json in Go so you need to create structs, everytime for a kind of data you are getting
type APIResponse struct {
	Status string      `json:"status"`
	Data   []MatchData `json:"data"`
}

type MatchData struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Venue  string   `json:"venue"`
	Teams  []string `json:"teams"`
	Score  []Score  `json:"score"`
}
