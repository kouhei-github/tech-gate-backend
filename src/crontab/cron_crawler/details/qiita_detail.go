package details

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetQiitaDetail(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	titleImage := getOgpImage(string(body))
	return titleImage, nil
}

func getOgpImage(html string) string {
	result := strings.Split(html, `<meta property="og:image" content="`)[1]
	titleImage := strings.Split(result, `"><meta `)[0]
	titleImage = strings.Replace(titleImage, "amp;", "", -1)
	return titleImage
}
