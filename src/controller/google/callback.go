package google

import (
	"context"
	"encoding/json"
	v2 "google.golang.org/api/oauth2/v2"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/google"
	"net-http/myapp/utils/jwt"
	"net/http"
	"strconv"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外受け付けない
	header := w.Header()
	header.Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	code := r.URL.Query().Get("code")

	// ユーザーの情報の取得
	config := google.GetConnect()
	ctx := context.Background()
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	if tok.Valid() == false {
		utils.WriteLogFile("vaild token")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	service, _ := v2.New(config.Client(ctx, tok))
	tokenInfo, _ := service.Tokeninfo().AccessToken(tok.AccessToken).Context(ctx).Do()

	// ユーザーが存在するか確認
	user := &repository.User{}
	users, err := user.FindByEmail(tokenInfo.Email)
	if err != nil {
		utils.WriteLogFile("ユーザーの検索で失敗しました")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// ユーザーの詳細情報をAPIで取得
	userInfo, err := service.Userinfo.Get().Context(ctx).Do()
	if err != nil {
		utils.WriteLogFile("ユーザー情報の取得で失敗しました")
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	// ユーザーの情報を保存 or 更新
	if len(users) == 0 {
		// ユーザーが存在しないなら保存
		user.UserName = userInfo.Name
		user.Email = userInfo.Email
		user.Image = userInfo.Picture
		user.IsLogin = true // ログイン状態に変更
		if err := user.Save(); err != nil {
			utils.WriteLogFile("ユーザー情報をDBに保存できませんでした")
			utils.WriteLogFile(err.Error())
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
			return
		}
	} else {
		// ユーザーが存在するなら更新
		user = &users[0]
		user.IsLogin = true // ログイン状態に変更
		if err := user.Update(); err != nil {
			utils.WriteLogFile("ユーザー情報を更新できませんでした")
			utils.WriteLogFile(err.Error())
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
			return
		}
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
