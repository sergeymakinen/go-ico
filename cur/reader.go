// Package cur implements a CUR file decoder and encoder.
package cur

import (
	"image"
	"io"

	"github.com/sergeymakinen/go-ico/internal/icondir"
)

// FormatError reports that the input is not a valid CUR.
type FormatError string

func (e FormatError) Error() string { return "cur: invalid format: " + string(e) }

// UnsupportedError reports that the input uses a valid but unimplemented CUR feature.
type UnsupportedError string

func (e UnsupportedError) Error() string { return "cur: unsupported feature: " + string(e) }

// Hotspot represents the coordinates of the cursor hotspot.
type Hotspot struct {
	X, Y int
}

// CUR represents the multiple cursors stored in an CUR file.
type CUR struct {
	Cursor  []image.Image
	Hotspot []Hotspot
}

// DecodeAll reads a CUR image from r and returns the stored cursors.
func DecodeAll(r io.Reader) (*CUR, error) {
	d := icondir.NewDecoder(r, false)
	if err := d.DecodeDir(); err != nil {
		return nil, convertErr(err)
	}
	entries, mm, err := d.DecodeAll()
	if err != nil {
		return nil, convertErr(err)
	}
	cur := &CUR{}
	for i := 0; i < len(entries); i++ {
		cur.Cursor = append(cur.Cursor, mm[i])
		cur.Hotspot = append(cur.Hotspot, Hotspot{
			X: entries[i].XHotspot,
			Y: entries[i].YHotspot,
		})
	}
	return cur, nil
}

// Decode reads a CUR image from r and returns the largest stored cursor
// as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	d := icondir.NewDecoder(r, false)
	if err := d.DecodeDir(); err != nil {
		return nil, convertErr(err)
	}
	e, err := d.Best()
	if err != nil {
		return nil, convertErr(err)
	}
	m, err := d.Decode(e)
	if err != nil {
		return nil, convertErr(err)
	}
	return m, nil
}

// DecodeConfig returns the color model and dimensions of the largest cursor
// stored in a CUR image without decoding the entire cursor.
func DecodeConfig(r io.Reader) (image.Config, error) {
	d := icondir.NewDecoder(r, false)
	if err := d.DecodeDir(); err != nil {
		return image.Config{}, convertErr(err)
	}
	e, err := d.Best()
	if err != nil {
		return image.Config{}, convertErr(err)
	}
	config, err := d.DecodeConfig(e)
	if err != nil {
		return image.Config{}, convertErr(err)
	}
	return config, nil
}

func convertErr(err error) error {
	switch err := err.(type) {
	case icondir.FormatError:
		return FormatError(err.Error())
	case icondir.UnsupportedError:
		return UnsupportedError(err.Error())
	default:
		return err
	}
}

func init() {
	image.RegisterFormat("cur", "\x00\x00\x02\x00??", Decode, DecodeConfig)
}
