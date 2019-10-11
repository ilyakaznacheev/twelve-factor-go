package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	uuid "github.com/satori/go.uuid"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	envName := os.Getenv("ENV_NAME")
	timeAPIurl := os.Getenv("TIME_API")
	tc := newTimeClient(timeAPIurl)

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := tc.getTime()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hello World! It's %s", t.Format(time.Kitchen))
		fmt.Fprintf(w, "\nRequest UUID %s", uuid.NewV4())
		if envName != "" {
			fmt.Fprintf(w, "\nRunning on %s", envName)
		}
	})

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: router,
	}

	// handle SIGTERM
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Server is ready to handle requests at", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
	<-done
	log.Println("Server stopped")
}

type timeResponse struct {
	Time time.Time `json:"datetime"`
}

type timeClient struct {
	url string
}

func newTimeClient(url string) timeClient {
	return timeClient{url}
}

func (c timeClient) getTime() (*time.Time, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData timeResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return &respData.Time, nil
}
