package model

type MovieResponse []struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type MovieToRadarrResponse struct {
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

type AddMovieToRadarrModel MovieToRadarrResponse

type GetMovieRadarrModel []RadarrModel

type RadarrModel struct {
	Title     string       `json:"title"`
	Overview  string       `json:"overview"`
	TmdbId    int          `json:"tmdbId"`
	Path      string       `json:"path"`
	HasFile   bool         `json:"hasFile"`
	InCinemas string       `json:"inCinemas"`
	Images    []ImageModel `json:"images"`
}

type ImageModel struct {
	RemoteUrl string `json:"remoteUrl"`
}

type MovieLoginResponse struct {
	Token string `json:"token"`
}
