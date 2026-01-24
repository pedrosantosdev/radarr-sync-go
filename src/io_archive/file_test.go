package io_archive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWildcardInNestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directories
	subDir1 := filepath.Join(tmpDir, "level1")
	subDir2 := filepath.Join(subDir1, "level2")

	err := os.MkdirAll(subDir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Create files at different levels
	files := map[string]string{
		filepath.Join(tmpDir, "file1.tar.gz"):      "test",
		filepath.Join(subDir1, "file2.tar.gz"):     "test",
		filepath.Join(subDir2, "file3.tar.gz"):     "test",
		filepath.Join(subDir2, "other.txt"):        "test",
	}

	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
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

func TestFileStatModificationTime(t *testing.T) {
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

	if info.ModTime().IsZero() {
		t.Error("Expected valid ModTime")
	}
}

func TestGZIPPreservesFileContent(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	testContent := []byte("This is test content for verification")
	sourceFile := filepath.Join(sourceDir, "content.txt")

	err := os.WriteFile(sourceFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	err = GZIP(sourceFile, targetDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	compressedFile := filepath.Join(targetDir, "content.txt."+Extension)
	stat, err := os.Stat(compressedFile)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if stat.Size() == 0 {
		t.Error("Compressed file is empty")
	}
}

func TestExtensionConstant(t *testing.T) {
	if Extension != "tar.gz" {
		t.Errorf("Expected Extension constant 'tar.gz', got '%s'", Extension)
	}
}
