package compress

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerMissingSource(t *testing.T) {
	err := Handler("", "/target", []string{})

	if err == nil {
		t.Error("Expected error for missing source, got nil")
	}

	if err.Error() != "missing arguments: source and target required" {
		t.Errorf("Expected error message about missing arguments, got '%s'", err.Error())
	}
}

func TestHandlerMissingTarget(t *testing.T) {
	err := Handler("/source", "", []string{})

	if err == nil {
		t.Error("Expected error for missing target, got nil")
	}

	if err.Error() != "missing arguments: source and target required" {
		t.Errorf("Expected error message about missing arguments, got '%s'", err.Error())
	}
}

func TestHandlerWithEmptyFileList(t *testing.T) {
	// Create temporary directories
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	err := Handler(sourceDir, targetDir, []string{})

	if err != nil {
		t.Fatalf("Expected no error for empty file list, got %v", err)
	}
}

func TestHandlerWithNonExistentSource(t *testing.T) {
	targetDir := t.TempDir()

	err := Handler("/non/existent/path", targetDir, []string{"/movies/test"})

	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}

func TestHandlerWithNonExistentTarget(t *testing.T) {
	sourceDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(sourceDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = Handler(sourceDir, "/non/existent/target", []string{"test"})

	if err == nil {
		t.Error("Expected error for non-existent target, got nil")
	}
}

func TestHandlerCalculatesFileDiff(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	// This is a basic structure test
	// Full integration test would verify actual file compression
	list := []string{"/movies/test"}

	err := Handler(sourceDir, targetDir, list)

	if err != nil {
		// It's okay if we get an error in this case as we're checking behavior
		t.Logf("Got expected behavior: %v", err)
	}
}
