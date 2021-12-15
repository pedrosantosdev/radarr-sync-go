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
		return nil
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
		Tar(s, target)
	}
}

func Missing(filename, target string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s.%s", target, filename, Extension)); err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func Handler(source, target string, list []string) error {
	if source == "" || target == "" {
		return fmt.Errorf("Missing arguments: source and target required")
	}
	fmt.Println("Init Mapping Folder")
	var diff []string
	for _, name := range list {
		filename := filepath.Base(name)
		if Missing(filename, target) {
			diff = append(diff, name)
		}
	}
	fmt.Println("Verify Diff Source to Target")
	Compress(diff, source, target)
	return nil
}
