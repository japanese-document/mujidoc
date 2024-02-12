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

// GetMarkDownFileNames searches the specified root directory and all of its subdirectories
// for files with the ".md" extension and returns a slice containing the paths of all markdown files found.
// Parameters:
// - root: The root directory from which the search will begin.
func GetMarkDownFileNames(root string) ([]string, error) {
	paths := []string{}
	walkDirFunc, err := createWalkDirFunc(&paths, ".md")
	if err != nil {
		return paths, err
	}
	if err := filepath.WalkDir(root, walkDirFunc); err != nil {
		return paths, errors.WithStack(err)
	}
	return paths, nil
}
