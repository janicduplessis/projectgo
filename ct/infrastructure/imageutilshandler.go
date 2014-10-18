package infrastructure

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"code.google.com/p/go.image/bmp"
	"code.google.com/p/go.image/tiff"
	"github.com/disintegration/imaging"
)

type ImageUtilsHandler struct {
}

func (handler *ImageUtilsHandler) Load(file io.Reader) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (handler *ImageUtilsHandler) Resize(image image.Image, maxWidth int, maxHeight int) image.Image {

	return imaging.Thumbnail(image, maxWidth, maxHeight, imaging.Lanczos)
}

func (handler *ImageUtilsHandler) Save(file io.Writer, img image.Image, format string) error {
	// Based on github.com/disintegration/imaging/helpers.go Save()
	var err error
	switch format {
	case ".jpg", ".jpeg":
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
		}

	case ".png":
		err = png.Encode(file, img)
	case ".tif", ".tiff":
		err = tiff.Encode(file, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case ".bmp":
		err = bmp.Encode(file, img)
	default:
		return errors.New("Invalid image format")
	}

	return err
}
