package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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

func FetchMoviesListToSync(baseUrl, token string) ([]model.MovieToRadarrResponse, error) {
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
	var cResp []model.MovieToRadarrResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}

	return cResp, nil
}

func AddMovieToServer(baseUrl, token string, data model.GetMovieRadarrModel) error {
	URL := fmt.Sprintf("%s/movies", baseUrl)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := Client.Do(req)
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
