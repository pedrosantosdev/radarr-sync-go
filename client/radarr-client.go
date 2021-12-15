package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/pedrosantosdev/radarr-sync-go/model"
)

// curl -H "Content-Type: application/json" -X POST -d '{"title":"Proof","qualityProfileId":"4", "tmdbid":"14904","titleslug":"proof-14904", "monitored":"true", "rootFolderPath":"/Volume1/Movies/", "images":[{"covertype":"poster","url":"https://image.tmdb.org/t/p/w640/ghPbOsvg8WrJQBSThtNakBGuDi4.jpg"}]}' http://192.168.1.10:8310/api/movie?apikey=YOUAPIKEYHERE
func AddMovieOnRadarr(baseUrl, token string, data model.AddMovieToRadarrModel) (string, error) {
	URL := fmt.Sprintf("%s/api/v3/movie?apikey=%s", baseUrl, token)
	onlyWords := regexp.MustCompile(`\w+`)
	FolderPath := fmt.Sprintf("/movies/%s (%s)", onlyWords.ReplaceAllString(data.Title, " "), data.Year)

	// "qualityProfileId":"6" HD - 720p/1080p
	values := map[string]interface{}{
		"tmdbid":           data.TmdbId,
		"path":             FolderPath,
		"monitored":        true,
		"qualityProfileId": "6",
	}
	fmt.Println(values)
	jsonData, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	resp, err := Client.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	// An error is returned if something goes wrong
	if err != nil {
		return "", err
	}
	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	} else {
		return "", nil
	}
}
