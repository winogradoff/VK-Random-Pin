package main

import (
	"database/sql"
	"fmt"
	lib "github.com/winogradoff/go_vk_random_pin"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func worker() {
	var db *sql.DB
	var err error

	db, err = lib.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	token := os.Getenv("VK_TOKEN")
	username := os.Getenv("VK_USERNAME")
	delay, _ := strconv.ParseInt(os.Getenv("VK_DELAY"), 10, 64)

	api := lib.ApiVK{
		Token:    token,
		Username: username,
		Delay:    delay,
		Version:  lib.API_VERSION,
		APIUrl:   lib.API_METHOD_URL,
	}

	for {
		userId := api.GetUserId()
		numberOfPosts := api.GetNumberOfPosts(userId)
		postId := api.GetRandomPost(userId, numberOfPosts)
		api.PinPost(userId, postId)
		err = lib.InsertMessage(db, lib.Message{
			Time:          time.Now().UTC(),
			UserId:        userId,
			NumberOfPosts: numberOfPosts,
			PostId:        postId,
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * time.Duration(delay))
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	go worker()

	// Ждать сигнала завершения
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	fmt.Println("Signal received. Shutting down.")
}
