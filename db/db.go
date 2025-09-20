package db

import (
	"database/sql"
	"log"

	"brange.net/quicknote/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Some function for connecting and sending requests to database server.
// Edit .env to configure !
func InitDB(cfg config.Config) {
	var err error
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ")/" + cfg.DBName
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func SaveMessage(content, link string) error {
	_, err := DB.Exec("INSERT INTO messages (content, link) VALUES (?, ?)", content, link)
	return err
}

func GetMessage(link string) (string, error) {
	var content string
	err := DB.QueryRow("SELECT content FROM messages WHERE link = ?", link).Scan(&content)
	return content, err
}

func DeleteMessage(link string) error {
	_, err := DB.Exec("DELETE FROM messages WHERE link = ?", link)
	return err
}
