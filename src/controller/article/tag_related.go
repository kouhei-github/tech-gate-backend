package article

import (
	"encoding/json"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net/http"
	"strconv"
	"strings"
)

type searchTag struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type searchArticleTagResponse struct {
	Data         []response  `json:"data"`
	SearchTagUrl string      `json:"search_tag_url"`
	RelatedTags  []searchTag `json:"related_tags"`
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

	// article_tagsテーブル(Many to Many)から検索
	tagRecord, err := repository.FindRelatedByTagNames(tag, pageNation)
	if err != nil {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error(DB Error)"})
		return
	}

	var responses []response
	var relatedTags []searchTag
	for _, article := range tagRecord.Articles {
		var res response
		var redirectSite site
		if strings.Contains(article.Url, "qiita") {
			redirectSite.Image = "https://youliangdao.s3.ap-northeast-1.amazonaws.com/favicon.png"
			redirectSite.Name = "qiita.com"
		} else {
			redirectSite.Image = "https://youliangdao.s3.ap-northeast-1.amazonaws.com/logo-only.png"
			redirectSite.Name = "zenn.dev"
		}

		res.Site = redirectSite
		res.Url = article.Url
		res.ImageUrl = article.ImageUrl
		res.Tags = article.Tags
		res.Title = article.Title
		res.PublishedAt = article.PublishedAt
		res.Comments = article.Comments
		res.UserBookMarked = article.UserBookMarked
		res.UserLiked = article.UserLiked
		res.Tags = article.Tags

		for _, articleTag := range article.Tags {
			if articleTag.ImageURL == "" {
				continue
			}
			searchTag := searchTag{Name: articleTag.Name, Image: articleTag.ImageURL}
			relatedTags = append(relatedTags, searchTag)
		}

		responses = append(responses, res)
	}

	result := searchArticleTagResponse{
		Data:         responses,
		SearchTagUrl: tagRecord.ImageURL,
		RelatedTags:  relatedTags,
	}

	json.NewEncoder(w).Encode(result)
}
