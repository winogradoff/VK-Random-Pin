package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type message struct {
	Time          time.Time
	UserId        int64
	NumberOfPosts int64
	PostId        int64
}

const API_METHOD_URL = "https://api.vk.com/method/"
const API_VERSION = "5.52"
const MESSAGES_SIZE = 100

var (
	token    string
	username string
	delay    int64
)

// Запрос к API ВКонтакте
func request(methodName string, params map[string]string) []byte {
	values := url.Values{}
	values.Set("access_token", token)
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
	response := request("users.get", map[string]string{
		"user_ids": username,
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

func main() {
	// Значения из окружения
	token = os.Getenv("VK_TOKEN")
	username = os.Getenv("VK_USERNAME")
	delay, _ = strconv.ParseInt(os.Getenv("VK_DELAY"), 10, 64)

	// Значения из командной строки
	flag.StringVar(&token, "token", token, "VK authentication token")
	flag.StringVar(&username, "username", username, "VK username")
	flag.Int64Var(&delay, "delay", delay, "Delay in seconds")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	messages := make([]message, 0)

	go func() {
		for {
			userId := getUserId()
			numberOfPosts := getNumberOfPosts(userId)
			postId := getRandomPost(userId, numberOfPosts)
			pinPost(userId, postId)

			messages = append(messages, message{
				Time:          time.Now().UTC(),
				UserId:        userId,
				NumberOfPosts: numberOfPosts,
				PostId:        postId,
			})

			for len(messages) > MESSAGES_SIZE {
				messages = messages[1:]
			}

			time.Sleep(time.Second * time.Duration(delay))
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"time":           time.Now().UTC(),
			"API_METHOD_URL": API_METHOD_URL,
			"API_VERSION":    API_VERSION,
			"MESSAGES_SIZE":  MESSAGES_SIZE,
			"Messages":       messages,
		})
	})

	router.Run(":" + os.Getenv("PORT"))
}
