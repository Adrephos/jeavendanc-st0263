package utils_files

import (
	"fmt"
	"log"
	"os"
)

func DirFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	var files []string
	if err != nil {
		log.Println(err)
	} else {
		for _, entry := range entries {
			files = append(files, entry.Name())
		}
	}
	return files
}

func FileMetadata(dir string, filename string) string {
	path := dir + filename
	if dir[len(dir)-1] != os.PathSeparator {
		path = dir + string(os.PathSeparator) + filename
	}
	return fmt.Sprint(os.Stat(path))
}

func CreateFile(dir string, filename string) {
	path := dir + filename
	if dir[len(dir)-1] != os.PathSeparator {
		path = dir + string(os.PathSeparator) + filename
	}

	os.Create(path)
}
