package io_archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const Extension = "tar.gz"

// CompressOptions configures compression options
type CompressOptions struct {
	// CompressionLevel sets gzip compression level (1-9, default 7)
	CompressionLevel int
}

// Compress creates a tar.gz archive from source to target directory.
// If source is a directory, it compresses recursively preserving structure.
// If source is a file, it compresses just that file.
// Returns the path to created archive file.
//
// Returns error if:
// - source path is empty or does not exist
// - target path is empty or not a directory
// - compression fails
//
// Example: Compress("/data/movies", "/backups", nil)
func Compress(source, target string, opts *CompressOptions) (string, error) {
	if source == "" {
		return "", fmt.Errorf("source path cannot be empty")
	}
	if target == "" {
		return "", fmt.Errorf("target path cannot be empty")
	}

	// Validate source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("source does not exist: %s", source)
		}
		return "", fmt.Errorf("cannot access source: %w", err)
	}

	// Validate target is accessible directory
	targetInfo, err := os.Stat(target)
	if err != nil {
		return "", fmt.Errorf("target not accessible: %w", err)
	}
	if !targetInfo.IsDir() {
		return "", fmt.Errorf("target must be a directory: %s", target)
	}

	// Set compression level
	level := 7
	if opts != nil && opts.CompressionLevel > 0 && opts.CompressionLevel <= 9 {
		level = opts.CompressionLevel
	}

	// Generate output filename
	filename := filepath.Base(source)
	outputPath := filepath.Join(target, filename+"."+Extension)

	// Create archive file
	tarfile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive: %w", err)
	}
	defer tarfile.Close()

	// Create gzip writer
	gz, err := gzip.NewWriterLevel(tarfile, level)
	if err != nil {
		tarfile.Close()
		os.Remove(outputPath)
		return "", fmt.Errorf("failed to create gzip writer: %w", err)
	}
	defer gz.Close()

	// Create tar writer
	writer := tar.NewWriter(gz)
	defer writer.Close()

	var baseDir string
	if sourceInfo.IsDir() {
		baseDir = filepath.Base(source)
	}

	// Walk source and add files to archive
	walkErr := filepath.Walk(source,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk error at %s: %w", path, err)
			}

			// Skip symlinks
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				return nil
			}

			// Create tar header
			header, err := tar.FileInfoHeader(fileInfo, "")
			if err != nil {
				return fmt.Errorf("failed to create header for %s: %w", path, err)
			}

			// Set header name with correct path separators for tar format
			if baseDir != "" {
				relPath, err := filepath.Rel(source, path)
				if err != nil {
					return fmt.Errorf("failed to get relative path: %w", err)
				}
				// Always use forward slash in tar format
				header.Name = filepath.ToSlash(filepath.Join(baseDir, relPath))
			} else {
				header.Name = filepath.ToSlash(filepath.Base(path))
			}

			// Write header
			if err := writer.WriteHeader(header); err != nil {
				return fmt.Errorf("failed to write header for %s: %w", path, err)
			}

			// Skip directory content
			if fileInfo.IsDir() {
				return nil
			}

			// Copy file content
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			copied, err := io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}
			if copied != fileInfo.Size() {
				return fmt.Errorf("incomplete copy of %s: got %d bytes, expected %d", path, copied, fileInfo.Size())
			}

			return nil
		})

	// Handle walk errors
	if walkErr != nil {
		writer.Close()
		gz.Close()
		tarfile.Close()
		os.Remove(outputPath)
		return "", fmt.Errorf("compression failed: %w", walkErr)
	}

	// Close writers in correct order
	closeErr := writer.Close()
	gzipErr := gz.Close()

	if closeErr != nil {
		tarfile.Close()
		os.Remove(outputPath)
		return "", fmt.Errorf("failed to close tar writer: %w", closeErr)
	}
	if gzipErr != nil {
		tarfile.Close()
		os.Remove(outputPath)
		return "", fmt.Errorf("failed to close gzip writer: %w", gzipErr)
	}

	tarfile.Close()
	return outputPath, nil
}

// GZIP is deprecated. Use Compress() instead.
// Kept for backwards compatibility.
func GZIP(source, target string) error {
	_, err := Compress(source, target, nil)
	return err
}
