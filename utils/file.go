package utils

import (
	"os"
	"path/filepath"
)

type FileInfo struct {
	Info os.FileInfo
	Path string
}

func GetDirFiles(dir string) ([]FileInfo, error) {
	var files []FileInfo
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{Info: info, Path: filepath.Join(dir, entry.Name())})
	}

	return files, nil
}
