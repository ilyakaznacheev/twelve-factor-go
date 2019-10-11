package main

import (
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! It's %s", time.Now().Format(time.Kitchen))
		fmt.Fprintf(w, "\nRequest UUID %s", uuid.NewV4())
		if envName != "" {
			fmt.Fprintf(w, "\nRunning on  %s", envName)
		}
	})
	http.ListenAndServe(":"+port, nil)
}
