package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"
)

type LikeCount struct {
	Id           int `json:"id"`
	UserId       int `json:"user_id"`
	LikeInUserId int `json:"like_in_user_id"`
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO likes(user_id,like_in_user_id) VALUES(?,?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	stmt.Close()

	_, err = stmt.Exec(likeCount.UserId, likeCount.LikeInUserId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Like created successfully"})
}

// func (a *application) GetLikeCount(w http.ResponseWriter, r *http.Request) {
// 	// var ?? ??
// 	// リクエストの内容を変数に格納する
// 	loggedInUserId, err := strconv.Atoi(r.URL.Query().Get("loggedInUserId"))
// 	if err != nil {
// 		http.Error(w, "Invalid user id", http.StatusBadRequest)
// 		return
// 	}
// 	// デコード
// 	likeInUserId, err := strconv.Atoi(r.URL.Query().Get("likeINuserId"))

// 	if err != nil {
// 		http.Error(w, "Invalid target user id", http.StatusBadRequest)
// 		return
// 	}
// 	db := database.Connect()
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	defer db.Close()

// 	// データベースを接続する
// 	// データベース接続を閉じる
// 	// クエリパラメーターを解析し変数に入れる
// 	// Sql文を書く
// 	// INSERT INTO Liks ()
// 	// exec
// 	//    http.status(Ok)
// }
