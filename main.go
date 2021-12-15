package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pedrosantosdev/radarr-sync-go/client"
	"github.com/pedrosantosdev/radarr-sync-go/compress"
)

func main() {
	fmt.Println("Init app")
	url := flag.String("url", "", "")
	login := flag.String("login", "", "")
	password := flag.String("password", "", "")
	source := flag.String("source", "", "")
	target := flag.String("target", "", "")
	radarrUrl := flag.String("radarr-url", "", "")
	radarrKey := flag.String("radarr-key", "", "")
	flag.Parse()
	if *url == "" || *source == "" || *target == "" || *login == "" || *password == "" || *radarrUrl == "" || *radarrKey == "" {
		log.Fatal("Missing arguments: url, soruce, target, radarrUrl, radarrKey required")
	}
	token, err := client.Login(*url, *login, *password)
	if err != nil {
		log.Fatal(err)
	}
	syncWithRadarr(*url, token.Token, *radarrUrl, *radarrKey)
	compressNSyncRemote(*url, token.Token, *source, *target)
}

func syncWithRadarr(url, token, radarrUrl, radarrKey string) {
	movies, err := client.FetchMoviesListToSync(url, token)
	if err != nil {
		log.Fatal(err)
	}
	for _, movie := range movies {
		_, err := client.AddMovieOnRadarr(radarrUrl, radarrKey, movie)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func compressNSyncRemote(url, token, source, target string) {
	movies, err := client.FetchMoviesListToCompress(url, token)
	if err != nil {
		log.Fatal(err)
	}
	var listMovies []string
	for _, movie := range movies {
		listMovies = append(listMovies, movie.Path)
	}
	compress.Handler(source, target, listMovies)
	if err != nil {
		log.Fatal(err)
	}
}
