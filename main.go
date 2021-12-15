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
	flag.Parse()
	if *url == "" || *source == "" || *target == "" || *login == "" || *password == "" {
		log.Fatal("Missing arguments: url, soruce and target required")
	}
	token, err := client.Login(*url, *login, *password)
	if err != nil {
		log.Fatal(err)
	}
	syncWithRadarr(*url, token.Token)
	// compressNSyncRemote(*url, token.Token, *source, *target)
}

func syncWithRadarr(url, token string) {
	movies, err := client.FetchMoviesListToSync(url, token)
	if err != nil {
		log.Fatal(err)
	}
	for _, movie := range movies {
		fmt.Println(movie.TmdbId)
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
}
