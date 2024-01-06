package cur

import (
	"bytes"
	"image"
	"testing"

	"github.com/sergeymakinen/go-ico/internal/testutil"
)

func TestEncodeAll(t *testing.T) {
	var mm []image.Image
	for _, entry := range testutil.Cursor.Entries {
		mm = append(mm, entry.MustDecode())
	}
	var buf bytes.Buffer
	if err := EncodeAll(&buf, mm); err != nil {
		t.Fatalf("EncodeAll() = %v; want nil", err)
	}
	cur, err := DecodeAll(&buf)
	if err != nil {
		t.Fatalf("DecodeAll() = _, %v; want nil", err)
	}
	testutil.CompareIconDir(t, testutil.Cursor, nil, cur.Cursor)
}

func TestEncodeAllShouldFail(t *testing.T) {
	mm := []image.Image{&image.Gray{}}
	var buf bytes.Buffer
	err := EncodeAll(&buf, mm)
	if expected := "cur: invalid format: invalid image size: 0x0"; err == nil || err.Error() != expected {
		t.Fatalf("EncodeAll() = %v; want %s", err, expected)
	}
}

func TestEncode(t *testing.T) {
	m := testutil.Cursor.Entries[0].MustDecode()
	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode() = %v; want nil", err)
	}
	m2, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode() = _, %v; want nil", err)
	}
	testutil.Compare(t, m, m2)
}

func TestEncodeShouldFail(t *testing.T) {
	var buf bytes.Buffer
	err := Encode(&buf, &image.Gray{})
	if expected := "cur: invalid format: invalid image size: 0x0"; err == nil || err.Error() != expected {
		t.Fatalf("Encode() = %v; want %s", err, expected)
	}
}
