package compress

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pedrosantosdev/radarr-sync-go/src/io_archive"
)

// Tests for SyncAndCompress function

func TestSyncAndCompressValidateSourceEmpty(t *testing.T) {
	targetDir := t.TempDir()

	err := SyncAndCompress("", targetDir, []string{})

	if err == nil {
		t.Error("Expected error for empty source")
	}
}

func TestSyncAndCompressValidateTargetEmpty(t *testing.T) {
	sourceDir := t.TempDir()

	err := SyncAndCompress(sourceDir, "", []string{})

	if err == nil {
		t.Error("Expected error for empty target")
	}
}

func TestSyncAndCompressEmptyFileList(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	err := SyncAndCompress(sourceDir, targetDir, []string{})

	if err != nil {
		t.Fatalf("Expected no error for empty file list, got %v", err)
	}
}

func TestSyncAndCompressSourceNotFound(t *testing.T) {
	targetDir := t.TempDir()

	err := SyncAndCompress("/non/existent/source", targetDir, []string{"test"})

	if err == nil {
		t.Error("Expected error for non-existent source")
	}
}

func TestSyncAndCompressTargetNotFound(t *testing.T) {
	sourceDir := t.TempDir()

	testFile := filepath.Join(sourceDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = SyncAndCompress(sourceDir, "/non/existent/target", []string{"test.txt"})

	if err == nil {
		t.Error("Expected error for non-existent target")
	}
}

func TestSyncAndCompressDoesNotRecompressAbsolutePaths(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	movieDir := filepath.Join(sourceDir, "movies", "Film")
	if err := os.MkdirAll(movieDir, 0o755); err != nil {
		t.Fatalf("Failed to create movie directory: %v", err)
	}

	movieFile := filepath.Join(movieDir, "film.mkv")
	if err := os.WriteFile(movieFile, []byte("movie content"), 0o644); err != nil {
		t.Fatalf("Failed to create movie file: %v", err)
	}

	oldTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(movieFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set movie modification time: %v", err)
	}

	moviePaths := []string{"/movies/Film/film.mkv"}

	if err := SyncAndCompress(sourceDir, targetDir, moviePaths); err != nil {
		t.Fatalf("First compression failed: %v", err)
	}

	archivePath := filepath.Join(targetDir, "film.mkv."+io_archive.Extension)
	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("Expected archive to exist: %v", err)
	}

	// Mark the archive with a timestamp newer than the source. If the archive
	// gets deleted and recompressed, Compress resets the timestamp to the source
	// modification time and this sentinel is lost.
	sentinel := oldTime.Add(24 * time.Hour)
	if err := os.Chtimes(archivePath, sentinel, sentinel); err != nil {
		t.Fatalf("Failed to set archive sentinel time: %v", err)
	}

	if err := SyncAndCompress(sourceDir, targetDir, moviePaths); err != nil {
		t.Fatalf("Second compression failed: %v", err)
	}

	secondInfo, err := os.Stat(archivePath)
	if err != nil {
		t.Fatalf("Expected archive to still exist: %v", err)
	}

	if !secondInfo.ModTime().Equal(sentinel) {
		t.Error("Expected archive to not be recompressed when source did not change")
	}
}

func TestIdentifyFilesToCompressOrdersLightestFirst(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	files := map[string]string{
		"movies/Large/large.mkv":   "large content",
		"movies/Small/small.mkv":   "s",
		"movies/Medium/medium.mkv": "medium content here",
	}

	var moviePaths []string
	for relativePath, content := range files {
		fullPath := filepath.Join(sourceDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		moviePaths = append(moviePaths, "/"+relativePath)
	}

	identified, err := identifyFilesToCompress(sourceDir, targetDir, moviePaths)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(identified) != len(files) {
		t.Fatalf("Expected %d files, got %d", len(files), len(identified))
	}

	for i := 1; i < len(identified); i++ {
		if identified[i-1].size > identified[i].size {
			t.Error("Expected files to be ordered from lightest to heaviest")
		}
	}
}

func TestSyncAndCompressMultipleFilesConcurrently(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	var moviePaths []string
	for i := 0; i < 6; i++ {
		relativePath := filepath.Join("movies", "Film"+string(rune('A'+i)))
		fullPath := filepath.Join(sourceDir, relativePath)
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			t.Fatalf("Failed to create movie directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join(fullPath, "film.mkv"), []byte("content"), 0o644); err != nil {
			t.Fatalf("Failed to create movie file: %v", err)
		}
		moviePaths = append(moviePaths, "/"+filepath.ToSlash(relativePath))
	}

	if err := SyncAndCompress(sourceDir, targetDir, moviePaths); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	archives, err := io_archive.FindWildcard(targetDir, "*."+io_archive.Extension)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	if len(archives) != len(moviePaths) {
		t.Errorf("Expected %d archives, got %d", len(moviePaths), len(archives))
	}
}
