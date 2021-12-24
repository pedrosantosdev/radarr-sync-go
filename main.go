package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
		log.Fatalln("Missing arguments: url, soruce, target, radarrUrl, radarrKey required")
		os.Exit(1)
	}
	token, err := client.Login(*url, *login, *password)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	syncWithRadarr(*url, token.Token, *radarrUrl, *radarrKey)
	compressNSyncRemote(*url, token.Token, *source, *target)
	fmt.Println("Finish app")
	log.Default()
	os.Exit(0)
}

func syncWithRadarr(url, token, radarrUrl, radarrKey string) {
	movies, err := client.FetchMoviesListToSync(url, token)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	for _, movie := range movies {
		if movie.HasFile {
			continue
		}
		err := client.AddMovieOnRadarr(radarrUrl, radarrKey, movie)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func compressNSyncRemote(url, token, source, target string) {
	movies, err := client.FetchMoviesListToCompress(url, token)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	var listMovies []string
	for _, movie := range movies {
		listMovies = append(listMovies, movie.Path)
	}
	e := compress.Handler(source, target, listMovies)
	if e != nil {
		log.Fatalln(e)
		os.Exit(1)
	}
}
