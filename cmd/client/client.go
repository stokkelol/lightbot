package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const authHeader = "X-Client-Auth"

func main() {
	token := os.Getenv("AUTH_TOKEN")
	addr := os.Getenv("SERVER_ADDR")

	client := http.Client{}

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		req, err := http.NewRequest("GET", addr, nil)
		if err != nil {
			slog.Error("new request", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			continue
		}

		req.Header.Set(authHeader, token)
		if _, err = client.Do(req); err != nil {
			continue
		}
	}
}
