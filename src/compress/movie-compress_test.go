package compress

import (
	"os"
	"path/filepath"
	"testing"
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
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = SyncAndCompress(sourceDir, "/non/existent/target", []string{"test.txt"})

	if err == nil {
		t.Error("Expected error for non-existent target")
	}
}
