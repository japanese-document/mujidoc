package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// CopyDir copies a directory from src to dist. It supports both Windows and Unix-like operating systems.
// On Windows, it uses robocopy, and on Unix-like systems, it uses cp. It returns an error with a stack trace
func CopyDir(src, dist string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		dist := filepath.Join(dist, IMAGE_DIR)
		cmd = exec.Command("robocopy", src, dist, "/E")
	} else {
		cmd = exec.Command("cp", "-r", src, dist)
	}

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// IsDirExists checks if the directory at the given path exists. It returns true if the directory exists,
// and false otherwise.
func IsDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
