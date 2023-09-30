package article

import (
	"encoding/json"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net-http/myapp/utils/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type articleResponse struct {
	Id             uint                  `json:"id"`
	Title          string                `json:"title"`
	ImageUrl       string                `json:"image"`
	Url            string                `json:"url"`
	PublishedAt    time.Time             `gorm:"not null" json:"date"`
	Tags           []*repository.Tag     `json:"tags"`
	UserLiked      []*repository.User    `json:"good"`
	UserBookMarked []*repository.User    `gojson:"book_marked"`
	Comments       []*repository.Comment `json:"comment"`
	Site           site                  `json:"site"`
	GoodNum        int                   `json:"good_num"`
	BookMarkedNum  int                   `json:"book_marked_num"`
	CommentNum     int                   `json:"comment_num"`
	NowBookmarked  bool                  `json:"now_bookmarked"`
	NowLiked       bool                  `json:"now_liked"`
}

type site struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

// GetArticleLatest は最新の記事を取得する関数です。
// この関数は HTTP GET リクエストのみを受け付けます。
// リクエストには "page" クエリパラメータが含まれている必要があります。
//
// レスポンスは記事のリストを JSON 形式で返します。
// 各記事オブジェクトはサイト情報、URL、イメージURL、タグ、タイトル、公開日時、
// コメント、ユーザーがブックマークしたかどうか、ユーザーがいいねしたかどうかを含みます。
//
// "page" クエリパラメータが存在しない、または数値に変換できない場合、ページ1がデフォルトとして使用されます。
//
// エラーの場合、該当するHTTPステータスコードとともにエラーメッセージを JSON 形式で返します。

func GetArticleLatest(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外受け付けない
	header := w.Header()
	header.Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}
	// クエリパラメータの取得
	page := r.URL.Query().Get("page")
	pageNation, err := strconv.Atoi(page)
	fmt.Println(pageNation)
	if err != nil {
		fmt.Println(err)
		pageNation = 1
	}

	// 記事の取得
	articles, err := repository.FindByArticles(pageNation)
	if err != nil {
		utils.WriteLogFile(err.Error())
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
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
	utils.WriteLogFile(string(userId))

	var responses []articleResponse
	for _, article := range *articles {
		res := createArticleResponse(article, userId)
		responses = append(responses, res)
	}
	json.NewEncoder(w).Encode(responses)
}
