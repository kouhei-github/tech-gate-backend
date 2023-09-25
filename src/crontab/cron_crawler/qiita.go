package cron_crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net-http/myapp/crontab/cron_crawler/details"
	"net-http/myapp/repository"
	"net/http"
	"time"
)

type Qiita struct {
	Url string
}

type QiitaArticle struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Tags      []Tag     `json:"tags"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}
type Tag struct {
	Name string `json:"name"`
}

func (e Qiita) Run() {

	url := e.Url + "?page=1&per_page=50"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var qiitaArticle []QiitaArticle
	if err := json.Unmarshal(body, &qiitaArticle); err != nil {
		fmt.Println(err)
		return
	}

	var articles []repository.Article
	for _, article := range qiitaArticle {

		tags := make([]*repository.Tag, len(article.Tags))
		for i, tag := range article.Tags {
			finds, err := repository.FindByTagNames(tag.Name)
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(finds) != 0 {
				tags[i] = &finds[0]
			} else {
				tag := &repository.Tag{Name: tag.Name}
				if err := tag.Save(); err != nil {
					fmt.Println(err)
					return
				}
				tags[i] = tag
			}
		}

		// OGP Imageを取得する[Qiitaの記事タイトルの画像]
		titleImage, err := details.GetQiitaDetail(article.URL)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 記事の作成
		article := repository.Article{
			Title:       article.Title,
			ZennSlug:    article.ID,
			PublishedAt: article.CreatedAt,
			Url:         article.URL,
			ImageUrl:    titleImage,
			Tags:        tags,
		}

		articles = append(articles, article)
	}

	// 記事の保存
	for _, article := range articles {
		finds, err := repository.FindByTitles(article.Title)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(finds) != 0 {
			continue
		}
		if err := article.Save(); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Done")
}
