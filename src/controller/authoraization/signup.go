package authoraization

import (
	"encoding/json"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/jwt"
	"net-http/myapp/utils/password"
	"net/http"
	"strconv"
)

type body struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	var payload body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteLogFile("Request Bodyを取得できませんでした")
		w.WriteHeader(500)
		w.Header().Set("its", "error")
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}
	hashPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		utils.WriteLogFile("Passwordのハッシュ化に失敗しました")
		w.WriteHeader(500)
		w.Header().Set("its", "error")
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}
	// ユーザーが存在するか確認
	user := &repository.User{
		Email:    payload.Email,
		UserName: payload.UserName,
		Password: hashPassword,
		IsLogin:  true,
	}
	users, err := user.FindByEmail(payload.Email)
	if err != nil {
		utils.WriteLogFile("ユーザーの検索で失敗しました")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	if len(users) != 0 {
		utils.WriteLogFile("既にユーザー名は存在します")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Emailアドレスは既に存在します"})
		return
	}

	if err := user.Save(); err != nil {
		utils.WriteLogFile("ユーザー情報をDBに保存できませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// jwt-tokenを返却 && ログインの期限を設定
	token, err := jwt.CreateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		utils.WriteLogFile("JWT Tokenが発行できませんでした")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(utils.MyError{Message: "tokenが向こうです"})
		return
	}

	response := struct {
		JwtToken         string `json:"jwtToken"`
		UserName         string `json:"userName"`
		UserImage        string `json:"userImage"`
		GithubUser       string `json:"githubUser"`
		TwitterUser      string `json:"twitterUser"`
		SelfIntroduction string `json:"selfIntroduction"`
		Email            string `json:"email"`
	}{JwtToken: token, UserName: user.UserName, UserImage: user.Image, GithubUser: user.GithubUser, TwitterUser: user.TwitterUser, SelfIntroduction: user.SelfIntroduction, Email: user.Email}
	utils.WriteLogFile(token)
	json.NewEncoder(w).Encode(response)
}
