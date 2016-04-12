package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"flag"
	"github.com/antonholmquist/jason"
	"github.com/bogdansolomykin/vk_wrapper/vk"
	"github.com/jasonlvhit/gocron"
)

func task(authToken string, profileUrl string) {
	fmt.Println("===")
	fmt.Println(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Local())

	var jasonObject *jason.Object
	var jasonValueArray []*jason.Value

	// Получить ник пользователя
	parts := strings.Split(profileUrl, "/")
	userName := parts[len(parts)-1]

	api := vk.Api{
		AccessToken: authToken,
		UserId: "",
		ExpiresIn: "",
	}

	// Получить UID пользователя по nickname
	jasonObject, _ = jason.NewObjectFromBytes([]byte(api.Request(
		"users.get",
		map[string]string{
			"user_ids": userName,
		},
	)))

	jasonValueArray, _ = jasonObject.GetValueArray("response")
	jasonObject, _ = jasonValueArray[0].Object()
	userId, _ := jasonObject.GetInt64("uid")
	userIdString := strconv.FormatInt(userId, 10)

	fmt.Println("userId:", userIdString)

	// Получить общее количество записей на стене
	jasonObject, _ = jason.NewObjectFromBytes([]byte(api.Request(
		"wall.get",
		map[string]string{
			"owner_id": userIdString,
			"count": "1",
		},
	)))

	jasonValueArray, _ = jasonObject.GetValueArray("response")
	numberOfPosts, _ := jasonValueArray[0].Int64()

	fmt.Println("numberOfPosts:", numberOfPosts)

	// Получить случайный пост
	jasonObject, _ = jason.NewObjectFromBytes([]byte(api.Request(
		"wall.get",
		map[string]string{
			"owner_id": userIdString,
			"count": "1",
			"offset": strconv.Itoa(rand.Intn(int(numberOfPosts))),
		},
	)))

	jasonValueArray, _ = jasonObject.GetValueArray("response")
	jasonObject, _ = jasonValueArray[1].Object()
	postId, _ := jasonObject.GetInt64("id")
	postIdString := strconv.FormatInt(postId, 10)

	fmt.Println("postId:", postIdString)

	// Закрепить пост
	api.Request(
		"wall.pin",
		map[string]string{
			"owner_id": userIdString,
			"post_id": postIdString,
		},
	)

	fmt.Print("pinned post: ", profileUrl, "?w=wall", userIdString, "_", postIdString)
	fmt.Println()
	fmt.Println("===")
}

func main() {
	var (
		authTokenEnv string
		profileUrlEnv string
		intervalEnv uint64
		authToken string
		profileUrl string
		interval uint64
	)

	// Значения из окружения
	authTokenEnv = os.Getenv("VK_AUTH_TOKEN")
	profileUrlEnv = os.Getenv("VK_PROFILE_URL")
	intervalEnv, _ = strconv.ParseUint(os.Getenv("VK_SCHEDULER_INTERVAL_SECONDS"), 10, 64)

	// Значения из командной строки
	flag.StringVar(&authToken, "token", authTokenEnv, "VK authentication token")
	flag.StringVar(&profileUrl, "profile", profileUrlEnv, "VK profile URL (vk.com/user)")
	flag.Uint64Var(&interval, "time", intervalEnv, "Scheduler interval in seconds")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	gocron.Every(interval).Seconds().Do(task, authToken, profileUrl)
	<-gocron.Start()
}
