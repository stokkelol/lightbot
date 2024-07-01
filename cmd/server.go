package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/stokkelol/lightbot/pkg/watcher"
	"log"
	"net/http"
	"os"
	"time"
)

const authHeader = "X-Client-Auth"

func main() {
	server := gin.Default()
	pprof.Register(server)

	token := os.Getenv("AUTH_TOKEN")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	cache := watcher.NewWatcher()

	server.GET("/ping", func(c *gin.Context) {
		header := c.GetHeader(authHeader)
		if header != token {
			c.String(http.StatusUnauthorized, "unauthorized")

			return
		}

		cache.SetLastTimestamp(time.Now())
		c.String(http.StatusOK, "ok")
	})

	updateRoute := fmt.Sprintf("/%s/get", telegramToken)
	server.GET(updateRoute, func(c *gin.Context) {
		if time.Now().Add(-60 * time.Second).After(time.Unix(cache.GetLastTimestamp(), 0)) {
			c.String(http.StatusOK, "new updates")

			return
		}

		c.String(http.StatusOK, "no updates")
	})

	bot := watcher.NewWatcher()
	go bot.Run()

	if err := server.Run(os.Getenv("SERVER_ADDR")); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
