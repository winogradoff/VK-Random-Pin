package main

import (
	"flag"
	"fmt"
	"github.com/elgs/gojq"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const API_METHOD_URL = "https://api.vk.com/method/"

var (
	authToken  string
	profileUrl string
	interval   uint64
)

func vkRequest(methodName string, params map[string]string) *gojq.JQ {
	u, _ := url.Parse(API_METHOD_URL + methodName)
	q := u.Query()
	q.Set("access_token", authToken)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	resp, _ := http.Get(u.String())
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	parser, _ := gojq.NewStringQuery(string(content))
	return parser
}

// Получить id пользователя по ссылке на профиль
func getUserId() int64 {
	urlParts := strings.Split(profileUrl, "/")
	userId, _ := vkRequest(
		"users.get",
		map[string]string{
			"user_ids": urlParts[len(urlParts)-1],
		},
	).QueryToInt64("response.[0].uid")
	return userId
}

// Получить количество записей на стене пользователя
func getNumberOfPosts(userId int64) int64 {
	numberOfPosts, _ := vkRequest(
		"wall.get",
		map[string]string{
			"owner_id": strconv.FormatInt(userId, 10),
			"count":    "1",
		},
	).QueryToInt64("response.[0]")
	return numberOfPosts
}

// Получить случайный пост
func getRandomPost(userId int64, numberOfPosts int64) int64 {
	postId, _ := vkRequest(
		"wall.get",
		map[string]string{
			"owner_id": strconv.FormatInt(userId, 10),
			"offset":   strconv.FormatInt(rand.Int63n(numberOfPosts), 10),
			"count":    "1",
		},
	).QueryToInt64("response.[1].id")
	return postId
}

// Закрепить пост
func pinPost(userId int64, postId int64) {
	vkRequest("wall.pin", map[string]string{
		"owner_id": strconv.FormatInt(userId, 10),
		"post_id":  strconv.FormatInt(postId, 10),
	})
}

func task() {
	fmt.Println("===")
	fmt.Println(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Local())

	userId := getUserId()
	numberOfPosts := getNumberOfPosts(userId)
	postId := getRandomPost(userId, numberOfPosts)
	pinPost(userId, postId)

	fmt.Println("userId:", userId)
	fmt.Println("numberOfPosts:", numberOfPosts)
	fmt.Println("postId:", postId)
	fmt.Print("pinned post: ", profileUrl, "?w=wall", userId, "_", postId)
	fmt.Println()
	fmt.Println("===")
}

func main() {
	// Значения из окружения
	authToken = os.Getenv("VK_AUTH_TOKEN")
	profileUrl = os.Getenv("VK_PROFILE_URL")
	interval, _ = strconv.ParseUint(os.Getenv("VK_SCHEDULER_INTERVAL_SECONDS"), 10, 64)

	// Значения из командной строки
	flag.StringVar(&authToken, "token", authToken, "VK authentication token")
	flag.StringVar(&profileUrl, "profile", profileUrl, "VK profile URL (vk.com/user)")
	flag.Uint64Var(&interval, "time", interval, "Scheduler interval in seconds")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	gocron.Every(interval).Seconds().Do(task)
	<-gocron.Start()
}
