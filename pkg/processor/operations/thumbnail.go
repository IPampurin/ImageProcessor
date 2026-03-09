package operations

import (
	"image"

	"github.com/disintegration/imaging"
)

func Thumbnail(img image.Image, width, height int) (image.Image, error) {
	// Используем Fit, чтобы сохранить пропорции и вписать в прямоугольник
	return imaging.Fit(img, width, height, imaging.Lanczos), nil
}
