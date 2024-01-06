package ico_test

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/sergeymakinen/go-ico"
)

func Example() {
	b, _ := os.ReadFile("icon_32x32-32.png")
	m1, _ := png.Decode(bytes.NewReader(b))
	b, _ = os.ReadFile("icon_256x256-32.png")
	m2, _ := png.Decode(bytes.NewReader(b))
	f, _ := os.Create("icon.ico")
	ico.EncodeAll(f, []image.Image{m1, m2})
	f.Close()
}
