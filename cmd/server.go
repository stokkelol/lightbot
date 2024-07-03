package main

import (
	"github.com/stokkelol/lightbot/pkg/bot"
	"github.com/stokkelol/lightbot/pkg/cache"
	"log"
	"net/http"
	"os"
	"time"
)

const authHeader = "X-Client-Auth"

func main() {

	token := os.Getenv("AUTH_TOKEN")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	timeCache := cache.NewCache()

	http.HandleFunc("/ping", func(rw http.ResponseWriter, req *http.Request) {
		header := req.Header.Get(authHeader)
		if header != token {
			http.Error(rw, "unauthorized", http.StatusUnauthorized)

			return
		}

		timeCache.SetLastTimestamp(time.Now())

		rw.WriteHeader(http.StatusOK)
	})

	b, err := bot.New(telegramToken, timeCache)
	if err != nil {
		log.Fatalf("new bot: %v", err)
	}

	go b.Run()

	if err = http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
