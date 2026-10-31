package compress

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/pedrosantosdev/radarr-sync-go/src/io_archive"
)

// fileToCompress describes a source path that needs a new archive.
type fileToCompress struct {
	path string
	size int64
}

// SyncAndCompress synchronizes compressed archives and compresses new files.
//
// Logic:
// 1. Removes compressed files in target that are not in moviePaths
// 2. Identifies files in moviePaths that need compression (don't exist or are outdated)
// 3. Compresses identified files concurrently, lightest files first
//
// Parameters:
// - source: directory containing original files
// - target: directory to save compressed files
// - moviePaths: list of relative file paths to compress
//
// Returns error if any operation fails.
func SyncAndCompress(source, target string, moviePaths []string) error {
	if source == "" {
		return fmt.Errorf("source path cannot be empty")
	}
	if target == "" {
		return fmt.Errorf("target path cannot be empty")
	}
	if len(moviePaths) == 0 {
		return nil // Nothing to do
	}

	// Build the set using the archive base name, because archives are stored as
	// "<base name>.tar.gz" regardless of the original directory structure.
	movieSet := make(map[string]bool, len(moviePaths))
	for _, path := range moviePaths {
		movieSet[filepath.Base(path)] = true
	}

	// Phase 1: Remove compressed files not in moviePaths
	if err := cleanupObsoleteArchives(target, movieSet); err != nil {
		return fmt.Errorf("cleanup phase failed: %w", err)
	}

	// Phase 2: Identify files to compress
	needsCompress, err := identifyFilesToCompress(source, target, moviePaths)
	if err != nil {
		return fmt.Errorf("diff phase failed: %w", err)
	}

	// Phase 3: Compress files
	if err := compressFiles(source, target, needsCompress); err != nil {
		return fmt.Errorf("compression phase failed: %w", err)
	}

	return nil
}

// cleanupObsoleteArchives removes compressed files in target that are not in movieSet.
func cleanupObsoleteArchives(target string, movieSet map[string]bool) error {
	compressedFiles, err := io_archive.FindWildcard(target, "*."+io_archive.Extension)
	if err != nil {
		return err
	}

	for _, compressedPath := range compressedFiles {
		filename := filepath.Base(compressedPath)
		// Remove extension to get original filename
		movieName := strings.TrimSuffix(filename, "."+io_archive.Extension)

		if !movieSet[movieName] {
			if err := os.Remove(compressedPath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", compressedPath, err)
			}
		}
	}

	return nil
}

// identifyFilesToCompress returns files that need to be compressed, ordered
// from lightest to heaviest so the fastest jobs finish first.
//
// A file needs compression if:
// - Compressed file doesn't exist
// - Original file is newer than compressed file
func identifyFilesToCompress(source, target string, moviePaths []string) ([]fileToCompress, error) {
	var needsCompress []fileToCompress

	for _, moviePath := range moviePaths {
		filename := filepath.Base(moviePath)

		// Check if compressed file exists
		compressedInfo, err := io_archive.GetFileInfo(target, filename, io_archive.Extension)
		if err != nil {
			return nil, fmt.Errorf("failed to check compressed file for %s: %w", moviePath, err)
		}

		fullPath := filepath.Join(source, moviePath)
		originalInfo, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				if compressedInfo == nil {
					return nil, fmt.Errorf("source file does not exist: %s", fullPath)
				}
				// Original doesn't exist but archive does, skip it
				continue
			}
			return nil, fmt.Errorf("failed to check source file for %s: %w", moviePath, err)
		}

		// Skip when the archive is already up to date
		if compressedInfo != nil && !originalInfo.ModTime().After(compressedInfo.ModTime()) {
			continue
		}

		needsCompress = append(needsCompress, fileToCompress{
			path: moviePath,
			size: sourceSize(fullPath, originalInfo),
		})
	}

	sort.Slice(needsCompress, func(i, j int) bool {
		if needsCompress[i].size == needsCompress[j].size {
			return needsCompress[i].path < needsCompress[j].path
		}
		return needsCompress[i].size < needsCompress[j].size
	})

	return needsCompress, nil
}

// sourceSize returns the total size of a file, or the sum of all regular files
// when the source is a directory.
func sourceSize(path string, info os.FileInfo) int64 {
	if !info.IsDir() {
		return info.Size()
	}

	var total int64
	_ = filepath.WalkDir(path, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil || !entry.Type().IsRegular() {
			return nil
		}
		if fileInfo, infoErr := entry.Info(); infoErr == nil {
			total += fileInfo.Size()
		}
		return nil
	})
	return total
}

// compressFiles compresses files from source to target concurrently.
func compressFiles(source, target string, files []fileToCompress) error {
	files = dedupeByArchiveName(files)
	if len(files) == 0 {
		return nil
	}

	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}
	if workers > len(files) {
		workers = len(files)
	}

	jobs := make(chan fileToCompress)
	errCh := make(chan error, len(files))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				fullPath := filepath.Join(source, file.path)

				outputPath, err := io_archive.Compress(fullPath, target, nil)
				if err != nil {
					errCh <- fmt.Errorf("failed to compress %s: %w", file.path, err)
					continue
				}

				fmt.Printf("Compressed: %s -> %s\n", file.path, filepath.Base(outputPath))
			}
		}()
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// dedupeByArchiveName keeps a single job per output archive, since archives are
// stored flat in target and duplicate names would write to the same file.
func dedupeByArchiveName(files []fileToCompress) []fileToCompress {
	if len(files) < 2 {
		return files
	}

	seen := make(map[string]bool, len(files))
	unique := files[:0]
	for _, file := range files {
		archiveName := filepath.Base(file.path)
		if seen[archiveName] {
			continue
		}
		seen[archiveName] = true
		unique = append(unique, file)
	}
	return unique
}
