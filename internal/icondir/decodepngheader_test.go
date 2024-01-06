package icondir

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

func TestDecodePNGHeader(t *testing.T) {
	files, err := filepath.Glob("testdata/*.png")
	if err != nil {
		panic("failed to list test files: " + err.Error())
	}
	re := regexp.MustCompile(`basn(\d)\w(\d+)\.png$`)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			b, err := os.ReadFile(file)
			if err != nil {
				panic("failed to read " + file + ": " + err.Error())
			}
			var e Entry
			if err = decodePNGHeader(bytes.NewReader(b), &e); err != nil {
				t.Fatalf("decodePNGHeader() = %v; want nil", err)
			}
			matches := re.FindStringSubmatch(file)
			ct, _ := strconv.Atoi(matches[1])
			bpp, _ := strconv.Atoi(matches[2])
			switch ct {
			case 2:
				bpp *= 3
			case 4:
				bpp *= 2
			case 6:
				bpp *= 4
			}
			if e.BPP != bpp {
				t.Errorf("Entry.BPP = %d; want %d", bpp, e.BPP)
			}
		})
	}
}
