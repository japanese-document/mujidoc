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

// ImageSize takes the path to an image file and returns the dimensions of the image.
// Supported image formats include JPEG, PNG, GIF, BMP, TIFF, and WEBP, thanks to
// the standard image package and additional formats provided by the golang.org/x/image package.
//
// Parameters:
// - path: A string representing the filesystem path to the image file.
//
// Returns:
// - The width (Dx) and height (Dy) of the image as integers.
// - An error if the file cannot be opened, or if the image format is not recognized or is corrupt.
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
