package article

import (
	"encoding/json"
	"fmt"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type response struct {
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

	var responses []response
	for _, article := range *articles {
		var res response
		var redirectSite site
		if strings.Contains(article.Url, "qiita") {
			redirectSite.Image = "https://youliangdao.s3.ap-northeast-1.amazonaws.com/favicon.png"
			redirectSite.Name = "qiita.com"
		} else {
			redirectSite.Image = "https://youliangdao.s3.ap-northeast-1.amazonaws.com/logo-only.png"
			redirectSite.Name = "zenn.dev"
		}
		res.Id = article.ID
		res.Site = redirectSite
		res.Url = article.Url
		res.ImageUrl = article.ImageUrl
		res.Tags = article.Tags
		res.Title = article.Title
		res.PublishedAt = article.PublishedAt
		res.Comments = article.Comments
		res.UserBookMarked = article.UserBookMarked
		res.UserLiked = article.UserLiked
		responses = append(responses, res)
	}
	json.NewEncoder(w).Encode(responses)
}
