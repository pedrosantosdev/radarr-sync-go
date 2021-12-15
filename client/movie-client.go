package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pedrosantosdev/radarr-sync-go/model"
)

var Client = &http.Client{}

// TODO: client set baseUrl, Header Auth between request https://pkg.go.dev/net/http?utm_source=gopls
func FetchMoviesListToCompress(baseUrl, token string) (model.MovieResponse, error) {
	URL := fmt.Sprintf("%s/movies/sync", baseUrl)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := Client.Do(req)
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

func FetchMoviesListToSync(baseUrl, token string) (model.MovieToRadarrResponse, error) {
	URL := fmt.Sprintf("%s/movies", baseUrl)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp model.MovieToRadarrResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}

	return cResp, nil
}

func Login(baseUrl, login, password string) (model.MovieLoginResponse, error) {
	URL := fmt.Sprintf("%s/login", baseUrl)
	resp, err := http.PostForm(URL, url.Values{
		"username": {login},
		"password": {password},
	})
	if err != nil {
		return model.MovieLoginResponse{}, err
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp model.MovieLoginResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return model.MovieLoginResponse{}, err
	}

	return cResp, nil
}
