package article

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

type searchTag struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type searchArticleTagResponse struct {
	Data         []articleResponse `json:"data"`
	SearchTagUrl string            `json:"search_tag_url"`
	RelatedTags  []searchTag       `json:"related_tags"`
}

func SearchArticlesByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	// クエリパラメータの取得
	tag := r.URL.Query().Get("tag")
	tag = strings.ToLower(tag) // 小文字に変換

	// クエリパラメータの取得
	page := r.URL.Query().Get("page")
	pageNation, err := strconv.Atoi(page)
	if err != nil {
		pageNation = 1
	}

	// headerからBearer Token読み出し
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// tokenの認証
	token, err := jwt.VerifyToken(tokenString)
	var userId uint64
	userId = 0
	if err != nil {
		if err.Error() == "Token is expired" {
			utils.WriteLogFile("JWT Tokenが失効しています")
			utils.WriteLogFile(err.Error())
		} else {
			utils.WriteLogFile("JWT Tokenを取得できませんでした")
			utils.WriteLogFile(err.Error())
		}

	} else {
		claims := token.Claims.(jwt2.MapClaims)
		searchId := claims["user"].(string)
		userId, err = strconv.ParseUint(searchId, 10, 64)
		if err != nil {
			utils.WriteLogFile("interfaceをuintに変更できませんでした")
			utils.WriteLogFile(err.Error())
		}
	}

	// article_tagsテーブル(Many to Many)から検索
	tagRecord, err := repository.FindRelatedByTagNames(tag, pageNation)
	if err != nil {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error(DB Error)"})
		return
	}

	var responses []articleResponse
	var relatedTags []searchTag
	for _, article := range tagRecord.Articles {
		for _, articleTag := range article.Tags {
			if articleTag.ImageURL == "" || contains(relatedTags, articleTag.Name) {
				continue
			}
			searchTag := searchTag{Name: articleTag.Name, Image: articleTag.ImageURL}
			relatedTags = append(relatedTags, searchTag)
		}
		res := createArticleResponse(*article, userId)

		responses = append(responses, res)
	}

	result := searchArticleTagResponse{
		Data:         responses,
		SearchTagUrl: tagRecord.ImageURL,
		RelatedTags:  relatedTags,
	}

	json.NewEncoder(w).Encode(result)
}

func contains(searchTags []searchTag, str string) bool {
	for _, searchTag := range searchTags {
		if searchTag.Name == str {
			return true
		}
	}
	return false
}
