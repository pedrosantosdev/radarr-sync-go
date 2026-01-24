package model

import (
	"testing"
)

func TestMovieResponseType(t *testing.T) {
	movie := MovieToRadarrResponse{
		Title:   "Test Movie",
		TmdbId:  12345,
		Year:    "2024",
		HasFile: true,
	}

	if movie.Title != "Test Movie" {
		t.Errorf("Expected title 'Test Movie', got '%s'", movie.Title)
	}

	if movie.TmdbId != 12345 {
		t.Errorf("Expected TmdbId 12345, got %d", movie.TmdbId)
	}

	if movie.Year != "2024" {
		t.Errorf("Expected year '2024', got '%s'", movie.Year)
	}

	if !movie.HasFile {
		t.Error("Expected HasFile to be true")
	}
}

func TestRadarrModelType(t *testing.T) {
	radarr := RadarrModel{
		Title:     "Inception",
		Overview:  "A movie about dreams",
		TmdbId:    27205,
		Path:      "/movies/Inception",
		HasFile:   true,
		InCinemas: "2010-07-16",
		Images: []ImageModel{
			{
				RemoteUrl: "https://example.com/image.jpg",
			},
		},
	}

	if radarr.Title != "Inception" {
		t.Errorf("Expected title 'Inception', got '%s'", radarr.Title)
	}

	if radarr.TmdbId != 27205 {
		t.Errorf("Expected TmdbId 27205, got %d", radarr.TmdbId)
	}

	if len(radarr.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(radarr.Images))
	}

	if radarr.Images[0].RemoteUrl != "https://example.com/image.jpg" {
		t.Errorf("Expected RemoteUrl 'https://example.com/image.jpg', got '%s'", radarr.Images[0].RemoteUrl)
	}
}

func TestMovieResponseSlice(t *testing.T) {
	movies := MovieResponse{
		{
			Title: "Movie 1",
			Path:  "/path/to/movie1",
		},
		{
			Title: "Movie 2",
			Path:  "/path/to/movie2",
		},
	}

	if len(movies) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(movies))
	}

	if movies[0].Title != "Movie 1" {
		t.Errorf("Expected first movie title 'Movie 1', got '%s'", movies[0].Title)
	}
}

func TestAddMovieToRadarrModel(t *testing.T) {
	movie := MovieToRadarrResponse{
		Title:   "Test",
		TmdbId:  123,
		Year:    "2024",
		HasFile: false,
	}

	addMovie := AddMovieToRadarrModel(movie)

	if addMovie.Title != "Test" {
		t.Errorf("Expected title 'Test', got '%s'", addMovie.Title)
	}

	if addMovie.TmdbId != 123 {
		t.Errorf("Expected TmdbId 123, got %d", addMovie.TmdbId)
	}
}

func TestMovieLoginResponse(t *testing.T) {
	login := MovieLoginResponse{
		Token: "test-token-123",
	}

	if login.Token != "test-token-123" {
		t.Errorf("Expected token 'test-token-123', got '%s'", login.Token)
	}
}

func TestRadarrResponseError(t *testing.T) {
	errors := RadarrResponseError{
		{
			PropertyName:   "path",
			ErrorMessage:   "Path is required",
			AttemptedValue: nil,
			Severity:       "error",
			ErrorCode:      "PathRequired",
		},
	}

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	if errors[0].PropertyName != "path" {
		t.Errorf("Expected PropertyName 'path', got '%s'", errors[0].PropertyName)
	}

	if errors[0].ErrorCode != "PathRequired" {
		t.Errorf("Expected ErrorCode 'PathRequired', got '%s'", errors[0].ErrorCode)
	}
}
