package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func handleMatches(w http.ResponseWriter, r *http.Request) {
	matches, _ := cache.Get()
	writeJSON(w, http.StatusOK, matches)
}

func handleMatchByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/matches/")
	// trim anything after the id (like /score or /subscribe)
	id = strings.Split(id, "/")[0]

	match, found := cache.GetByID(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "match not found"})
		return
	}
	writeJSON(w, http.StatusOK, match)
}

func handleScore(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/matches/")
	id = strings.Split(id, "/")[0]

	match, found := cache.GetByID(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "match not found"})
		return
	}
	writeJSON(w, http.StatusOK, match.Score)
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/matches/")
	id = strings.Split(id, "/")[0]

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := ps.Subscribe(id)
	defer ps.Unsubscribe(id, ch)

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming not supported"})
		return
	}

	for {
		select {
		case match := <-ch:
			data, _ := json.Marshal(match)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, updatedAt := cache.Get()
	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"cacheAge": time.Since(updatedAt).String(),
	})
}
