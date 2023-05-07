package main

import (
	"net/http"
	"react-go-mybackend/database"
)

func (a *application) ListingCreate(w http.ResponseWriter, r *http.Request) {
	// データーベースに接続する
	db := database.Connect()
	// データーベース接続を閉じる
	defer db.Close()
	// レスポンスを受け取る
	// リクエストボディを読み込む
	// ニューデコードでデコードして実際にデコードする
}
