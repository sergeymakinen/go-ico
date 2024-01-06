package icondir_test

import (
	"bytes"
	"encoding/binary"
	"image/color"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/icondir"
	"github.com/sergeymakinen/go-ico/internal/testutil"
)

func TestDecoder_DecodeAllIcon(t *testing.T) {
	b := testutil.Icon.MustRead()
	d := icondir.NewDecoder(bytes.NewReader(b), true)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeDir() = %v; want nil", err)
	}
	entries, mm, err := d.DecodeAll()
	if err != nil {
		t.Fatalf("Decoder.DecodeAll() = _, _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Icon, entries, mm)
}

func TestDecoder_DecodeAllCursor(t *testing.T) {
	b := testutil.Cursor.MustRead()
	d := icondir.NewDecoder(bytes.NewReader(b), false)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeDir() = %v; want nil", err)
	}
	entries, mm, err := d.DecodeAll()
	if err != nil {
		t.Fatalf("Decoder.DecodeAll() = _, _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Cursor, entries, mm)
}

func TestDecoder_BestConfigIcon(t *testing.T) {
	d := icondir.NewDecoder(bytes.NewReader(testutil.Icon.MustRead()), true)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeDir() = %v; want nil", err)
	}
	e, err := d.Best()
	if err != nil {
		t.Fatalf("Decoder.Best() = _, %v; want nil", err)
	}
	if expected := 256; e.Width != expected {
		t.Errorf("Entry.Width = %d; want %d", e.Width, expected)
	}
	if expected := 256; e.Height != expected {
		t.Errorf("Entry.Height = %d; want %d", e.Height, expected)
	}
	if expected := 32; e.BPP != expected {
		t.Errorf("Entry.BPP = %d; want %d", e.BPP, expected)
	}
	config, err := d.DecodeConfig(e)
	if err != nil {
		t.Fatalf("Decoder.DecodeConfig() = _, %v; want nil", err)
	}
	if expected := 256; config.Width != expected {
		t.Errorf("image.Config.Width = %d; want %d", config.Width, expected)
	}
	if expected := 256; config.Height != expected {
		t.Errorf("image.Config.Height = %d; want %d", config.Height, expected)
	}
	if config.ColorModel != color.NRGBAModel {
		t.Errorf("image.Config.ColorModel = %T; want color.NRGBAModel", config.ColorModel)
	}
}

func TestDecoder_BestConfigCursor(t *testing.T) {
	d := icondir.NewDecoder(bytes.NewReader(testutil.Cursor.MustRead()), false)
	if err := d.DecodeDir(); err != nil {
		t.Fatalf("Decoder.DecodeDir() = %v; want nil", err)
	}
	e, err := d.Best()
	if err != nil {
		t.Fatalf("Decoder.Best() = _, %v; want nil", err)
	}
	if expected := 128; e.Width != expected {
		t.Errorf("Entry.Width = %d; want %d", e.Width, expected)
	}
	if expected := 128; e.Height != expected {
		t.Errorf("Entry.Height = %d; want %d", e.Height, expected)
	}
	if expected := 32; e.BPP != expected {
		t.Errorf("Entry.BPP = %d; want %d", e.BPP, expected)
	}
	config, err := d.DecodeConfig(e)
	if err != nil {
		t.Fatalf("Decoder.DecodeConfig() = _, %v; want nil", err)
	}
	if expected := 128; config.Width != expected {
		t.Errorf("image.Config.Width = %d; want %d", config.Width, expected)
	}
	if expected := 128; config.Height != expected {
		t.Errorf("image.Config.Height = %d; want %d", config.Height, expected)
	}
	if config.ColorModel != color.RGBAModel {
		t.Errorf("image.Config.ColorModel = %T; want color.RGBAModel", config.ColorModel)
	}
}

func TestDecoder_DecodeDirShouldFailIcon(t *testing.T) {
	b := make([]byte, 23)
	expect := func(msg string) {
		d := icondir.NewDecoder(bytes.NewReader(b), true)
		if err := d.DecodeDir(); err == nil || err.Error() != msg {
			t.Fatalf("Decoder.DecodeDir() = %v; want %s", err, msg)
		}
	}
	expect("not an ICO file")
	b[2] = 1
	expect("no icons")
	binary.LittleEndian.PutUint16(b[4:], 1)
	expect("unexpected EOF")
}

func TestDecoder_DecodeDirShouldFailCursor(t *testing.T) {
	b := make([]byte, 23)
	expect := func(msg string) {
		d := icondir.NewDecoder(bytes.NewReader(b), false)
		if err := d.DecodeDir(); err == nil || err.Error() != msg {
			t.Fatalf("Decoder.DecodeDir() = %v; want %s", err, msg)
		}
	}
	expect("not a CUR file")
	b[2] = 2
	expect("no cursors")
	binary.LittleEndian.PutUint16(b[4:], 1)
	expect("unexpected EOF")
}
