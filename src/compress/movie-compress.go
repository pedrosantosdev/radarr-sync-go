package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"
)

const Extension = "tar.gz"

func compress(source, target string) error {
	filename := filepath.Base(source)
	fmt.Println("Compressing: ", filename)
	target = filepath.Join(target, fmt.Sprintf("%s.%s", filename, Extension))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()
	gz, err := gzip.NewWriterLevel(tarfile, 6)
	if err != nil {
		return err
	}
	defer gz.Close()
	tarball := tar.NewWriter(gz)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func compressList(listNames []string, source, target string) {
	for _, name := range listNames {
		s := filepath.Join(source, name)
		err := compress(s, target)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func missing(filename, source, target string, isDirectory bool) bool {
	var search string
	if isDirectory {
		search = fmt.Sprintf("%s/%s", target, filename)
	} else {
		search = fmt.Sprintf("%s/%s.%s", target, filename, Extension)
	}
	fileInfo, err := os.Stat(search)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	folderInfo, err := os.Stat(fmt.Sprintf("%s/%s", source, filename))
	if err != nil {
		return false
	}
	if fileInfo.ModTime().After(folderInfo.ModTime()) {
		return false
	}
	return true
}

func walkMatch(root, pattern string) ([]string, error) {
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

func Handler(source, target string, list []string) error {
	if source == "" || target == "" {
		return fmt.Errorf("missing arguments: source and target required")
	}
	fmt.Println("Init Mapping Folder")
	var diff []string
	for _, name := range list {
		filename := filepath.Base(name)
		if missing(filename, source, target, false) {
			diff = append(diff, name)
		}
	}
	fmt.Println("Verify Diff Target to Source")
	compressedFiles, err := walkMatch(target, fmt.Sprintf("*.%s", Extension))
	if err != nil {
		return err
	}
	for _, compressed := range compressedFiles {
		filename := fmt.Sprintf("/movies/%s", strings.Replace(filepath.Base(compressed), fmt.Sprintf(".%s", Extension), "", -1))
		fmt.Printf("%s %d \n", filename, funk.IndexOfString(list, filename))
		if funk.IndexOfString(list, filename) < 0 {
			err := os.Remove(compressed)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Verify Diff Source to Target")
	compressList(diff, source, target)
	return nil
}
