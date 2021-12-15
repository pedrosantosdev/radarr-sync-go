package model

type MovieResponse []struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type MovieToRadarrResponse []struct {
	TmdbId string `json:"tmdbId"`
}

type MovieLoginResponse struct {
	Token string `json:"token"`
}
