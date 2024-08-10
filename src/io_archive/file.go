package io_archive

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func FindWildcard(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func FileStat(filename, fileExtension, root string) fs.FileInfo {
	var search string
	if fileExtension == "" {
		search = fmt.Sprintf("%s/%s", root, filename)
	} else {
		search = fmt.Sprintf("%s/%s.%s", root, filename, fileExtension)
	}
	fileInfo, err := os.Stat(search)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	return fileInfo
}
