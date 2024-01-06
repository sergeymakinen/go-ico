package cur

import (
	"bytes"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/testutil"
)

func TestDecodeAll(t *testing.T) {
	b := testutil.Cursor.MustRead()
	cur, err := DecodeAll(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("DecodeAll() = _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Cursor, nil, cur.Cursor)
}

func TestDecode(t *testing.T) {
	b := testutil.Cursor.MustRead()
	m, err := Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("Decode() = _, %v; want nil", err)
	}
	testutil.Compare(t, testutil.Cursor.Entries[0].MustDecode(), m)
}

func TestDecodeShouldFail(t *testing.T) {
	b := make([]byte, 10)
	_, err := Decode(bytes.NewReader(b))
	if expected := "cur: invalid format: not a CUR file"; err == nil || err.Error() != expected {
		t.Fatalf("Decode() = _, %v; want %s", err, expected)
	}
}

func TestDecodeConfig(t *testing.T) {
	b := testutil.Cursor.MustRead()
	config, err := DecodeConfig(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("DecodeConfig() = _, %v; want nil", err)
	}
	if expected := 128; config.Width != expected {
		t.Errorf("image.Config.Width = %d; want %d", config.Width, expected)
	}
	if expected := 128; config.Height != expected {
		t.Errorf("image.Config.Height = %d; want %d", config.Height, expected)
	}
}

func TestDecodeConfigShouldFail(t *testing.T) {
	b := make([]byte, 10)
	_, err := DecodeConfig(bytes.NewReader(b))
	if expected := "cur: invalid format: not a CUR file"; err == nil || err.Error() != expected {
		t.Fatalf("DecodeConfig() = _, %v; want %s", err, expected)
	}
}
