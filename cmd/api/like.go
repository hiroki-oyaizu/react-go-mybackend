package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"

	"github.com/go-chi/chi/v5"
)

type LikeCount struct {
	Id         int `json:"id"`
	UserId     int `json:"user_id"`
	LikeUserId int `json:"tweet_id"`
}

// いいねを誰がしたかを保存する関数
func (a *application) CreateLikeCount(w http.ResponseWriter, r *http.Request) {
	var likeCount LikeCount

	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&likeCount)
	if err != nil {
		fmt.Println(err)
	}
	db := database.Connect()
	if err != nil {
		http.Error(w, "Internal server error1", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Likes(user_id,tweet_id) VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error2", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(likeCount.UserId, likeCount.LikeUserId)
	if err != nil {
		http.Error(w, "Internal server error3", http.StatusInternalServerError)
		return

	}
	stmt.Close()
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Like created successfully"})
}

func (a *application) GetAllLikes(w http.ResponseWriter, r *http.Request) {

	db := database.Connect()
	defer db.Close()

	query := "SELECT tweet_id, COUNT(*) as like_count FROM Likes GROUP BY tweet_id"
	rows, err := db.Query(query)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	likeMap := make(map[string]int)
	for rows.Next() {
		var tweetID string
		var likeCount int
		if err := rows.Scan(&tweetID, &likeCount); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		likeMap[tweetID] = likeCount
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(likeMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a *application) GetLikeCount(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	db := database.Connect()
	id := chi.URLParam(r, "id")
	fmt.Println("id:", id)
	if id == "" {
		http.Error(w, "???", http.StatusBadRequest)
		return
	}

	// データベースに接続
	defer db.Close()

	// SQLクエリを準備し、指定されたtweetIdに一致するいいねの数を取得
	var count int
	query := "SELECT COUNT(*) FROM Likes WHERE tweet_id = ?"
	row := db.QueryRow(query, id)

	err := row.Scan(&count)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// カウントの値をキー'likeCount'と一緒にオブジェクトにラップ
	countObj := map[string]int{
		"likeCount": count,
	}

	// オブジェクトをJSONに変換
	res, err := json.Marshal(countObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 結果を返す
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
