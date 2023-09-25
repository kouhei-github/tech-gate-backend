package controller

import (
	"encoding/json"
	"fmt"
	"net-http/myapp/repository"
	"net/http"
)

// Handler HandlerはHTTPリクエストを処理します。
//
// HTTPリクエストのヘッダー情報を取得し、コンソールに出力します。
// その後、"TEST"の文字列も出力します。
// 最後に、空のUserオブジェクトをJSON形式でレスポンスとして戻します。
//
// w: HTTPレスポンスを書き込むためのResponseWriter
// r: 受け取ったHTTPリクエスト
func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.Header
	fmt.Println(query)
	fmt.Println("TEST")
	u := repository.User{}
	//data := []interface{}{
	//	"Authテスト", "認証", "認可", 1997, true,
	//}

	json.NewEncoder(w).Encode(u)
}

func HandlerTwo(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	var data []interface{}
	data = []interface{}{
		"Authテスト", "認証", "認可", 1997, true,
	}

	json.NewEncoder(w).Encode(data)
}
