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
	"time"
)

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

func HttpClient() *http.Client {
	c := &http.Client{Timeout: 10 * time.Second}
	return c
}

func handleJson(respBody io.ReadCloser, cResp interface{}) error {
	if err := json.NewDecoder(respBody).Decode(&cResp); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return fmt.Errorf(msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "request body contains badly-formed JSON"
			return fmt.Errorf(msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return fmt.Errorf(msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			return fmt.Errorf(msg)

		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return fmt.Errorf(msg)

		}
	}
	return nil
}

func SendRequest(c *http.Client, url, method string, cResp interface{}, data *map[string]interface{}, headerMap *map[string]string) error {
	endpoint := url

	var body []byte
	if method == "POST" && data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = jsonData

	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if headerMap != nil {
		for headerKey, headerValue := range *headerMap {
			req.Header.Add(headerKey, headerValue)
		}
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	// Close the connection to reuse it
	defer resp.Body.Close()

	//Decode the data
	if err := handleJson(resp.Body, &cResp); err != nil && (resp.StatusCode < 200 && resp.StatusCode > 299) {
		msg := fmt.Sprintf("%s | %d", err, resp.StatusCode)
		return fmt.Errorf(msg)
	}

	return nil
}

func PostFormEncoded(URL string, cResp interface{}, data map[string]string) error {
	body := url.Values{}
	for key, value := range data {
		body.Add(key, value)
	}
	resp, err := http.PostForm(URL, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := handleJson(resp.Body, &cResp); err != nil || (resp.StatusCode < 200 && resp.StatusCode > 299) {
		return err
	}

	return nil
}
