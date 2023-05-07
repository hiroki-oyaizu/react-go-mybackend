package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"react-go-mybackend/database"
)

type Tweets struct {
	ID            int     `json:"id"`
	Tweet_Content string  `json:"tweet_content"`
	Image         string  `json:"image"`
	UserID        int     `json:"user_id"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	ProfileImage  *string `json:"profileImage"`
}

func (a *application) AllGetTweet(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()
	var tweets []Tweets
	rows, err := db.Query("SELECT tweets.id, tweets.tweet_content, tweets.image, tweets.user_id, users.firstName, users.lastName, users.profileImage FROM tweets JOIN users ON tweets.user_id = users.id")

	if err != nil {
		log.Println("Error preparing statement:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var tweet Tweets
		err := rows.Scan(&tweet.ID, &tweet.Tweet_Content, &tweet.Image, &tweet.UserID, &tweet.FirstName, &tweet.LastName, &tweet.ProfileImage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	err = rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(tweets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *application) CrateTweet(w http.ResponseWriter, r *http.Request) {
	// Content-Type をチェック
	fmt.Println("Received Content-Type:", r.Header.Get("Content-Type"))
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var tweet Tweets
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := database.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tweets (tweet_content, image, user_id) VALUES(?, ?, ?)")

	if err != nil {
		log.Println("Error preparing statement:", err) // ログに詳細な情報を出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tweet.Tweet_Content, tweet.Image, tweet.UserID)
	if err != nil {
		log.Println("Error executing statement:", err) // ログに詳細な情報を出力
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
