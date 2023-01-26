package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Age       string
}

type Tweet struct {
	ID int
}

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

// dbとはどのポートのSQLなのかの情報が入っている
func GetRows(db *sql.DB) {
	// Queryメソッド複数レコードを取得したいときに活用できます
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatalf("getRows db.Query error err:%v", err)
	}
	// 最終的に閉じる
	defer rows.Close()
	// rows.Next()各レコードに対して操作できる
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Age); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		fmt.Println(u)
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("getRows rows.Err error err:%v", err)
	}
}

func GetSingleRow(db *sql.DB, userID int) {
	u := &User{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", userID).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Age)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("getSingleRow no records.")
		return
	}
	if err != nil {
		log.Fatalf("getSingleRow db.QueryRow error err:%v", err)
	}
	fmt.Println(u)
}

func FetchRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("Err2")
		panic(err.Error())
	}
	return rows
}
