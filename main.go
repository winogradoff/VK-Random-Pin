package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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
const API_VERSION = "5.52"

var (
	authToken  string
	profileUrl string
	seconds    int64
)

// Запрос к API ВКонтакте
func request(methodName string, params map[string]string) []byte {
	values := url.Values{}
	values.Set("access_token", authToken)
	values.Set("v", API_VERSION)
	for k, v := range params {
		values.Set(k, v)
	}
	response, _ := http.PostForm(API_METHOD_URL+methodName, values)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return body
}

// Получить id пользователя по ссылке на профиль
func getUserId() int64 {
	urlParts := strings.Split(profileUrl, "/")
	response := request("users.get", map[string]string{
		"user_ids": urlParts[len(urlParts)-1],
	})
	var result struct {
		Response []struct {
			UserId int64 `json:"id"`
		} `json:"response"`
	}
	json.Unmarshal(response, &result)
	return result.Response[0].UserId
}

// Получить количество записей на стене пользователя
func getNumberOfPosts(userId int64) int64 {
	response := request("wall.get", map[string]string{
		"owner_id": strconv.FormatInt(userId, 10),
		"count":    "1",
	})
	var result struct {
		Response struct {
			Count int64 `json:"count"`
		} `json:"response"`
	}
	json.Unmarshal(response, &result)
	return result.Response.Count
}

// Получить случайный пост
func getRandomPost(userId int64, numberOfPosts int64) int64 {
	response := request("wall.get", map[string]string{
		"owner_id": strconv.FormatInt(userId, 10),
		"offset":   strconv.FormatInt(rand.Int63n(numberOfPosts), 10),
		"count":    "1",
	})
	var result struct {
		Response struct {
			Items []struct {
				PostId int64 `json:"id"`
			} `json:"items"`
		} `json:"response"`
	}
	json.Unmarshal(response, &result)
	return result.Response.Items[0].PostId
}

// Закрепить пост
func pinPost(userId int64, postId int64) {
	request("wall.pin", map[string]string{
		"owner_id": strconv.FormatInt(userId, 10),
		"post_id":  strconv.FormatInt(postId, 10),
	})
}

func task() {
	userId := getUserId()
	numberOfPosts := getNumberOfPosts(userId)
	postId := getRandomPost(userId, numberOfPosts)
	pinPost(userId, postId)

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s\n", time.Now().UTC()))
	buffer.WriteString(fmt.Sprintf("uid: %d, posts: %d, pinned: %d\n", userId, numberOfPosts, postId))
	buffer.WriteString(fmt.Sprintf("url: %s?w=wall%d_%d\n\n", profileUrl, userId, postId))
	fmt.Print(buffer.String())
}

func main() {
	// Значения из окружения
	authToken = os.Getenv("VK_AUTH_TOKEN")
	profileUrl = os.Getenv("VK_PROFILE_URL")
	seconds, _ = strconv.ParseInt(os.Getenv("VK_SCHEDULER_INTERVAL_SECONDS"), 10, 64)

	// Значения из командной строки
	flag.StringVar(&authToken, "token", authToken, "VK authentication token")
	flag.StringVar(&profileUrl, "profile", profileUrl, "VK profile URL (vk.com/user)")
	flag.Int64Var(&seconds, "time", seconds, "Scheduler interval in seconds")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	for {
		task()
		time.Sleep(time.Second * time.Duration(seconds))
	}
}
