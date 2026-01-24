package client

import (
	"fmt"
	"regexp"

	"github.com/pedrosantosdev/radarr-sync-go/src/model"
)

var radarrURI string

func SetRadarrUri(baseUrl, token string) {
	radarrURI = fmt.Sprintf("%s/api/v3/movie?apikey=%s", baseUrl, token)
}

func AddMovieOnRadarr(data model.MovieToRadarrResponse) error {
	onlyWords := regexp.MustCompile(`\W+`)
	FolderPath := fmt.Sprintf("/movies/%s (%s)", onlyWords.ReplaceAllString(data.Title, " "), data.Year)

	// "qualityProfileId":6 HD - 720p/1080p
	values := map[string]interface{}{
		"tmdbid":              data.TmdbId,
		"path":                FolderPath,
		"monitored":           true,
		"qualityProfileId":    6,
		"minimumAvailability": 2,
	}

	var cResp model.RadarrResponseError

	err := SendRequest("POST", radarrURI, &cResp, values, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetAllMoviesOnRadarr() (model.GetMovieRadarrModel, error) {
	var cResp model.GetMovieRadarrModel

	err := SendRequest("GET", radarrURI, &cResp, nil, nil)
	if err != nil {
		return nil, err
	}
	return cResp, nil
}
