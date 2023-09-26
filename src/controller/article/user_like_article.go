package article

import (
	"encoding/json"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"net-http/myapp/controller"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/jwt"
	"net/http"
	"strconv"
	"strings"
)

type request struct {
	ArticleId uint `json:"article_id"`
}

func FindUserLikedArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response := controller.Response{Status: 401, Text: "入力内容をお確かめください"}
		w.WriteHeader(response.Status)
		json.NewEncoder(w).Encode(response)
		utils.WriteLogFile(err.Error())
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
	fmt.Println(userId)
	// いいねする処理を書く
	//記事の検索
	article, err := repository.FindLikeArticleById(body.ArticleId)
	if err != nil {
		utils.WriteLogFile("article not found")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(utils.MyError{Message: "article not found"})
		return
	}

	// Userモデルの取得
	user := repository.User{}
	err = user.FindById(uint(userId))
	if err != nil {
		utils.WriteLogFile("userを取得することができませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// 記事をいいねする
	article.UserLiked = append(article.UserLiked, &user)
	if err := article.Update(); err != nil {
		utils.WriteLogFile("記事を更新することができませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "記事を更新することができませんでした"})
		return
	}
	json.NewEncoder(w).Encode(article)
}
