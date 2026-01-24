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

func TestGetFileInfoExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename+".tar.gz")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info, err := io_archive.GetFileInfo(tmpDir, filename, "tar.gz")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info == nil {
		t.Error("Expected FileInfo, got nil")
	}

	if info.Name() != filename+".tar.gz" {
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

func TestGetFileInfoNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()

	info, err := io_archive.GetFileInfo(tmpDir, "nonexistent", "tar.gz")

	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got %v", err)
	}

	if info != nil {
		t.Error("Expected nil FileInfo for non-existent file")
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

func TestGetFileInfoWithoutExtension(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info, err := io_archive.GetFileInfo(tmpDir, filename, "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info == nil {
		t.Error("Expected FileInfo, got nil")
	}

	if info.Name() != filename {
		t.Errorf("Expected file name '%s', got '%s'", filename, info.Name())
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

func TestGetFileInfoEmptyExtension(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "testfile"

	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info, err := io_archive.GetFileInfo(tmpDir, filename, "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info == nil {
		t.Error("Expected FileInfo for file without extension")
	}
}

func TestGetFileInfoValidationEmptyRoot(t *testing.T) {
	_, err := io_archive.GetFileInfo("", "file", "txt")

	if err == nil {
		t.Error("Expected error for empty root")
	}
}

func TestGetFileInfoValidationEmptyFilename(t *testing.T) {
	_, err := io_archive.GetFileInfo("/tmp", "", "txt")

	if err == nil {
		t.Error("Expected error for empty filename")
	}
}

func TestFindWildcardValidationEmptyRoot(t *testing.T) {
	_, err := io_archive.FindWildcard("", "*.txt")

	if err == nil {
		t.Error("Expected error for empty root")
	}
}

func TestFindWildcardValidationEmptyPattern(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := io_archive.FindWildcard(tmpDir, "")

	if err == nil {
		t.Error("Expected error for empty pattern")
	}

func TestGZIPFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	err := GZIP("/path/to/nonexistent/file", tmpDir)

	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}

func TestCompressValidateSourceEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := io_archive.Compress("", tmpDir, nil)

	if err == nil {
		t.Error("Expected error for empty source")
	}
}

func TestCompressValidateTargetEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := io_archive.Compress(tmpDir, "", nil)

	if err == nil {
		t.Error("Expected error for empty target")
	}
}

func TestCompressSourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := io_archive.Compress("/path/to/nonexistent/file", tmpDir, nil)

	if err == nil {
		t.Error("Expected error for non-existent source")
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

func TestCompressTargetNotDirectory(t *testing.T) {
	sourceDir := t.TempDir()
	targetFile := filepath.Join(t.TempDir(), "file.txt")

	err := os.WriteFile(targetFile, []byte("target"), 0644)
	if err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	_, err = io_archive.Compress(sourceDir, targetFile, nil)

	if err == nil {
		t.Error("Expected error when target is not a directory")
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

	

func TestCompressSimpleFile(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := os.WriteFile(sourceFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	outputPath, err := io_archive.Compress(sourceFile, targetDir, nil)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the tar.gz file was created
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Expected compressed file to have content")
	}

	expectedName := "test.txt." + Extension
	if filepath.Base(outputPath) != expectedName {
		t.Errorf("Expected output name '%s', got '%s'", expectedName, filepath.Base(outputPath))
	}
}err = GZIP(sourceFile, targetDir)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	

func TestCompressDirectory(t *testing.T) {
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

	outputPath, err := io_archive.Compress(testDir, targetDir, nil)

	if err != nil {
		t.Fatalf("Expected no error for directory compression, got %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Expected compressed file to have content")
	}

	expectedName := "testdir." + Extension
	if filepath.Base(outputPath) != expectedName {
		t.Errorf("Expected output name '%s', got '%s'", expectedName, filepath.Base(outputPath))
	}
}

func TestCompressWithCompressionLevel(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := os.WriteFile(sourceFile, []byte("test content test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	opts := &io_archive.CompressOptions{CompressionLevel: 1}
	outputPath, err := io_archive.Compress(sourceFile, targetDir, opts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Expected compressed file to exist, got error: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Expected compressed file to have content")
	}
}// Verify the tar.gz file was created
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
