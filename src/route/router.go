package route

import (
	"net-http/myapp/controller"
	"net-http/myapp/controller/article"
	"net/http"
)

type Router struct {
	Mutex *http.ServeMux
}

func (router *Router) GetRouter() {

	// 練習
	router.Mutex.HandleFunc("/", controller.Handler)
	router.Mutex.HandleFunc("/two", controller.HandlerTwo)

	// 最新の記事一覧
	router.Mutex.HandleFunc("/api/v1/article/latest", article.GetArticleLatest)

	//// おすすめの記事一覧
	//router.Mutex.HandleFunc("/api/v1/article/recommend", controller.HandlerTwo)

	// タグで絞り込み
	router.Mutex.HandleFunc("/api/v1/article", article.SearchArticlesByTag)

	// 記事をDBに保存

	// 人気の記事一覧

	// タグの管理
}
