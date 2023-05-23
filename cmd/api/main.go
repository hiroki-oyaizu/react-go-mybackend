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

	err = database.CreateTableIfNotExists(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	err = database.CreateTweetTableIfNotExists(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
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

// 関数を作って
