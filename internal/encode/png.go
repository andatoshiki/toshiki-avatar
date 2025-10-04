
package encode

import (
	"image"
	"io"
	"image/png"
)

func EncodePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
