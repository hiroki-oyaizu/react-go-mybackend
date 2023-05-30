package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"

	"github.com/go-chi/chi/v5"
)

type Comment struct {
	UserId  int    `json:"user_id"`
	TweetId int    `json:"tweet_id"`
	Comment string `json:"comment"`
}

func (a *application) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&comment)
	fmt.Println(comment.Comment)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	db := database.Connect()
	if comment.UserId == 0 || comment.TweetId == 0 || comment.Comment == "" {
		http.Error(w, "Invalid user id or tweet id or comment", http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare("INSERT INTO comments (user_id, tweet_id, comment) VALUES(?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error222", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(comment.UserId, comment.TweetId, comment.Comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	defer db.Close()
	w.WriteHeader(http.StatusOK)
}

func (a *application) GetComments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetComments")
	comment := Comment{}
	// データーベースを接続する
	db := database.Connect()
	// データーベースを閉じる

	id := chi.URLParam(r, "id")
	fmt.Println(id)

	if id == "" {
		http.Error(w, "idが空です。", http.StatusBadRequest)
		return
	}

	// SQLを書いてステートメントを用意する
	err := db.QueryRow("SELECT comment FROM comments WHERE tweet_id = ?", id).Scan(&comment.Comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// データーベースからデータを取得する
	fmt.Println(comment)
	commentJson, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(commentJson)
	// データをJSONに変換する
	// ステータスOKとデータを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(commentJson)

	defer db.Close()
}
