package testutil

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/icondir"
)

type entry struct {
	prefix string

	Width, Height, BPP, XHotspot, YHotspot int
}

func (e entry) MustDecode() image.Image {
	file := filepath.Join(dir, "testdata", e.String())
	b, err := os.ReadFile(file)
	if err != nil {
		panic("failed to read " + file + ": " + err.Error())
	}
	m, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		panic("failed to decode " + file + ": " + err.Error())
	}
	return m
}

func (e entry) String() string {
	return fmt.Sprintf("%s_%dx%d-%d.png", e.prefix, e.Width, e.Height, e.BPP)
}

type IconDir struct {
	Name    string
	Entries []entry
}

func (d IconDir) MustRead() []byte {
	file := filepath.Join(dir, "testdata", d.Name)
	b, err := os.ReadFile(file)
	if err != nil {
		panic("failed to read " + file + ": " + err.Error())
	}
	return b
}

var (
	Icon = newIconDir("icon.ico", []entry{
		{Width: 64, Height: 64, BPP: 1},
		{Width: 16, Height: 16, BPP: 1},
		{Width: 64, Height: 64, BPP: 4},
		{Width: 32, Height: 32, BPP: 4},
		{Width: 16, Height: 16, BPP: 4},
		{Width: 64, Height: 64, BPP: 8},
		{Width: 32, Height: 32, BPP: 8},
		{Width: 16, Height: 16, BPP: 8},
		{Width: 64, Height: 64, BPP: 24},
		{Width: 32, Height: 32, BPP: 24},
		{Width: 16, Height: 16, BPP: 24},
		{Width: 256, Height: 256, BPP: 32},
		{Width: 64, Height: 64, BPP: 32},
		{Width: 32, Height: 32, BPP: 32},
		{Width: 16, Height: 16, BPP: 32},
	})
	Cursor = newIconDir("cursor.cur", []entry{
		{Width: 128, Height: 128, BPP: 32, XHotspot: 20, YHotspot: 20},
		{Width: 64, Height: 64, BPP: 32, XHotspot: 10, YHotspot: 10},
		{Width: 48, Height: 48, BPP: 32, XHotspot: 5, YHotspot: 5},
		{Width: 32, Height: 32, BPP: 32, XHotspot: 3, YHotspot: 3},
	})
)

func CompareIconDir(t *testing.T, dir IconDir, entries []*icondir.Entry, mm []image.Image) {
	if actual, expected := len(mm), len(dir.Entries); actual != expected {
		t.Fatalf("len([]image.Image) = %d; want %d", actual, expected)
	}
	for i, e := range dir.Entries {
		t.Run(e.String(), func(t *testing.T) {
			m := e.MustDecode()
			Compare(t, m, mm[i])
			if entries != nil {
				if entries[i].BPP != e.BPP {
					t.Errorf("Entry.BPP = %d; want %d", entries[i].BPP, e.BPP)
				}
				if e.XHotspot > 0 || e.YHotspot > 0 {
					if entries[i].XHotspot != e.XHotspot {
						t.Errorf("Entry.XHotspot = %d; want %d", entries[i].XHotspot, e.XHotspot)
					}
					if entries[i].YHotspot != e.YHotspot {
						t.Errorf("Entry.YHotspot = %d; want %d", entries[i].YHotspot, e.YHotspot)
					}
				}
			}
		})
	}
}

func Compare(t *testing.T, expected, actual image.Image) {
	if !expected.Bounds().Eq(actual.Bounds()) {
		t.Fatalf("Bounds() = %s; want %s", actual.Bounds(), expected.Bounds())
	}
	errors := 0
	for y := expected.Bounds().Min.Y; y < expected.Bounds().Max.Y; y++ {
		for x := expected.Bounds().Min.X; x < expected.Bounds().Max.X; x++ {
			expectedR, expectedG, expectedB, expectedA := expected.At(x, y).RGBA()
			actualR, actualG, actualB, actualA := actual.At(x, y).RGBA()
			if expectedR != actualR || expectedG != actualG || expectedB != actualB || expectedA != actualA {
				t.Errorf("At(%d, %d) = %v; want %v", x, y, actual.At(x, y), expected.At(x, y))
				if errors >= 50 {
					t.Fatalf("Too many At() errors")
				} else {
					errors++
				}
			}
		}
	}
}

func newIconDir(name string, entries []entry) IconDir {
	var ee []entry
	for _, e := range entries {
		ee = append(ee, entry{
			prefix: strings.TrimSuffix(name, filepath.Ext(name)),
			Width:  e.Width,
			Height: e.Height,
			BPP:    e.BPP,
		})
	}
	return IconDir{
		Name:    name,
		Entries: ee,
	}
}

func Palettize(m image.Image) image.Image {
	var palette color.Palette
	colors := map[color.Color]bool{}
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := m.At(x, y)
			if _, ok := colors[c]; ok {
				continue
			}
			palette = append(palette, c)
			colors[c] = true
		}
	}
	paletted := image.NewPaletted(m.Bounds(), palette)
	draw.Draw(paletted, paletted.Bounds(), m, m.Bounds().Min, draw.Src)
	return paletted
}

var dir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir = filepath.Dir(filename)
	println(dir)
}
