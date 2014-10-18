package infrastructure

import (
	"bytes"
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

func (handler *ImageUtilsHandler) Save(img image.Image, format string) ([]byte, error) {
	// Based on github.com/disintegration/imaging/helpers.go Save()
	var err error
	writer := new(bytes.Buffer)
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
			err = jpeg.Encode(writer, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 95})
		}

	case ".png":
		err = png.Encode(writer, img)
	case ".tif", ".tiff":
		err = tiff.Encode(writer, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case ".bmp":
		err = bmp.Encode(writer, img)
	default:
		return nil, errors.New("Invalid image format")
	}

	return writer.Bytes(), err
}
