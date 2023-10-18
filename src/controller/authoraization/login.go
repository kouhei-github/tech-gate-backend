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

type loginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
}

func NormalLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	var payload loginBody
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteLogFile("Request Bodyを取得できませんでした")
		w.WriteHeader(500)
		w.Header().Set("its", "error")
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	user := &repository.User{}
	users, err := user.FindByEmail(payload.Email)
	if err != nil {
		utils.WriteLogFile("ユーザーの検索で失敗しました")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	if len(users) == 0 {
		utils.WriteLogFile("存在しないメールアドレスです")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.MyError{Message: "存在しないメールアドレスです"})
		return
	}

	user = &users[0]

	// passwordが正しいか確認
	ok := password.VerifyPassword(payload.Password, user.Password)
	if !ok {
		utils.WriteLogFile("パスワードが正しくないです")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(utils.MyError{Message: "パスワードが正しくないです"})
		return
	}

	user.IsLogin = true // ログイン状態に変更
	if err := user.Update(); err != nil {
		utils.WriteLogFile("ユーザー情報を更新できませんでした")
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
