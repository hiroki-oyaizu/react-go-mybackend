package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"
)

type Users struct {
	Id        int    `json:"id"`
	FirstNane string `json:"firstNane"`
	LastNane  string `json:"lastNane"`
	Age       int    `json:"age"`
}

func (a application) AllGetUsers(w http.ResponseWriter, r *http.Request) {
	// データベースに接続する
	db := database.Connect()
	// データベース接続を閉じる
	defer db.Close()

	// ユーザー情報を格納するためのスライスを宣言する
	var users []Users

	// データベースからユーザー情報を取得するためのクエリを実行する
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("error", err)
	}

	// ループを使用して、すべてのユーザー情報をスキャンし、スライスに追加する
	for rows.Next() {
		// 新しいユーザー情報を格納するための変数を宣言する
		var user Users
		// データベースからユーザー情報をスキャンする
		// はい、その通りです。rows.Scan()関数は、データベースから取得した行の各カラムの値を、指定した変数に代入します。この場合、rows.Scan()はuser.Id, user.FirstNane, user.LastNane, user.Ageに、データベースから取得した対応するカラムの値を代入しています。
		err := rows.Scan(&user.Id, &user.FirstNane, &user.LastNane, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// スライスにユーザー情報を追加する
		users = append(users, user)
	}
	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *application) PostUser(w http.ResponseWriter, r *http.Request) {
	var user Users
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&user)
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
	stmt, err := db.Prepare("INSERT INTO users (firstNane, lastNane, age) VALUES (?, ?, ?)")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstNane, user.LastNane, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
