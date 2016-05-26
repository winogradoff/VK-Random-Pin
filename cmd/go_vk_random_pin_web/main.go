package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	lib "github.com/winogradoff/go_vk_random_pin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const API_METHOD_URL = "https://api.vk.com/method/"
const API_VERSION = "5.52"
const MESSAGES_SIZE = 100

func main() {
	// Подключение к БД
	db, err := lib.Connect()
	if errd != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"time":           time.Now().UTC(),
			"API_METHOD_URL": API_METHOD_URL,
			"API_VERSION":    API_VERSION,
			"MESSAGES_SIZE":  MESSAGES_SIZE,
			"Messages":       lib.FetchMessages(db),
		})
	})

	router.Run(":" + os.Getenv("PORT"))
}
