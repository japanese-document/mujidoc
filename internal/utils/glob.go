package utils

import (
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
)

// createWalkDirFunc creates and returns a fs.WalkDirFunc that appends file paths with a specific suffix to the paths slice.
// It returns an error if the provided paths slice pointer is nil, ensuring safety against nil pointer dereference.
// The function is designed to be used with filepath.WalkDir to collect file paths that end with the given suffix.
func createWalkDirFunc(paths *[]string, suffix string) (func(string, fs.DirEntry, error) error, error) {
	if paths == nil {
		return nil, errors.WithStack(errors.New("paths is nil"))
	}
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if d.IsDir() {
			return nil
		}
		lastIndex := len(path) - len(suffix)
		if lastIndex >= 0 && path[lastIndex:] == suffix {
			*paths = append(*paths, path)
		}
		return nil
	}
	return walkDirFunc, nil
}

// GetMarkDownFileNames searches for files with a specified suffix within a root directory and returns their paths.
// It utilizes createWalkDirFunc to generate a function that filters and collects file paths during a directory walk.
// An error is returned if any issues arise during the creation of the walk function or the directory walk itself.
func GetMarkDownFileNames(root, suffix string) ([]string, error) {
	paths := []string{}
	walkDirFunc, err := createWalkDirFunc(&paths, suffix)
	if err != nil {
		return paths, err
	}
	if err := filepath.WalkDir(root, walkDirFunc); err != nil {
		return paths, errors.WithStack(err)
	}
	return paths, nil
}
