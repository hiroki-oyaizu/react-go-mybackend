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
	mux.Get("/article", app.GetArticle)
	mux.Post("/article", app.PostArticle)
	mux.Get("/article/{id}", app.GetArticleByID)
	mux.Put("/article/{id}", app.UpdateArticleByID)
	mux.Delete("/article/{id}", app.DeleteArticleByID)
	mux.Get("/users", app.AllGetUsers)
	mux.Post("/users", app.PostUser)
	mux.Get("/users/{id}", app.GetDetailUser)
	mux.Delete("/users/{id}", app.DeleteUserID)
	mux.Put("/users/{id}", app.UpdateUser)
	mux.Post("/tweets/new", app.PostTweet)
	mux.Post("/login", app.LoginUser)

	mux.Get("/tweet", app.AllGetTweet)
	mux.Post("/tweet/new", app.CrateTweet)
	mux.Post("/follow/new", app.CreateFollow)
	mux.Get("/follow/isFollowing", app.IsFollowing)
	mux.Post("/follow/unfollow", app.Unfolded)

	mux.Post("/like", app.CreateLikeCount)
	mux.Get("/like/all", app.GetAllLikes)
	mux.Get("/like/{id}", app.GetLikeCount)

	mux.Get("/comment/{id}", app.GetComments)
	mux.Post("/comment", app.CreateComment)
	return mux
}
