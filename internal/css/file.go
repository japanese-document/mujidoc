package css

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateWriteTask returns a closure function that when executed, creates or ensures the existence
// of the specified output directory and writes CSS content to a file within that directory.
// The function takes two parameters: 'outputDirDir' is the path to the output directory where the CSS file
// should be created, and 'fileName' is the name of the CSS file to be created.
func CreateWriteTask(outputDirDir, fileName string) func() error {
	return func() error {
		err := os.MkdirAll(outputDirDir, os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}
		cssFileName := filepath.Join(outputDirDir, fileName)
		file, err := os.Create(cssFileName)
		if err != nil {
			return errors.WithStack(err)
		}
		defer file.Close()

		_, err = file.WriteString(CSS_CONTENT)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func Version() string {
	return uuid.NewSHA1(uuid.Nil, []byte(CSS_CONTENT)).String()
}
