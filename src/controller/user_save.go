package controller

import (
	"encoding/json"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/jwt"
	"net/http"
	"strings"
)

type body struct {
	Name     string
	Email    string `json:"email" binding:"required"`
	Password string
}

func UserSaveHandler(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外受け付けない
	header := w.Header()
	header.Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	// headerからBearer Token読み出し
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.VerifyToken(tokenString)
	if err != nil {
		utils.WriteLogFile("Request Bodyを取得できませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}
	claims := token.Claims.(jwt2.MapClaims)
	fmt.Println(claims["user"])
	// headerからBearer Token読み出し

	// リクエストボディの取得
	var payload body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteLogFile("Request Bodyを取得できませんでした")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	u := &repository.User{Email: payload.Email}
	err = u.Save()
	if err != nil {
		utils.WriteLogFile("ユーザーの保存に失敗しました")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}
	//data := []interface{}{
	//	"Authテスト", "認証", "認可", 1997, true,
	//}

	json.NewEncoder(w).Encode(u)
}
