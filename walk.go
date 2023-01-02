package scanner

import (
	"io/fs"
	"path/filepath"

	"os"
)

type WalkFilter func(path string, d fs.DirEntry) (bool, error)
type WalkErrorHandler func(path string, d fs.DirEntry, err error) error
type WalkFileHandler func(path string, d fs.DirEntry)

func Walk(rootDir string, filter WalkFilter, fileHandler WalkFileHandler, errorHandler WalkErrorHandler) {

	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err == nil {

			includeFile, err := filter(path, d)
			if err == nil && includeFile {
				fileHandler(filepath.Join(rootDir, path), d)
				// fileHandler(path, d)
			}
		} else {
			err = errorHandler(path, d, err)
		}

		return err
	}

	filesystem := os.DirFS(rootDir)
	fs.WalkDir(filesystem, ".", walkDirFunc)

}
