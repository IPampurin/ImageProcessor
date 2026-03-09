package operations

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

func Resize(img image.Image, width, height int) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: %dx%d", width, height)
	}
	return imaging.Resize(img, width, height, imaging.Lanczos), nil
}

// EncodeJPEG encodes image to JPEG bytes.
func EncodeJPEG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	return buf.Bytes(), err
}

// EncodePNG encodes image to PNG bytes.
func EncodePNG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	return buf.Bytes(), err
}
