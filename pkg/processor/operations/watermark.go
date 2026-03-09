package operations

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// AddTextWatermark накладывает текст поверх изображения.
// Для простоты используем встроенный шрифт или загружаем из файла.
func AddTextWatermark(img image.Image, text string) (image.Image, error) {
	// Создаём RGBA-копию, чтобы рисовать поверх
	dst := imaging.Clone(img)
	bounds := dst.Bounds()

	// Загружаем шрифт (можно встроить простой ttf или использовать системный)
	// Здесь для примера используем встроенный через go-bindata или заглушку.
	// В реальности лучше загрузить из файла, указанного в конфиге.
	fontBytes, err := loadFont() // заглушка – нужно реализовать
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// Параметры текста
	size := 20.0
	col := image.White
	point := fixed.Point26_6{
		X: fixed.Int26_6(20 * 64),
		Y: fixed.Int26_6(float64(bounds.Dy()-20) * 64),
	}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: truetype.NewFace(f, &truetype.Options{Size: size}),
		Dot:  point,
	}
	d.DrawString(text)

	return dst, nil
}

// loadFont – заглушка, в реальности читает TTF-файл.
func loadFont() ([]byte, error) {
	// Можно загрузить из embed или из файла.
	// Пока вернём nil – нужно реализовать отдельно.
	return nil, nil
}
