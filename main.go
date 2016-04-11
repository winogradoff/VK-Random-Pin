package main

import (
	"os"
	"log"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"github.com/antonholmquist/jason"
	"github.com/bogdansolomykin/vk_wrapper/vk"
	"github.com/jasonlvhit/gocron"
)

var (
    InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func task() {
	InfoLogger.Println("===")

	InfoLogger.Println(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Local())

	var jasonObject *jason.Object
	var jasonValueArray []*jason.Value

	api := vk.Api{
		AccessToken: os.Getenv("VK_AUTH_TOKEN"),
		UserId: "",
		ExpiresIn: "",
	}

	profileUrl := os.Getenv("VK_PROFILE_URL")
	parts := strings.Split(profileUrl, "/")
	userName := parts[len(parts)-1]

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

	InfoLogger.Println("userId:", userIdString)

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

	InfoLogger.Println("numberOfPosts:", numberOfPosts)

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

	InfoLogger.Println("postId:", postIdString)

	// Закрепить пост
	api.Request(
		"wall.pin",
		map[string]string{
			"owner_id": userIdString,
			"post_id": postIdString,
		},
	)

	InfoLogger.Print("pinned post: ", profileUrl, "?w=wall", userIdString, "_", postIdString)
	InfoLogger.Println("===")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	gocron.Every(5).Seconds().Do(task)
	<-gocron.Start()
}
