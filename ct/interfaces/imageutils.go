package interfaces

import (
	"image"
	"io"
)

// ImageUtils interface for image manipulations
type ImageUtils interface {
	Load(reader io.Reader) (image.Image, error)
	Resize(image image.Image, maxWidth int, maxHeight int) image.Image
	Save(image image.Image, fileName string) error
}
