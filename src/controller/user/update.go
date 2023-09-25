package user

import (
	"encoding/json"
	jwt2 "github.com/dgrijalva/jwt-go"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/jwt"
	"net/http"
	"strconv"
	"strings"
)

type body struct {
	UserName         string `json:"userName" binding:"required"`
	SelfIntroduction string `json:"selfIntroduction" binding:"required"`
	GithubUser       string `json:"githubUser" binding:"required"`
	TwitterUser      string `json:"twitterUser" binding:"required"`
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// PUTメソッド以外受け付けない
	header := w.Header()
	header.Set("Content-Type", "application/json")
	if r.Method != "PUT" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}
	// headerからBearer Token読み出し
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// tokenの認証
	token, err := jwt.VerifyToken(tokenString)
	if err != nil {
		if err.Error() == "Token is expired" {
			utils.WriteLogFile("JWT Tokenが失効しています")
			utils.WriteLogFile(err.Error())
			w.WriteHeader(403)
			json.NewEncoder(w).Encode(utils.MyError{Message: "Token is expired"})
			return
		}
		utils.WriteLogFile("JWT Tokenを取得できませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// tokenからuser_idの取得
	claims := token.Claims.(jwt2.MapClaims)
	searchId := claims["user"].(string)
	userId, err := strconv.ParseUint(searchId, 10, 64)
	if err != nil {
		utils.WriteLogFile("interfaceをuintに変更できませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// userモデルの取得
	user := repository.User{}
	err = user.FindById(uint(userId))
	if err != nil {
		utils.WriteLogFile("userを取得することができませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// リクエストボディの取得
	var payload body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteLogFile("Request Bodyを取得できませんでした")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// ユーザー情報の更新処理を行います。
	// payloadから受け取った情報で、userの情報を更新します。
	user.UserName = payload.UserName                 // ユーザー名を更新
	user.TwitterUser = payload.TwitterUser           // Twitterのユーザー名を更新
	user.GithubUser = payload.GithubUser             // Githubのユーザー名を更新
	user.SelfIntroduction = payload.SelfIntroduction // 自己紹介文を更新

	// userの情報をデータベースに更新します。
	err = user.Update()

	// 更新処理が失敗した場合は、エラーログを出力し500エラー(Internal Server Error)を返します。
	if err != nil {
		utils.WriteLogFile("user情報を更新することができませんでした")                               // エラーログの出力
		utils.WriteLogFile(err.Error())                                            // エラーメッセージのログ出力
		w.WriteHeader(500)                                                         // 500エラーを返す
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"}) // エラーメッセージを返す
		return                                                                     // 処理を終了
	}
	json.NewEncoder(w).Encode(user)
}
