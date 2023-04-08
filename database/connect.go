package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {

	user := "webuser"
	password := "webpass"
	host := "localhost"
	port := "3307"
	database_name := "go_mysql8_development"

	dbconf := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database_name + "?charset=utf8mb4"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func FetchRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("Err2")
		panic(err.Error())
	}
	return rows
}

// CreateTableIfNotExists関数を追加
func CreateTableIfNotExists(db *sql.DB) error {
	// テーブルが存在しない場合に作成するクエリ
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id int AUTO_INCREMENT,
			firstNane varchar(100),
			lastNane varchar(100),
			age int,
			mail varchar(255) UNIQUE,
			password varchar(255),
			PRIMARY KEY(id)
		);
	`

	// クエリを実行
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func CreateTweetTableIfNotExists(db *sql.DB) error {
	createTweetsTableQuery := `
	  CREATE TABLE IF NOT EXISTS tweets (
		id int AUTO_INCREMENT,
		title varchar(255),
		content varchar(255),
		gender ENUM('male', 'female'),
		toggle BOOLEAN,
		PRIMARY KEY(id)
	  );
	`

	_, err := db.Exec(createTweetsTableQuery)
	if err != nil {
		return err
	}

	createTweetDaysTableQuery := `
	  CREATE TABLE IF NOT EXISTS tweet_days (
		id int AUTO_INCREMENT,
		tweet_id int,
		day ENUM('sun', 'mon', 'tue', 'wed', 'thu', 'fri', 'sat'),
		PRIMARY KEY(id),
		FOREIGN KEY (tweet_id) REFERENCES tweets(id)
	  );
	`

	_, err = db.Exec(createTweetDaysTableQuery)
	if err != nil {
		return err
	}

	return nil
}
