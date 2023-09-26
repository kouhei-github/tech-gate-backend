package route

import (
	"net-http/myapp/controller"
	"net-http/myapp/controller/article"
	"net-http/myapp/controller/google"
	"net-http/myapp/controller/user"
)

func (router *Router) GetAuthRouter() {
	router.Mutex.HandleFunc("/api/v1/auth", controller.HandlerTwo)
	router.Mutex.HandleFunc("/api/v1/google/login", google.GoogleOauth)
	router.Mutex.HandleFunc("/api/v1/google/callback", google.GoogleLoginHandler)
	// ユーザー情報の更新
	router.Mutex.HandleFunc("/api/v1/user/update", user.UserUpdateHandler)

	// いいね　アクション
	router.Mutex.HandleFunc("/api/v1/article/like", article.SaveUserLikedArticle)
	// ユーザーのいいね一覧
	router.Mutex.HandleFunc("/api/v1/article-like", article.FindUserLikedArticle)

	// ブックマーク アクション
	router.Mutex.HandleFunc("/api/v1/article/book-mark", article.SaveUserBookMarkArticle)
	// ユーザーのブックマーク一覧
	router.Mutex.HandleFunc("/api/v1/book-mark", article.FindUserBookMarkedArticle)

	// コメント アクション

	// ユーザーのコメント一覧
}
