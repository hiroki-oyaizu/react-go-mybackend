package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// 新規ルータ作成
	mux := chi.NewRouter()
	// サーバは落とさずにエラーレスポンスを返せるようにリカバリーするmiddlewareです。。ログを記録する
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Get("/", app.Home)
	mux.Get("/new", app.New)
	mux.Get("/hello", app.checkJson)
	mux.Get("/tweets", app.AllTweets)
	mux.Get("/users", app.AllUser)
	mux.Get("/article", app.GetArticle)
	mux.Post("/article", app.PostArticle)
	mux.Get("/article/{id}", app.GetArticleByID)
	mux.Put("/article/{id}", app.UpdateArticleByID)
	mux.Delete("/article/{id}", app.DeleteArticleByID)
	return mux
}
