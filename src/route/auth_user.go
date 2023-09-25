package route

import (
	"net-http/myapp/controller"
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

	// コメント アクション

	// ブックマーク アクション

	// ユーザーのいいね一覧

	// ユーザーのブックマーク一覧

	// ユーザーのコメント一覧
}
