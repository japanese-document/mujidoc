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
		cmd = exec.Command("robocopy", src, dist, "/E")
	} else {
		cmd = exec.Command("cp", "-r", src, dist)
	}

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CreateSrcAndOutputDir constructs and returns the source and output directory paths
// by joining the environment variables "SOURCE_DIR" and "OUTPUT_DIR" with imageDir, respectively.
// It assumes that "SOURCE_DIR" and "OUTPUT_DIR" are set in the environment.
func CreateSrcAndOutputDir() (string, string) {
	return filepath.Join(os.Getenv("SOURCE_DIR"), imageDir), filepath.Join(os.Getenv("OUTPUT_DIR"), imageDir)
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
