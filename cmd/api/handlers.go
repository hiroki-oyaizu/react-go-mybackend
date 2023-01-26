package main

import (
	"encoding/json"
	"fmt"
	"log"
	"react-go-mybackend/database"
	"react-go-mybackend/internal/models"

	// "log"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status string `json:"status"`
	}{
		Status: "actice",
	}
	// json.Marshalは構造体をjsonに変換します。
	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

}

func (app *application) New(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "new page %s", app.Domain)
}

type Jtest struct {
	Status  int
	Message string
}

func (a *application) checkJson(w http.ResponseWriter, r *http.Request) {
	test := Jtest{Status: http.StatusOK, Message: "確認OK"}

	res, err := json.Marshal(test)

	if err != nil {
		log.Fatalln(err)
		return
	}
	//ファイルの種類を表している
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app *application) AllTweets(w http.ResponseWriter, r *http.Request) {
	var tweets []models.Tweet

	tweet_data1 := models.Tweet{
		ID:    1,
		Title: "表示確認",
		Tweet: "ツイート",
	}
	// 　　　　　　　　　どこに（追加元）何を
	tweets = append(tweets, tweet_data1)

	tweet_data2 := models.Tweet{
		ID:    2,
		Title: "表示確認2",
		Tweet: "ツイート2",
	}

	tweets = append(tweets, tweet_data2)

	json_tweets, err := json.Marshal(tweets)

	if err != nil {
		log.Fatalln(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_tweets)
}

type User struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
	AGE  int    `json:"age"`
}

func (a *application) AllUser(w http.ResponseWriter, r *http.Request) {
	// どこのデーターベースか示す情報
	db := database.Connect()
	// 最終的に閉じる
	defer db.Close()
	// MySQLからusersテーブルの情報を取ってくる
	rows := database.FetchRows(db)
	// 1人のユーザー情報を格納する空のオブジェクトを用意
	only_user := User{}
	// 全部のユーザーの情報を格納するスライスを定義
	var users []User
	for rows.Next() {
		error := rows.Scan(&only_user.ID, &only_user.NAME, &only_user.AGE)
		if error != nil {
			fmt.Println("scan error")
		} else {
			users = append(users, only_user)
		}
	}

	res, err := json.Marshal(users)
	if err != nil {
		log.Fatalln(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
