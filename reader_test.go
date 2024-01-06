package ico

import (
	"bytes"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/testutil"
)

func TestDecodeAll(t *testing.T) {
	b := testutil.Icon.MustRead()
	mm, err := DecodeAll(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("DecodeAll() = _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Icon, nil, mm)
}

func TestDecode(t *testing.T) {
	b := testutil.Icon.MustRead()
	m, err := Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("Decode() = _, %v; want nil", err)
	}
	testutil.Compare(t, testutil.Icon.Entries[11].MustDecode(), m)
}

func TestDecodeShouldFail(t *testing.T) {
	b := make([]byte, 10)
	_, err := Decode(bytes.NewReader(b))
	if expected := "ico: invalid format: not an ICO file"; err == nil || err.Error() != expected {
		t.Fatalf("Decode() = _, %v; want %s", err, expected)
	}
}

func TestDecodeConfig(t *testing.T) {
	b := testutil.Icon.MustRead()
	config, err := DecodeConfig(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("DecodeConfig() = _, %v; want nil", err)
	}
	if expected := 256; config.Width != expected {
		t.Errorf("image.Config.Width = %d; want %d", config.Width, expected)
	}
	if expected := 256; config.Height != expected {
		t.Errorf("image.Config.Height = %d; want %d", config.Height, expected)
	}
}

func TestDecodeConfigShouldFail(t *testing.T) {
	b := make([]byte, 10)
	_, err := DecodeConfig(bytes.NewReader(b))
	if expected := "ico: invalid format: not an ICO file"; err == nil || err.Error() != expected {
		t.Fatalf("DecodeConfig() = _, %v; want %s", err, expected)
	}
}
