package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pedrosantosdev/radarr-sync-go/src/io_archive"
)

// SyncAndCompress synchronizes compressed archives and compresses new files.
//
// Logic:
// 1. Removes compressed files in target that are not in moviePaths
// 2. Identifies files in moviePaths that need compression (don't exist or are outdated)
// 3. Compresses identified files
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

	// Create map for O(1) lookup
	movieSet := make(map[string]bool)
	for _, path := range moviePaths {
		movieSet[path] = true
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

// identifyFilesToCompress returns list of files that need to be compressed.
// A file needs compression if:
// - Compressed file doesn't exist
// - Original file is newer than compressed file
func identifyFilesToCompress(source, target string, moviePaths []string) ([]string, error) {
	var needsCompress []string

	for _, moviePath := range moviePaths {
		filename := filepath.Base(moviePath)

		// Check if compressed file exists
		compressedInfo, err := io_archive.GetFileInfo(target, filename, io_archive.Extension)
		if err != nil {
			return nil, fmt.Errorf("failed to check compressed file for %s: %w", moviePath, err)
		}

		if compressedInfo == nil {
			// Doesn't exist, needs compression
			needsCompress = append(needsCompress, moviePath)
			continue
		}

		// Check if original file is newer
		originalInfo, err := io_archive.GetFileInfo(source, filename, "")
		if err != nil {
			return nil, fmt.Errorf("failed to check source file for %s: %w", moviePath, err)
		}

		if originalInfo == nil {
			// Original doesn't exist, skip
			continue
		}

		// If original is newer, needs recompression
		if originalInfo.ModTime().After(compressedInfo.ModTime()) {
			needsCompress = append(needsCompress, moviePath)
		}
	}

	return needsCompress, nil
}

// compressFiles compresses list of files from source to target.
func compressFiles(source, target string, moviePaths []string) error {
	for _, moviePath := range moviePaths {
		fullPath := filepath.Join(source, moviePath)

		outputPath, err := io_archive.Compress(fullPath, target, nil)
		if err != nil {
			return fmt.Errorf("failed to compress %s: %w", moviePath, err)
		}

		fmt.Printf("Compressed: %s -> %s\n", moviePath, filepath.Base(outputPath))
	}

	return nil
}
