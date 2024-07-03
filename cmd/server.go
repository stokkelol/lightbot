package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/stokkelol/lightbot/pkg/bot"
	"github.com/stokkelol/lightbot/pkg/cache"
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

	timeCache := cache.NewCache()

	server.GET("/ping", func(c *gin.Context) {
		header := c.GetHeader(authHeader)
		if header != token {
			c.String(http.StatusUnauthorized, "unauthorized")

			return
		}

		timeCache.SetLastTimestamp(time.Now())

		c.String(http.StatusOK, "ok")
	})

	b, err := bot.New(telegramToken, timeCache)
	if err != nil {
		log.Fatalf("new bot: %v", err)
	}

	go b.Run()

	if err = server.Run("0.0.0.0:9000"); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
