package fstream

import "testing"

func FuzzStoreChunk(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
	})
}
