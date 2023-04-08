package main

import (
	"encoding/json"
	"io"
	"net/http"
	"react-go-mybackend/database"
)

type Tweet struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Gender  string `json:"gender"`
	Toggle  bool   `json:"toggle"`
}

func (a *application) PostTweet(w http.ResponseWriter, r *http.Request) {
	var tweet Tweet
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := database.Connect()
	// データベースに接続

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// データベースにデータを挿入
	stmt, err := db.Prepare("INSERT INTO tweets (title, content, gender, toggle) VALUES (?, ?, ?, ?)")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tweet.Title, tweet.Content, tweet.Gender, tweet.Toggle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
