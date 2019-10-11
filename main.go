package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	http.ListenAndServe(":"+port, nil)
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
