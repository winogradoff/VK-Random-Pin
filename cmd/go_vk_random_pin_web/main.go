package main

import (
	"github.com/gin-gonic/gin"
	lib "github.com/winogradoff/go_vk_random_pin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func webServer() {
	db, err := lib.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	location, _ := time.LoadLocation("Europe/Moscow")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"time":           time.Now().In(location),
			"API_METHOD_URL": lib.API_METHOD_URL,
			"API_VERSION":    lib.API_VERSION,
			"MESSAGES_SIZE":  lib.MESSAGES_SIZE,
			"Messages":       lib.FetchMessages(db),
		})
	})

	router.Run(":" + os.Getenv("PORT"))
}

func main() {
	go webServer()

	// Ждать сигнала завершения
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	log.Println("Signal received. Shutting down.")
}
