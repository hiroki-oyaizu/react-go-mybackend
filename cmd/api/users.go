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

// 誕生日情報の構造体
type Birthday struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}
type Users struct {
	Id           int      `json:"id"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Age          int      `json:"age"`
	Mail         string   `json:"mail"`
	Password     string   `json:"password"`
	ProfileImage *string  `json:"profileImage"`
	Birthday     Birthday `json:"birthday"`
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
		var user Users
		var profileImageBytes []byte // この行を for ループ内に移動
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Mail, &user.Password, &user.Birthday.Year, &user.Birthday.Month, &user.Birthday.Day, &profileImageBytes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if profileImageBytes != nil {
			profileImage := string(profileImageBytes)
			user.ProfileImage = &profileImage
		}

		users = append(users, user)
	}

	// クエリ結果を閉じる
	err = rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// スライスをJSONに変換する
	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスヘッダーにContent-Typeを設定し、HTTPステータスコードを設定して、JSONデータを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *application) PostUser(w http.ResponseWriter, r *http.Request) {
	// Content-Type をチェック
	fmt.Println("Received Content-Type:", r.Header.Get("Content-Type"))
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
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
	stmt, err := db.Prepare("INSERT INTO users (firstName, lastName, age, mail, password, profileImage, year, month, day) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Age, user.Mail, user.Password, user.ProfileImage, user.Birthday.Year, user.Birthday.Month, user.Birthday.Day)

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

	err := db.QueryRow("SELECT id,firstName,lastName,age,profileImage FROM users WHERE id = ?", id).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age)

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
	stmt, err := db.Prepare("UPDATE users SET id=?, firstName=?, lastName=?, age=?,profileImage=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(u.Id, u.FirstName, u.LastName, u.Age, u.ProfileImage, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User updated successfully.")
}
