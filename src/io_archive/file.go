package io_archive

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FindWildcard searches recursively for files matching pattern in root directory.
// Returns absolute paths of matching files. Returns nil slice if no matches found.
// Errors are returned only for path traversal issues, not missing files.
//
// Example: FindWildcard("/movies", "*.tar.gz") returns ["/movies/file.tar.gz"]
func FindWildcard(root, pattern string) ([]string, error) {
	if root == "" {
		return nil, fmt.Errorf("root path cannot be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	var matches []string
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		matched, patternErr := filepath.Match(pattern, filepath.Base(path))
		if patternErr != nil {
			return patternErr
		}

		if matched {
			matches = append(matches, path)
		}
		return nil
	})

	if walkErr != nil {
		return nil, walkErr
	}
	return matches, nil
}

// GetFileInfo returns file information from root directory with optional extension.
// Returns (nil, nil) if file does not exist. Returns error only for access issues.
//
// Parameters:
// - root: directory path
// - filename: file name (without extension)
// - extension: optional file extension (e.g., "tar.gz" or "" for no extension)
//
// Example: GetFileInfo("/movies", "file", "tar.gz") checks for "/movies/file.tar.gz"
func GetFileInfo(root, filename, extension string) (fs.FileInfo, error) {
	if root == "" {
		return nil, fmt.Errorf("root path cannot be empty")
	}
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	var path string
	if extension != "" {
		path = filepath.Join(root, filename+"."+extension)
	} else {
		path = filepath.Join(root, filename)
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	return info, nil
}

// FileStat is deprecated. Use GetFileInfo instead.
// Kept for backwards compatibility.
func FileStat(filename, fileExtension, root string) fs.FileInfo {
	info, err := GetFileInfo(root, filename, fileExtension)
	if err != nil {
		return nil
	}
	return info
}
