package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"react-go-mybackend/database"
	"react-go-mybackend/internal/models"

	// "log"

	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status string `json:"status"`
	}{
		Status: "actice",
	}
	// json.Marshalは構造体をjsonに変換します。
	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

}

func (app *application) New(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "new page %s", app.Domain)
}

type Jtest struct {
	Status  int
	Message string
}

func (a *application) checkJson(w http.ResponseWriter, r *http.Request) {
	test := Jtest{Status: http.StatusOK, Message: "確認OK"}

	res, err := json.Marshal(test)

	if err != nil {
		log.Fatalln(err)
		return
	}
	//ファイルの種類を表している
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app *application) AllTweets(w http.ResponseWriter, r *http.Request) {
	var tweets []models.Tweet

	tweet_data1 := models.Tweet{
		ID:    1,
		Title: "表示確認",
		Tweet: "ツイート",
	}
	// 　　　　　　　　　どこに（追加元）何を
	tweets = append(tweets, tweet_data1)

	tweet_data2 := models.Tweet{
		ID:    2,
		Title: "表示確認2",
		Tweet: "ツイート2",
	}

	tweets = append(tweets, tweet_data2)

	json_tweets, err := json.Marshal(tweets)

	if err != nil {
		log.Fatalln(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_tweets)
}

type User struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
	AGE  int    `json:"age"`
}

func (a *application) AllUser(w http.ResponseWriter, r *http.Request) {
	// どこのデーターベースか示す情報
	db := database.Connect()
	// 最終的に閉じる
	defer db.Close()
	// MySQLからusersテーブルの情報を取ってくる
	rows := database.FetchRows(db)
	// 1人のユーザー情報を格納する空のオブジェクトを用意
	only_user := User{}
	// 全部のユーザーの情報を格納するスライスを定義
	var users []User
	for rows.Next() {
		error := rows.Scan(&only_user.ID, &only_user.NAME, &only_user.AGE)
		if error != nil {
			fmt.Println("scan error")
		} else {
			users = append(users, only_user)
		}
	}

	res, err := json.Marshal(users)
	if err != nil {
		log.Fatalln(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type Article struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func (a *application) GetArticle(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		log.Fatalf("getRows db.Query error err:%v", err)
	}
	only_ar := Article{}
	var article []Article
	for rows.Next() {
		err := rows.Scan(&only_ar.Id, &only_ar.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		article = append(article, only_ar)
	}
	res, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *application) PostArticle(w http.ResponseWriter, r *http.Request) {
	// リクエストボディを解析
	var article Article
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&article)
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
	stmt, err := db.Prepare("INSERT INTO articles (title) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	a := &Article{}
	db := database.Connect()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing article ID", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("SELECT id, title FROM articles WHERE id = ?", id).Scan(&a.Id, &a.Title)
	if err != nil {
		fmt.Fprintf(w, "Error querying database: %s", err.Error())
		return
	}

	articleJSON, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(articleJSON)
	defer db.Close()
}

func (app *application) UpdateArticleByID(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()

	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	stmt, err := db.Prepare("UPDATE articles SET id=?, title=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(article.Id, article.Title, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User updated successfully.")
}

// func (app *application) DeleteArticleByID(w http.ResponseWriter, r *http.Request) {
// 	db := database.Connect()
// 	defer db.Close()

// 	var article Article
// 	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	id := chi.URLParam(r, "id")
// 	stmt, err := db.Prepare("UPDATE articles SET id=?, title=? WHERE id=?")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	if _, err := stmt.Exec(article.Id, article.Title, id); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, "User updated successfully.")
// }

func (app *application) DeleteArticleByID(w http.ResponseWriter, r *http.Request) {
	db := database.Connect()
	defer db.Close()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing article ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM articles WHERE id=?")
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
