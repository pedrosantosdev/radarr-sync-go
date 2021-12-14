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
	source := flag.String("source", "", "")
	target := flag.String("target", "", "")
	flag.Parse()
	if *url == "" || *source == "" || *target == "" {
		log.Fatal("Missing arguments: url, soruce and target required")
	}
	movies, err := client.FetchMoviesList(*url)
	if err != nil {
		log.Fatal(err)
	}
	var listMovies []string
	for _, movie := range movies {
		listMovies = append(listMovies, movie.Path)
	}
	compress.Handler(*source, *target, listMovies)
}
