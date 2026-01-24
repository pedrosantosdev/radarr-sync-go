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

// httpClient is reused with connection pooling for better performance
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// HTTPClient returns the global reusable HTTP client with timeout.
// Use this client for all requests to reuse connections.
func HTTPClient() *http.Client {
	return httpClient
}

// decodeJSON decodes response JSON with robust error handling.
func decodeJSON(respBody io.ReadCloser, result interface{}) error {
	if err := json.NewDecoder(respBody).Decode(&result); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown field %s", fieldName)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("request body must not be empty")

		default:
			return fmt.Errorf("failed to decode JSON: %w", err)
		}
	}
	return nil
}

// SendRequest sends HTTP request with automatic JSON encoding/decoding.
//
// Parameters:
// - method: HTTP method (GET, POST, PUT, DELETE)
// - endpoint: full URL to request
// - result: pointer to struct to decode response
// - data: optional request body (only for POST/PUT)
// - headers: optional custom headers
//
// Returns error if request fails or status code is not 2xx.
func SendRequest(method, endpoint string, result interface{}, data interface{}, headers map[string]string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	body, err := encodeRequestBody(method, data)
	if err != nil {
		return err
	}

	req, err := createRequest(method, endpoint, body, headers)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Decode response
	if err := decodeJSON(resp.Body, result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func encodeRequestBody(method string, data interface{}) ([]byte, error) {
	if data != nil && (method == "POST" || method == "PUT") {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		return body, nil
	}
	return []byte{}, nil
}

func createRequest(method, endpoint string, body []byte, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

// SendFormEncoded sends form-encoded POST request.
//
// Parameters:
// - endpoint: full URL to request
// - result: pointer to struct to decode response
// - data: form data
//
// Returns error if request fails or status code is not 2xx.
func SendFormEncoded(endpoint string, result interface{}, data map[string]string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	body := url.Values{}
	for key, value := range data {
		body.Add(key, value)
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Decode response
	if err := decodeJSON(resp.Body, result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// StructToMap converts struct to map using JSON marshaling.
// This is useful for dynamic field access.
//
// Note: This operation has overhead due to marshaling/unmarshaling.
// For performance-critical code, consider using reflection directly.
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	if obj == nil {
		return nil, errors.New("cannot convert nil to map")
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return result, nil
}
