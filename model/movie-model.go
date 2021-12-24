package model

type MovieResponse []struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type MovieToRadarrResponse []struct {
	Title   string `json:"title"`
	TmdbId  int    `json:"tmdbId"`
	Year    string `json:"year"`
	HasFile bool   `json:"hasFile"`
}

type RadarrResponseError []struct {
	PropertyName   string      `json:"propertyName"`
	ErrorMessage   string      `json:"errorMessage"`
	AttemptedValue interface{} `json:"attemptedValue"`
	Severity       string      `json:"severity"`
	ErrorCode      string      `json:"errorCode"`
}

type AddMovieToRadarrModel struct {
	Title   string `json:"title"`
	TmdbId  int    `json:"tmdbId"`
	Year    string `json:"year"`
	HasFile bool   `json:"hasFile"`
}

type MovieLoginResponse struct {
	Token string `json:"token"`
}
