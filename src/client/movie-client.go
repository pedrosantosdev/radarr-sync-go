//nolint:goimports
package client

import (
	"fmt"

	"github.com/pedrosantosdev/radarr-sync-go/src/model"
)

var serverURI string

func SetServerUri(baseUrl string) {
	serverURI = baseUrl
}

func FetchMoviesListToCompress(token string) (model.MovieResponse, error) {
	URL := fmt.Sprintf("%s/movies/sync", serverURI)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	var cResp model.MovieResponse

	err := SendRequest("GET", URL, &cResp, nil, headers)
	if err != nil {
		return nil, err
	}

	return cResp, nil
}

func FetchMoviesListToSync(token string) ([]model.MovieToRadarrResponse, error) {
	URL := fmt.Sprintf("%s/movies", serverURI)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	var cResp []model.MovieToRadarrResponse

	err := SendRequest("GET", URL, &cResp, nil, headers)
	if err != nil {
		return nil, err
	}

	return cResp, nil
}

func AddMovieToServer(token string, data *model.RadarrModel) error {
	URL := fmt.Sprintf("%s/movies", serverURI)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	inCinemas := "TBA"

	if data != nil && data.InCinemas != "" {
		inCinemas = data.InCinemas[0:10]
	}

	body := map[string]interface{}{
		"tmdbId":    data.TmdbId,
		"title":     data.Title,
		"overview":  data.Overview,
		"path":      data.Path,
		"hasFile":   data.HasFile,
		"image":     data.Images[0].RemoteUrl,
		"inCinemas": inCinemas,
		"needSync":  false,
	}

	var cResp interface{}
	err := SendRequest("POST", URL, &cResp, body, headers)
	if err != nil {
		return err
	}

	return nil
}

func Login(login, password string) (model.MovieLoginResponse, error) {
	URL := fmt.Sprintf("%s/login", serverURI)
	var cResp model.MovieLoginResponse
	data := map[string]string{
		"username": login,
		"password": password,
	}
	err := SendFormEncoded(URL, &cResp, data)
	if err != nil {
		return model.MovieLoginResponse{}, err
	}

	return cResp, nil
}