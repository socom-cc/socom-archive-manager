package socomarchive

import (
	_ "embed"
	"testing"
)

//go:embed test.zdb
var testArchive []byte

func TestLoadArchive(t *testing.T) {
	_, err := LoadSocomArchive(testArchive)
	if err != nil {
		t.Fatalf("LoadArchive(testArchive) = _, %v, want _, nil", err)
	}
}
