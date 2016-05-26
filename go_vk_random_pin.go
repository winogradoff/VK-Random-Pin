package go_vk_random_pin

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	API_METHOD_URL = "https://api.vk.com/method/"
	API_VERSION    = "5.52"
	MESSAGES_SIZE  = 100 // Количество хранимых сообщений в БД
	TABLE_SQL      = `
		CREATE TABLE IF NOT EXISTS log_messages
		(
			message_id      serial     PRIMARY KEY,
			time            timestamp  NOT NULL DEFAULT now(),
			user_id         integer    NOT NULL DEFAULT 0,
			number_of_posts integer    NOT NULL DEFAULT 0,
			post_id         integer    NOT NULL DEFAULT 0
		);`
)

type Message struct {
	Time          time.Time
	UserId        int64
	NumberOfPosts int64
	PostId        int64
}

type ApiVK struct {
	Token    string
	Username string
	Delay    int64
	Version  string
	APIUrl   string
}

// Запрос к API
func (api *ApiVK) Request(methodName string, params map[string]string) []byte {
	values := url.Values{}
	values.Set("access_token", api.Token)
	values.Set("v", api.Version)
	for k, v := range params {
		values.Set(k, v)
	}
	response, _ := http.PostForm(api.APIUrl+methodName, values)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return body
}

// ID пользователя
func (api *ApiVK) GetUserId() int64 {
	response := api.Request("users.get", map[string]string{
		"user_ids": api.Username,
	})
	var result struct {
		Response []struct {
			UserId int64 `json:"id"`
		} `json:"response"`
	}
	json.Unmarshal(response, &result)
	return result.Response[0].UserId
}

// Количество записей на стене пользователя
func (api *ApiVK) GetNumberOfPosts(userId int64) int64 {
	response := api.Request("wall.get", map[string]string{
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

// Случайный пост
func (api *ApiVK) GetRandomPost(userId int64, numberOfPosts int64) int64 {
	response := api.Request("wall.get", map[string]string{
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
func (api *ApiVK) PinPost(userId int64, postId int64) {
	api.Request("wall.pin", map[string]string{
		"owner_id": strconv.FormatInt(userId, 10),
		"post_id":  strconv.FormatInt(postId, 10),
	})
}

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(TABLE_SQL)
	if err != nil {
		return err
	}
	return nil
}

func Connect() (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = CreateSchema(db)
	if err != nil {
		return nil, err
	}

	return db, err
}

func InsertMessage(db *sql.DB, m Message) error {
	var err error

	_, err = db.Exec(`
		INSERT INTO log_messages
		(user_id, number_of_posts, post_id) VALUES ($1, $2, $3);`,
		m.UserId, m.NumberOfPosts, m.PostId)

	if err != nil {
		return err
	}

	// Очистить всё, кроме MESSAGES_SIZE последних
	_, err = db.Exec(`
		DELETE FROM log_messages
		WHERE message_id NOT IN (
			SELECT message_id FROM log_messages
			ORDER BY time DESC LIMIT $1
		);`, MESSAGES_SIZE)

	return err
}

func FetchMessages(db *sql.DB) []Message {
	rows, err := db.Query(`
		SELECT time, user_id, number_of_posts, post_id
		FROM log_messages ORDER BY time DESC;`)

	if err != nil {
		log.Fatal(err)
	}

	messages := make([]Message, 0)

	location, _ := time.LoadLocation("Europe/Moscow")

	for rows.Next() {
		var m Message
		err = rows.Scan(&m.Time, &m.UserId, &m.NumberOfPosts, &m.PostId)
		if err != nil {
			log.Fatal(err)
		}
		m.Time = m.Time.In(location)
		messages = append(messages, m)
	}

	return messages
}
