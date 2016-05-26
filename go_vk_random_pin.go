package go_vk_random_pin

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	API_METHOD_URL = "https://api.vk.com/method/"
	API_VERSION    = "5.52"
	MESSAGES_SIZE  = 100
	TABLE_SQL      = `
		CREATE TABLE IF NOT EXISTS log_messages
		(
			message_id      serial     default,
			time            timestamp  NOT NULL DEFAULT now(),
			user_id         integer    NOT NULL DEFAULT 0,
			number_of_posts integer    NOT NULL DEFAULT 0,
			post_id         integer    NOT NULL DEFAULT 0,
			PRIMARY KEY (message_id, user_id, post_id)
		);`
)

type Message struct {
	Time          time.Time
	UserId        int64
	NumberOfPosts int64
	PostId        int64
}

func CreateSchema(db *pg.DB) error {
	_, err := db.Exec(TABLE_SQL)
	if err != nil {
		return err
	}
	return nil
}

func Connect() (*pg.DB, error) {
	var (
		db  *pg.DB
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

func InsertMessage(db *pg.DB, m *Message) error {
	_, err := db.Exec(`
		INSERT INTO log_messages
		(time, user_id, number_of_posts, post_id) VALUES (?, ?, ?, ?)`,
		m.Time, m.UserId, m.NumberOfPosts, m.PostId)
	return err
}

func FetchMessages(db *pg.DB) []Message {
	rows, err := db.Query(`
		SELECT
		time, user_id, number_of_posts, post_id
		FROM log_messages
		ORDER BY record_date
		DESC LIMIT ?`,
		MESSAGES_SIZE)

	if err != nil {
		log.Fatal(err)
	}

	messages := make([]Message, 0)

	for rows.Next() {
		var m Message
		err = rows.Scan(&m.Time, &m.UserId, &m.NumberOfPosts, &m.PostId)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, m)
	}

	return messages
}
