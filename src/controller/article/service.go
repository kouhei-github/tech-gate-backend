package article

import (
	"net-http/myapp/repository"
	"strings"
)

func createArticleResponse(article repository.Article, userId uint64) articleResponse {
	var res articleResponse
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
	if article.Comments == nil {
		res.CommentNum = 0
	} else {
		res.CommentNum = len(article.Comments)
	}
	if article.UserLiked == nil {
		res.GoodNum = 0
	} else {
		res.GoodNum = len(article.UserLiked)
	}
	if userId != 0 {
		if includeLoginUser(article.UserBookMarked, uint(userId)) {
			res.NowBookmarked = true
		} else {
			res.NowBookmarked = false
		}

		if includeLoginUser(article.UserLiked, uint(userId)) {
			res.NowLiked = true
		} else {
			res.NowLiked = false
		}
	}

	return res
}

func includeLoginUser(users []*repository.User, userId uint) bool {
	for _, user := range users {
		if user.ID == userId {
			return true
		}
	}
	return false
}
