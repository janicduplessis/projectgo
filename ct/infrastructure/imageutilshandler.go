package infrastructure

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
)

type ImageUtilsHandler struct {
}

func (handler *ImageUtilsHandler) Load(reader io.Reader) (image.Image, error) {
	image, _, err := image.Decode(reader)

	return image, err
}

func (handler *ImageUtilsHandler) Resize(image image.Image, maxWidth int, maxHeight int) image.Image {

	return imaging.Thumbnail(image, maxWidth, maxHeight, imaging.Lanczos)
}

func (handler *ImageUtilsHandler) Save(image image.Image, fileName string) error {
	return imaging.Save(image, fileName)
}
