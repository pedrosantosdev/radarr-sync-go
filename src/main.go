package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pedrosantosdev/radarr-sync-go/src/client"
	"github.com/pedrosantosdev/radarr-sync-go/src/compress"
	"github.com/pedrosantosdev/radarr-sync-go/src/model"
	"github.com/thoas/go-funk"
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
	client.SetServerUri(*url)
	client.SetRadarrUri(*radarrUrl, *radarrKey)
	token, err := client.Login(*login, *password)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	e := syncWithRadarr(token.Token)
	if e != nil {
		log.Fatalln(e)
		os.Exit(1)
	}
	er := compressNSyncRemote(token.Token, *source, *target)
	if er != nil {
		log.Fatalln(er)
		os.Exit(1)
	}
	fmt.Println("Finish app")
	os.Exit(0)
}

func syncWithRadarr(token string) error {
	moviesOnServer, err := client.FetchMoviesListToSync(token)
	if err != nil {
		return err
	}
	moviesOnRadarr, err := client.GetAllMoviesOnRadarr()
	if err != nil {
		return err
	}
	fmt.Println("Server to Radarr")
	for _, movie := range moviesOnServer {
		if movie.HasFile || (funk.IndexOf(moviesOnRadarr, func(value model.RadarrModel) bool {
			return value.TmdbId == movie.TmdbId
		}) > -1) {
			continue
		}
		err := client.AddMovieOnRadarr(movie)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	fmt.Println("Radarr to Server")
	for _, movie := range moviesOnRadarr {
		if !movie.HasFile || (funk.IndexOf(moviesOnServer, func(value model.MovieToRadarrResponse) bool {
			return value.TmdbId == movie.TmdbId
		}) > -1) {
			continue
		}
		// Add To Server
		err := client.AddMovieToServer(token, movie)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func compressNSyncRemote(token, source, target string) error {
	movies, err := client.FetchMoviesListToCompress(token)
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
		return e
	}
	return nil
}
