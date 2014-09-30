package interfaces

import (
	"image"
	"io"
)

type ImageUtils interface {
	Load(reader io.Reader) (image.Image, error)
	Resize(image image.Image, maxWidth int, maxHeight int) image.Image
	Save(image image.Image, fileName string) error
}
