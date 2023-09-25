package details

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ZennMainArticle struct {
	Article Article `json:"article"`
}

type Topics struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	PublishedAt time.Time `json:"published_at"`
	Path        string    `json:"path"`
	OgImageURL  string    `json:"og_image_url"`
	Topics      []Topics  `json:"topics"`
}

func (receiver *ZennMainArticle) GetDetailPage(url string) {
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
	if err := json.Unmarshal(body, receiver); err != nil {
		fmt.Println(err)
		return
	}
}
