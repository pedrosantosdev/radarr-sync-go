package compress

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pedrosantosdev/radarr-sync-go/src/io_archive"
	"github.com/thoas/go-funk"
)

func Handler(source, target string, list []string) error {
	if source == "" || target == "" {
		return fmt.Errorf("missing arguments: source and target required")
	}
	fmt.Println("Init Mapping Folder")
	var diff []string
	for _, name := range list {
		filename := filepath.Base(name)
		targetFile := io_archive.FileStat(filename, io_archive.Extension, target)
		if targetFile == nil {
			diff = append(diff, name)
			continue
		}
		sourceFolder := io_archive.FileStat(filename, "", target)
		if sourceFolder.ModTime().After(targetFile.ModTime()) {
			diff = append(diff, name)
		}
	}
	fmt.Println("Verify Diff Target to Source")
	compressedFiles, err := io_archive.FindWildcard(target, fmt.Sprintf("*.%s", io_archive.Extension))
	if err != nil {
		return err
	}
	for _, compressed := range compressedFiles {
		filename := fmt.Sprintf("/movies/%s", strings.Replace(filepath.Base(compressed), fmt.Sprintf(".%s", io_archive.Extension), "", -1))
		fmt.Printf("%s %d \n", filename, funk.IndexOfString(list, filename))
		if funk.IndexOfString(list, filename) < 0 {
			err := os.Remove(compressed)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Verify Diff Source to Target")
	for _, name := range diff {
		fullPathToCompress := filepath.Join(source, name)
		fmt.Println("Compressing: ", fullPathToCompress)
		err := io_archive.GZIP(fullPathToCompress, target)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	return nil
}
