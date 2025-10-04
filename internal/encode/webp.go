package encode

import (
	"image"
	"io"
	"github.com/HugoSmits86/nativewebp"
)

func EncodeWebP(w io.Writer, img image.Image) error {
	return nativewebp.Encode(w, img, nil)
}
