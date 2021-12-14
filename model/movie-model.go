package model

type MovieResponse []struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
