package encode

import (
	"image"
	"io"
	"image/jpeg"
)

func EncodeJPEG(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
}
