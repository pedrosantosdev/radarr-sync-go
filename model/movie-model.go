package model

type MovieResponse []struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type MovieToRadarrResponse []struct {
	Title  string `json:"title"`
	TmdbId int    `json:"tmdbId"`
	Year   string `json:"year"`
}

type AddMovieToRadarrModel struct {
	Title  string `json:"title"`
	TmdbId int    `json:"tmdbId"`
	Year   string `json:"year"`
}

type MovieLoginResponse struct {
	Token string `json:"token"`
}
