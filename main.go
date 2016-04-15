package main

import (
	"flag"
	"fmt"
	"github.com/bogdansolomykin/vk_wrapper/vk"
	"github.com/elgs/gojq"
	"github.com/jasonlvhit/gocron"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func task(authToken string, profileUrl string) {
	fmt.Println("===")
	fmt.Println(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Local())

	// Получить ник пользователя
	parts := strings.Split(profileUrl, "/")
	userName := parts[len(parts)-1]

	api := vk.Api{
		AccessToken: authToken,
		UserId:      "",
		ExpiresIn:   "",
	}

	var json string
	var parser *gojq.JQ

	// Получить UID пользователя по nickname
	json = api.Request(
		"users.get",
		map[string]string{
			"user_ids": userName,
		},
	)
	parser, _ = gojq.NewStringQuery(json)
	userId, _ := parser.QueryToInt64("response.[0].uid")
	fmt.Println("uid:", userId)

	// Получить общее количество записей на стене
	json = api.Request(
		"wall.get",
		map[string]string{
			"owner_id": strconv.FormatInt(userId, 10),
			"count":    "1",
		},
	)
	parser, _ = gojq.NewStringQuery(json)
	numberOfPosts, _ := parser.QueryToInt64("response.[0]")
	fmt.Println("numberOfPosts:", numberOfPosts)

	// Получить случайный пост
	json = api.Request(
		"wall.get",
		map[string]string{
			"owner_id": strconv.FormatInt(userId, 10),
			"count":    "1",
			"offset":   strconv.FormatInt(rand.Int63n(numberOfPosts), 10),
		},
	)
	parser, _ = gojq.NewStringQuery(json)
	postId, _ := parser.QueryToInt64("response.[1].id")
	fmt.Println("postId:", postId)

	// Закрепить пост
	api.Request(
		"wall.pin",
		map[string]string{
			"owner_id": strconv.FormatInt(userId, 10),
			"post_id":  strconv.FormatInt(postId, 10),
		},
	)

	fmt.Print("pinned post: ", profileUrl, "?w=wall", userId, "_", postId)
	fmt.Println()
	fmt.Println("===")
}

func main() {
	// Значения из окружения
	authToken := os.Getenv("VK_AUTH_TOKEN")
	profileUrl := os.Getenv("VK_PROFILE_URL")
	interval, _ := strconv.ParseUint(os.Getenv("VK_SCHEDULER_INTERVAL_SECONDS"), 10, 64)

	// Значения из командной строки
	flag.StringVar(&authToken, "token", authToken, "VK authentication token")
	flag.StringVar(&profileUrl, "profile", profileUrl, "VK profile URL (vk.com/user)")
	flag.Uint64Var(&interval, "time", interval, "Scheduler interval in seconds")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	gocron.Every(interval).Seconds().Do(task, authToken, profileUrl)
	<-gocron.Start()
}
