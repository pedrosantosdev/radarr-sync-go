package client

import (
	"testing"
)

func TestSetServerUri(t *testing.T) {
	expectedURL := "http://localhost:8080"
	SetServerUri(expectedURL)

	if serverURI != expectedURL {
		t.Errorf("Expected serverURI '%s', got '%s'", expectedURL, serverURI)
	}
}

func TestSetServerUriDifferentURL(t *testing.T) {
	url1 := "http://server1.com"
	SetServerUri(url1)
	if serverURI != url1 {
		t.Errorf("Expected serverURI '%s', got '%s'", url1, serverURI)
	}

	url2 := "http://server2.com"
	SetServerUri(url2)
	if serverURI != url2 {
		t.Errorf("Expected serverURI '%s', got '%s'", url2, serverURI)
	}
}

func TestSetRadarrUri(t *testing.T) {
	baseUrl := "http://localhost:7878"
	token := "test-api-key-123"

	SetRadarrUri(baseUrl, token)

	expectedURI := "http://localhost:7878/api/v3/movie?apikey=test-api-key-123"
	if radarrURI != expectedURI {
		t.Errorf("Expected radarrURI '%s', got '%s'", expectedURI, radarrURI)
	}
}

func TestSetRadarrUriWithSpecialCharacters(t *testing.T) {
	baseUrl := "http://example.com:7878"
	token := "key-with-special-chars-!@#$"

	SetRadarrUri(baseUrl, token)

	expectedURI := "http://example.com:7878/api/v3/movie?apikey=key-with-special-chars-!@#$"
	if radarrURI != expectedURI {
		t.Errorf("Expected radarrURI '%s', got '%s'", expectedURI, radarrURI)
	}
}

func TestLoginIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}

func TestFetchMoviesListToSyncIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}

func TestFetchMoviesListToCompressIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}

func TestAddMovieToServerIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}

func TestAddMovieOnRadarrIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}

func TestGetAllMoviesOnRadarrIntegration(t *testing.T) {
	// This would require mocking HTTP response
	t.Skip("Integration test - requires HTTP mock server")
}
