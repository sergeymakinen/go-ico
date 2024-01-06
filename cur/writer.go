package cur

import (
	"image"
	"io"

	"github.com/sergeymakinen/go-ico/internal/icondir"
)

// EncodeAll writes the cursors in mm to w in CUR format.
func EncodeAll(w io.Writer, mm []image.Image) error {
	e := icondir.NewEncoder(w, false)
	for _, m := range mm {
		if err := e.Add(m, 0, 0); err != nil {
			return convertErr(err)
		}
	}
	return e.Encode()
}

// Encode writes the cursor m to w in CUR format.
func Encode(w io.Writer, m image.Image) error {
	e := icondir.NewEncoder(w, false)
	if err := e.Add(m, 0, 0); err != nil {
		return convertErr(err)
	}
	return e.Encode()
}
