package main

import (
	"encoding/json"
	"io"
	"net/http"
	"react-go-mybackend/database"
)

type Follow struct {
	Id           int `json:"id"`
	UserId       int `json:"user_id"`
	FollowUserId int `json:"follow_user_id"`
}

func (a *application) CreateFollow(w http.ResponseWriter, r *http.Request) {
	// 変数を準備する
	var follow Follow
	//request bodyからデータを取得
	// デコードする
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//データベースを接続する
	db := database.Connect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//データベース接続を閉じる
	defer db.Close()
	// ステートメントを用意する
	stmt, err := db.Prepare("INSERT INTO follows (user_id, follow_user_id) VALUES(?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer stmt.Close()
	//データベースにデータを保存する
	_, err = stmt.Exec(follow.UserId, follow.FollowUserId)
	// クエリを実行する
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// 200のステータスコードを返す
	w.WriteHeader(http.StatusCreated)
}
