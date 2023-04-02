package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react-go-mybackend/database"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
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

func (app *application) GetDetailUser(w http.ResponseWriter, r *http.Request) {
	u := &Users{}
	db := database.Connect()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "???", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("SELECT id,firstNane,lastNane,age FROM users WHERE id = ?", id).Scan(&u.Id, &u.FirstNane, &u.LastNane, &u.Age)

	if err != nil {
		fmt.Fprintf(w, "Error querying database: %s", err.Error())
		return
	}

	usersJSON, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
	defer db.Close()
}

func (app *application) DeleteUserID(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing article ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Article deleted successfully.")
}

func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()
	u := &Users{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	stmt, err := db.Prepare("UPDATE users SET id=?, firstNane=?, lastNane=?, age=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(u.Id, u.FirstNane, u.LastNane, u.Age, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User updated successfully.")
}
