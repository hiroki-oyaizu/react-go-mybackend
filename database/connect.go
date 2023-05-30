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
		firstName varchar(100),
		lastName varchar(100),
		age int,
		mail varchar(255) UNIQUE,
		password varchar(255),
		profileImage LONGTEXT,
		year int,
		month int,
		day int,
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
    tweet_content varchar(255),
    image text,
    user_id int,
    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
  );
`
	_, err := db.Exec(createTweetsTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateFollowTableIfNotExists(db *sql.DB) error {
	createFollowTableQuery := `
	 CREATE TABLE IF NOT EXISTS follows (
		 id int AUTO_INCREMENT,
		 user_id int,
		 follow_user_id int,
		 PRIMARY KEY(id),
		 FOREIGN KEY (user_id) REFERENCES users(id),
		 FOREIGN KEY (follow_user_id) REFERENCES users(id)
		 );
	`
	_, err := db.Exec(createFollowTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// いいねを誰がしたかを管理するLikesテーブルがなかったら作成する関数
func CreateLikesTableNotExists(db *sql.DB) error {
	createLikesTableQuery := `
	CREATE TABLE IF NOT EXISTS Likes (
		id int AUTO_INCREMENT PRIMARY KEY,
		user_id int,
		tweet_id int,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (tweet_id) REFERENCES tweets(id)
	);
`

	_, err := db.Exec(createLikesTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func CreateCommentsTableNotExists(db *sql.DB) error {
	createCommentsTableQuery := `
	CREATE TABLE IF NOT EXISTS comments (
		id int AUTO_INCREMENT PRIMARY KEY,
		user_id int NOT NULL,
		tweet_id int NOT NULL,
		comment varchar(255) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE
	);
		`
	_, err := db.Exec(createCommentsTableQuery)
	if err != nil {
		return err
	}
	return nil
}
