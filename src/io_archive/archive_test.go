package io_archive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWildcardValidPattern(t *testing.T) {
	// Create temporary directory with test files
	tmpDir := t.TempDir()

	// Create test files
	files := []string{
		"file1.tar.gz",
		"file2.tar.gz",
		"file3.txt",
		"file4.tar.gz",
	}

	for _, file := range files {
		filePath := filepath.Join(tmpDir, file)
		err := os.WriteFile(filePath, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	matches, err := FindWildcard(tmpDir, "*.tar.gz")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}
}

func TestFindWildcardNoMatches(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	matches, err := FindWildcard(tmpDir, "*.tar.gz")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestFindWildcardInvalidRoot(t *testing.T) {
	_, err := FindWildcard("/non/existent/path", "*.tar.gz")

	if err == nil {
		t.Error("Expected error for invalid root, got nil")
	}
}

func TestFindWildcardIgnoresDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files
	err = os.WriteFile(filepath.Join(tmpDir, "file.tar.gz"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	matches, err := FindWildcard(tmpDir, "*.tar.gz")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}

func TestFileStatExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename+".tar.gz")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info := FileStat(filename, "tar.gz", tmpDir)

	if info == nil {
		t.Error("Expected FileInfo, got nil")
	}

	if !info.IsDir() && info.Name() != filename+".tar.gz" {
		t.Errorf("Expected file name '%s.tar.gz', got '%s'", filename, info.Name())
	}
}

func TestFileStatNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()

	info := FileStat("nonexistent", "tar.gz", tmpDir)

	if info != nil {
		t.Error("Expected nil for non-existent file, got FileInfo")
	}
}

func TestFileStatWithoutExtension(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info := FileStat(filename, "", tmpDir)

	if info == nil {
		t.Error("Expected FileInfo, got nil")
	}

	if info.Name() != filename {
		t.Errorf("Expected file name '%s', got '%s'", filename, info.Name())
	}
}

func TestFileStatEmptyExtension(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info := FileStat(filename, "", tmpDir)

	if info == nil {
		t.Error("Expected FileInfo for file without extension, got nil")
	}
}

func TestGZIPFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	err := GZIP("/path/to/nonexistent/file", tmpDir)

	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}

func TestGZIPTargetDirectoryNotAccessible(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.txt")

	err := os.WriteFile(sourceFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Try to write to a non-existent directory
	err = GZIP(sourceFile, "/root/no/permission")

	if err == nil {
		t.Error("Expected error for non-accessible target directory")
	}
}

func TestGZIPSimpleFile(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := os.WriteFile(sourceFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	err = GZIP(sourceFile, targetDir)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the tar.gz file was created
	expectedFile := filepath.Join(targetDir, "test.txt."+Extension)
	info, err := os.Stat(expectedFile)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Expected compressed file to have content")
	}
}

func TestGZIPDirectory(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	// Create a source directory with files
	testDir := filepath.Join(sourceDir, "testdir")
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	err = GZIP(testDir, targetDir)

	if err != nil {
		t.Fatalf("Expected no error for directory compression, got %v", err)
	}

	expectedFile := filepath.Join(targetDir, "testdir."+Extension)
	info, err := os.Stat(expectedFile)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Expected compressed file to have content")
	}
}
