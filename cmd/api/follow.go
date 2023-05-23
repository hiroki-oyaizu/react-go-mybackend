package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"
	"strconv"
)

type Date struct {
	UserName  string  `json:"username"`
	Email     string  `json:"email"`
	Is_active bool    `json:"is_active"`
	Profile   Profile `json:"profile"`
}

type Profile struct {
	First_name string         `json:"first_name"`
	Last_name  string         `json:"last_name"`
	Age        int            `json:"age"`
	Location   Location       `json:"location"`
	Social     []Social_links `json:"social_links"`
}

type Location struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

type Social_links struct {
	Platform string `json:"platform"`
	Handle   string `json:"handle"`
}

type Follow struct {
	Id           int `json:"id"`
	UserId       int `json:"user_id"`
	FollowUserId int `json:"follow_user_id"`
}

func (a *application) CreateFollow(w http.ResponseWriter, r *http.Request) {
	// 変数を準備する
	var follow Follow

	// Request bodyからデータを取得
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&follow)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	fmt.Println(follow.UserId)
	fmt.Println(follow.FollowUserId)
	// Validate the request body
	if follow.UserId == 0 || follow.FollowUserId == 0 {
		http.Error(w, "Invalid user id or follow user id", http.StatusBadRequest)
		return
	}

	//データベースを接続する
	db := database.Connect()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//データベース接続を閉じる
	defer db.Close()

	// ステートメントを用意する
	stmt, err := db.Prepare("INSERT INTO follows (user_id, follow_user_id) VALUES(?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	//データベースにデータを保存する
	_, err = stmt.Exec(follow.UserId, follow.FollowUserId)
	// クエリを実行する
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Follow created successfully"})
}

// w.Header().Set("Content-Type", "application/json")
// この行は、レスポンスのヘッダーに"Content-Type"という項目を設定しています。
// "Content-Type"は、レスポンスの内容（ボディ）がどのような形式のデータであるかを指定するものです。
// ここでは"application/json"と指定しているので、「このレスポンスの内容はJSON形式のデータですよ」ということを
// クライアント（このレスポンスを受け取る側）に伝えています。

// json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Follow created successfully"})
// この行は、レスポンスの内容（ボディ）を作成しています。
// 具体的には、"status"と"message"という二つの項目を持つマップ（キーと値のペアの集まり）を作成し、
// それをJSON形式の文字列に変換（エンコード）して、レスポンスの内容としています。
// "status"は"success"という値を持ち、"message"は"Follow created successfully"という値を持っています。
// これにより、「処理が成功しました」という情報をクライアントに伝えることができます。
// json.NewEncoder(w)は、レスポンス（w）にJSON形式のデータを書き込むための準備を行っています。
// Encode関数は、引数として与えられたデータ（ここではマップ）をJSON形式に変換し、それをレスポンスに書き込みます。

// IsFollowing関数は、あるユーザーが別のユーザーをフォローしているかどうかをチェックします。
func (a *application) IsFollowing(w http.ResponseWriter, r *http.Request) {
	// ログインユーザーのIDと対象ユーザーのIDはクエリパラメータとして提供されるべきです。
	// これらのIDを取得します。
	loggedInUserId, err := strconv.Atoi(r.URL.Query().Get("loggedInUserId"))
	// もしIDが無効（文字列を数値に変換できない）ならエラーメッセージを返します。
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	targetUserId, err := strconv.Atoi(r.URL.Query().Get("targetUserId"))
	// もし対象ユーザーIDが無効ならエラーメッセージを返します。
	if err != nil {
		http.Error(w, "Invalid target user id", http.StatusBadRequest)
		return
	}
	// データベースに接続します。
	db := database.Connect()
	// もし接続に問題があったらエラーメッセージを返します。
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// データベース接続を閉じるためのコードです。これは最後に必ず実行されます。
	defer db.Close()

	// データベースに問い合わせるためのSQLクエリを準備します。
	query := "SELECT COUNT(*) FROM follows WHERE user_id = ? AND follow_user_id = ?"

	// クエリを実行します。ログインユーザーIDと対象ユーザーIDをパラメータとして渡します。
	row := db.QueryRow(query, loggedInUserId, targetUserId)

	var count int
	// クエリの結果（フォローの数）をcount変数に格納します。
	err = row.Scan(&count)
	// もし結果の取得に問題があったらエラーメッセージを返します。
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// もしcountが0より大きければ、ログインユーザーは対象ユーザーをフォローしています。
	isFollowing := count > 0

	// レスポンスを送信します。これにはフォローしているかどうかの情報が含まれます。
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isFollowing": isFollowing})
}

func (a *application) Unfolded(w http.ResponseWriter, r *http.Request) {
	// Request bodyからデータを取得
	var follow Follow
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the request body
	if follow.UserId == 0 || follow.FollowUserId == 0 {
		http.Error(w, "Invalid user id or follow user id", http.StatusBadRequest)
		return
	}

	//データベースを接続する
	db := database.Connect()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//データベース接続を閉じる
	defer db.Close()

	// ステートメントを用意する
	stmt, err := db.Prepare("DELETE FROM follows WHERE user_id = ? AND follow_user_id = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	//データベースからデータを削除する
	_, err = stmt.Exec(follow.UserId, follow.FollowUserId)
	// クエリを実行する
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Unfollowed successfully"})
}
