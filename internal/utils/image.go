package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func ImageSize(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	b := img.Bounds()
	return b.Dx(), b.Dy(), nil
}
