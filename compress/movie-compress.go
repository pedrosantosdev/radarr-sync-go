package compress

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const Extension = "tar"

func Tar(source, target string) error {
	filename := filepath.Base(source)
	fmt.Println("Compressing: ", filename)
	target = filepath.Join(target, fmt.Sprintf("%s.%s", filename, Extension))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
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

func Compress(listNames []string, source, target string) {
	for _, name := range listNames {
		s := filepath.Join(source, name)
		err := Tar(s, target)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

func Missing(filename, target string, isDirectory bool) bool {
	var search string
	if isDirectory {
		search = fmt.Sprintf("%s/%s", target, filename)
	} else {
		search = fmt.Sprintf("%s/%s.%s", target, filename, Extension)
	}
	if _, err := os.Stat(search); err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func WalkMatch(root, pattern string) ([]string, error) {
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
		return fmt.Errorf("Missing arguments: source and target required")
	}
	fmt.Println("Init Mapping Folder")
	var diff []string
	for _, name := range list {
		filename := filepath.Base(name)
		if Missing(filename, target, false) {
			diff = append(diff, name)
		}
	}
	fmt.Println("Verify Diff Target to Source")
	compressedFiles, err := WalkMatch(target, fmt.Sprintf("*.%s", Extension))
	if err != nil {
		return err
	}
	for _, compressed := range compressedFiles {
		filename := strings.Replace(filepath.Base(compressed), fmt.Sprintf(".%s", Extension), "", -1)
		folder := fmt.Sprintf("movies/%s", filename)
		if Missing(folder, source, true) {
			err := os.Remove(compressed)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Verify Diff Source to Target")
	Compress(diff, source, target)
	return nil
}
