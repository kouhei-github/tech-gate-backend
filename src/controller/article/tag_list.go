package article

import (
	"encoding/json"
	"net-http/myapp/repository"
	"net-http/myapp/utils"
	"net/http"
)

type TagAllResponse struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func TagListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Method Not Allowed"})
		return
	}

	records, err := repository.FindAllTags()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(utils.MyError{Message: "Internal Server Error"})
		return
	}

	var tags []TagAllResponse
	for _, record := range records {
		if record.ImageURL == "" {
			continue
		}
		tag := TagAllResponse{Url: record.ImageURL, Name: record.Name}
		tags = append(tags, tag)
	}

	json.NewEncoder(w).Encode(tags)
}
