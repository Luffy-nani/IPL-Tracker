package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	StartFetcher()
	mux := http.NewServeMux()

	mux.HandleFunc("/matches", Chain(handleMatches))
	mux.HandleFunc("/matches/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/subscribe") {
			Chain(handleSubscribe)(w, r)
		} else if strings.HasSuffix(path, "/score") {
			Chain(handleScore)(w, r)
		} else {
			Chain(handleMatchByID)(w, r)
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//We used go because the ListenAndServe blocks forever so we used go func here
	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error: ", err)
		}
	}() //This () at last calls the func immediately..since our func is anonymous...it doest get to call so we used this
	//srv.ListenAndServe() --> Starts the http server
	// http.ErrServerClosed()..it helps to check if the server is shutdown gracefully...if the server is shutdown gracefully then there is no error..

	//Graceful shutdown, meaning --> Stop the server without  suddenly killing the users"
	// Graceful shutdown says --> "don't accept NEW requests,but let CURRENT requests finish"

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan //this blocks till we recieve a signal from os..it we recieve then we can continue with shutdown

	fmt.Println("\nShutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //give ongoing requests at most 5 seconds to finish
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Done!")

}
