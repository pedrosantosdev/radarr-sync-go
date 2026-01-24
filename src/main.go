package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pedrosantosdev/radarr-sync-go/src/client"
	"github.com/pedrosantosdev/radarr-sync-go/src/compress"
	"github.com/pedrosantosdev/radarr-sync-go/src/model"
)

// Flag names
const (
	flagURL         = "url"
	flagLogin       = "login"
	flagPassword    = "password"
	flagSource      = "source"
	flagTarget      = "target"
	flagRadarrURL   = "radarr-url"
	flagRadarrKey   = "radarr-key"
	flagSkipCompress = "skip-compress"
	flagDebug       = "debug"
)

func main() {
	fmt.Println("Init app")
	url := flag.String(flagURL, "", "Server URL")
	login := flag.String(flagLogin, "", "Server username")
	password := flag.String(flagPassword, "", "Server password")
	source := flag.String(flagSource, "", "Folder with files to compress")
	target := flag.String(flagTarget, "", "Target path for compressed files")
	radarrUrl := flag.String(flagRadarrURL, "", "URL from Radarr")
	radarrKey := flag.String(flagRadarrKey, "", "API key from Radarr")
	skipCompress := flag.Bool(flagSkipCompress, false, "Skip the compression stage")
	debug := flag.Bool(flagDebug, false, "Enable debug mode")
	flag.Parse()

	if err := validateFlags(*url, *login, *password, *radarrUrl, *radarrKey,
		!*skipCompress, *source, *target); err != nil {
		log.Fatalf("Validation error: %v\n", err)
	}

	client.SetServerUri(*url)
	client.SetRadarrUri(*radarrUrl, *radarrKey)

	token, err := client.Login(*login, *password)
	if err != nil {
		log.Fatalf("Login failed: %v\n", err)
	}

	if err := syncWithRadarr(token.Token, *debug); err != nil {
		log.Fatalf("Sync failed: %v\n", err)
	}

	if !*skipCompress {
		if err := compressNSyncRemote(token.Token, *source, *target); err != nil {
			log.Fatalf("Compression failed: %v\n", err)
		}
	}

	fmt.Println("Finish app")
}

func syncWithRadarr(token string, debug bool) error {
	moviesOnServer, err := client.FetchMoviesListToSync(token)
	if err != nil {
		return fmt.Errorf("fetch server movies failed: %w", err)
	}

	moviesOnRadarr, err := client.GetAllMoviesOnRadarr()
	if err != nil {
		return fmt.Errorf("fetch radarr movies failed: %w", err)
	}

	if err := syncServerToRadarr(moviesOnServer, moviesOnRadarr, debug); err != nil {
		return err
	}

	return syncRadarrToServer(token, moviesOnServer, moviesOnRadarr, debug)
}

func syncServerToRadarr(moviesOnServer []model.MovieToRadarrResponse,
	moviesOnRadarr []model.RadarrModel, debug bool) error {
	fmt.Println("Syncing: Server to Radarr")
	for _, movie := range moviesOnServer {
		if debug {
			fmt.Printf("  Processing: %s\n", movie.Title)
		}

		// Skip if already has file or exists on Radarr
		if movie.HasFile || movieExistsOnRadarr(movie.TmdbId, moviesOnRadarr) {
			continue
		}

		if err := client.AddMovieOnRadarr(movie); err != nil {
			fmt.Printf("  Error adding %s to Radarr: %v\n", movie.Title, err)
			continue
		}
	}
	return nil
}

func syncRadarrToServer(token string, moviesOnServer []model.MovieToRadarrResponse,
	moviesOnRadarr []model.RadarrModel, debug bool) error {
	fmt.Println("Syncing: Radarr to Server")
	for _, movie := range moviesOnRadarr {
		if debug {
			fmt.Printf("  Processing: %s\n", movie.Title)
		}

		// Skip if no file on Radarr or already exists on server
		if !movie.HasFile || movieExistsOnServer(movie.TmdbId, moviesOnServer) {
			continue
		}

		if err := client.AddMovieToServer(token, &movie); err != nil {
			fmt.Printf("  Error adding %s to server: %v\n", movie.Title, err)
			continue
		}
	}
	return nil
}

// Helper functions

// validateFlags validates required command line flags.
func validateFlags(url, login, password, radarrUrl, radarrKey string,
	needSourceTarget bool, source, target string) error {
	if url == "" {
		return fmt.Errorf("url is required")
	}
	if login == "" {
		return fmt.Errorf("login is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if radarrUrl == "" {
		return fmt.Errorf("radarr-url is required")
	}
	if radarrKey == "" {
		return fmt.Errorf("radarr-key is required")
	}
	if needSourceTarget {
		if source == "" {
			return fmt.Errorf("source is required when compression is enabled")
		}
		if target == "" {
			return fmt.Errorf("target is required when compression is enabled")
		}
	}
	return nil
}

// movieExistsOnRadarr checks if a movie with given tmdbId exists on Radarr.
func movieExistsOnRadarr(tmdbId int, movies model.GetMovieRadarrModel) bool {
	for _, movie := range movies {
		if movie.TmdbId == tmdbId {
			return true
		}
	}
	return false
}

// movieExistsOnServer checks if a movie with given tmdbId exists on server.
func movieExistsOnServer(tmdbId int, movies []model.MovieToRadarrResponse) bool {
	for _, movie := range movies {
		if movie.TmdbId == tmdbId {
			return true
		}
	}
	return false
}

func compressNSyncRemote(token, source, target string) error {
	movies, err := client.FetchMoviesListToCompress(token)
	if err != nil {
		return fmt.Errorf("fetch movies list failed: %w", err)
	}

	var listMovies []string
	for _, movie := range movies {
		listMovies = append(listMovies, movie.Path)
	}

	if err := compress.SyncAndCompress(source, target, listMovies); err != nil {
		return fmt.Errorf("sync and compress failed: %w", err)
	}

	return nil
}
