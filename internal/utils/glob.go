package utils

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type IFileSystem interface {
	WalkDir(root string, fn fs.WalkDirFunc) error
}

type FileSystem struct{}

func (f FileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

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

// validateFileNames checks each file name in the provided slice of file paths for spaces (including tabs and newlines).
//
// Parameters:
// - paths: A slice of strings, each representing a file path to validate.
//
// Returns:
// - An error if any file name contains spaces, tabs, or newlines.
// - nil if all file names are valid and do not contain spaces.
func validateFileNames(paths []string) error {
	for _, path := range paths {
		if strings.ContainsAny(path, " \t\n") {
			return errors.WithStack(fmt.Errorf("there are spaces in file name. the file name is %s", path))
		}
	}
	return nil
}

// GetMarkDownFileNames searches the specified root directory and all of its subdirectories
// for files with the ".md" extension and returns a slice containing the paths of all markdown files found.
// Parameters:
// - root: The root directory from which the search will begin.
func GetMarkDownFileNames(fs IFileSystem, root string) ([]string, error) {
	paths := []string{}
	walkDirFunc, err := createWalkDirFunc(&paths, ".md")
	if err != nil {
		return paths, err
	}
	if err := fs.WalkDir(root, walkDirFunc); err != nil {
		return paths, errors.WithStack(err)
	}
	err = validateFileNames(paths)
	if err != nil {
		return paths, err
	}
	return paths, nil
}
