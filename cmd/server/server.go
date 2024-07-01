package main

import (
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

	bot := watcher.NewWatcher()
	go bot.Run()

	if err := server.Run(os.Getenv("SERVER_ADDR")); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
