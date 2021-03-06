package interfaces

import (
	"image"
	"io"
)

// ImageUtils interface for image manipulations
type ImageUtils interface {
	Load(file io.Reader) (image.Image, error)
	Resize(image image.Image, maxWidth int, maxHeight int) image.Image
	Save(img image.Image, format string) ([]byte, error)
}
