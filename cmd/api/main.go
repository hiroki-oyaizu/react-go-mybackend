package main

import (
	"fmt"

	"log"
	"net/http"

	"react-go-mybackend/database"
)

func init() {
	db := database.Connect()
	defer db.Close()

	err := db.Ping()

	if err != nil {
		fmt.Println("データベース接続失敗")
		return
	} else {
		fmt.Println("データベース接続成功")
	}
}

type application struct {
	Domain string
}

const port = 8080

func main() {
	var app application
	app.Domain = "example.com"

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

// func getRows(db *sql.DB) {
// 	rows, err := db.Query("SELECT * FROM users")
// 	if err != nil {
// 		log.Fatalf("getRows db.Query error err:%v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		u := &database.User{}
// 		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Age); err != nil {
// 			log.Fatalf("getRows rows.Scan error err:%v", err)
// 		}
// 		fmt.Println(u)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatalf("getRows rows.Err error err:%v", err)
// 	}
// }
