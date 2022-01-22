package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/pedrosantosdev/radarr-sync-go/model"
)

// curl -H "Content-Type: application/json" -X POST -d '{"title":"Proof","qualityProfileId":"4", "tmdbid":"14904","titleslug":"proof-14904", "monitored":"true", "rootFolderPath":"/Volume1/Movies/", "images":[{"covertype":"poster","url":"https://image.tmdb.org/t/p/w640/ghPbOsvg8WrJQBSThtNakBGuDi4.jpg"}]}' http://192.168.1.10:8310/api/movie?apikey=YOUAPIKEYHERE
func AddMovieOnRadarr(baseUrl, token string, data model.AddMovieToRadarrModel) error {
	URL := fmt.Sprintf("%s/api/v3/movie?apikey=%s", baseUrl, token)
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
	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}
	resp, err := Client.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	// An error is returned if something goes wrong
	if err != nil {
		return err
	}
	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	//Create a variable of the same type as our model
	var cResp model.RadarrResponseError

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		if err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError

			switch {
			case errors.As(err, &syntaxError):
				msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
				return fmt.Errorf(msg)

			case errors.Is(err, io.ErrUnexpectedEOF):
				msg := "Request body contains badly-formed JSON"
				return fmt.Errorf(msg)

			case errors.As(err, &unmarshalTypeError):
				msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
				return fmt.Errorf(msg)

			case strings.HasPrefix(err.Error(), "json: unknown field "):
				fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
				msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
				return fmt.Errorf(msg)

			case errors.Is(err, io.EOF):
				msg := "Request body must not be empty"
				return fmt.Errorf(msg)

			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
