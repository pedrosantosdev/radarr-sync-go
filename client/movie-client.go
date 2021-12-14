package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pedrosantosdev/radarr-sync-go/model"
)

func FetchMoviesList(baseUrl string) (model.MovieResponse, error) {
	URL := fmt.Sprintf("%s/movies/sync", baseUrl)
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp model.MovieResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}

	return cResp, nil
}
