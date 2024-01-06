package icondir_test

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/icondir"
	"github.com/sergeymakinen/go-ico/internal/testutil"
)

func TestEncoder_EncodeIcon(t *testing.T) {
	var buf bytes.Buffer
	e := icondir.NewEncoder(&buf, true)
	for _, entry := range testutil.Icon.Entries {
		m := entry.MustDecode()
		if entry.BPP <= 8 {
			m = testutil.Palettize(m)
		}
		if err := e.Add(m, 0, 0); err != nil {
			t.Fatalf("Encoder.Add() = %v; want nil", err)
		}
	}
	if err := e.Encode(); err != nil {
		t.Fatalf("Encoder.Encode() = %v; want nil", err)
	}
	d := icondir.NewDecoder(&buf, true)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeAll() = %v; want nil", err)
	}
	entries, mm, err := d.DecodeAll()
	if err != nil {
		t.Fatalf("Decoder.DecodeAll() = _, _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Icon, entries, mm)
}

func TestEncoder_EncodeCursor(t *testing.T) {
	var buf bytes.Buffer
	e := icondir.NewEncoder(&buf, false)
	for _, entry := range testutil.Cursor.Entries {
		if err := e.Add(entry.MustDecode(), 0, 0); err != nil {
			t.Fatalf("Encoder.Add() = %v; want nil", err)
		}
	}
	if err := e.Encode(); err != nil {
		t.Fatalf("Encoder.Encode() = %v; want nil", err)
	}
	d := icondir.NewDecoder(&buf, false)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeAll() = %v; want nil", err)
	}
	entries, mm, err := d.DecodeAll()
	if err != nil {
		t.Fatalf("Decoder.DecodeAll() = _, _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Cursor, entries, mm)
}

func TestEncoder_EncodePNG32Bit(t *testing.T) {
	files, err := filepath.Glob("testdata/basn[246]*.png")
	if err != nil {
		panic("failed to list test files: " + err.Error())
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			b, err := os.ReadFile(file)
			if err != nil {
				panic("failed to read " + file + ": " + err.Error())
			}
			m, err := png.Decode(bytes.NewReader(b))
			if err != nil {
				t.Fatalf("png.Decode() = _, %v; want nil", err)
			}
			var m2 image.Image
			switch m := m.(type) {
			case *image.Gray:
				m2 = image.NewGray(image.Rect(0, 0, 256, 256))
			case *image.Gray16:
				m2 = image.NewGray16(image.Rect(0, 0, 256, 256))
			case *image.RGBA:
				m2 = image.NewRGBA(image.Rect(0, 0, 256, 256))
			case *image.RGBA64:
				m2 = image.NewRGBA64(image.Rect(0, 0, 256, 256))
			case *image.Paletted:
				m2 = image.NewPaletted(image.Rect(0, 0, 256, 256), m.Palette)
			case *image.NRGBA:
				m2 = image.NewNRGBA(image.Rect(0, 0, 256, 256))
			case *image.NRGBA64:
				m2 = image.NewRGBA64(image.Rect(0, 0, 256, 256))
			}
			draw.Draw(m2.(draw.Image), m2.Bounds(), m, m.Bounds().Min, draw.Src)
			var buf bytes.Buffer
			e := icondir.NewEncoder(&buf, true)
			if err := e.Add(m2, 0, 0); err != nil {
				t.Fatalf("Encoder.Add() = %v; want nil", err)
			}
			if err := e.Encode(); err != nil {
				t.Fatalf("Encoder.Encode() = %v; want nil", err)
			}
			d := icondir.NewDecoder(&buf, true)
			if err := d.DecodeDir(); err != nil {
				t.Fatalf("Decoder.DecodeAll() = %v; want nil", err)
			}
			entries, _, err := d.DecodeAll()
			if err != nil {
				t.Fatalf("Decoder.DecodeAll() = _, _, %v; want nil", err)
			}
			if actual, expected := entries[0].BPP, 32; actual != expected {
				t.Errorf("Entry.BPP = %d; want %d", actual, expected)
			}
		})
	}
}
