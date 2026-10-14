package dataTypes

import (
	"bytes"
	"reflect"
	"testing"
)

// The encodings are what java.util.BitSet.toByteArray() produces, prefixed with its length.
func TestWriteBitSet(t *testing.T) {
	cases := []struct {
		name  string
		input []int64
		want  []byte
	}{
		{"empty", []int64{}, []byte{0x00}},
		{"only zero longs", []int64{0, 0}, []byte{0x00}},
		{"bit 0", []int64{1}, []byte{0x01, 0x01}},
		{"low 26 bits", []int64{1<<26 - 1}, []byte{0x04, 0xFF, 0xFF, 0xFF, 0x03}},
		{"bit 8 needs a second byte", []int64{1 << 8}, []byte{0x02, 0x00, 0x01}},
		{"bit 64 spills into the next long", []int64{0, 1}, []byte{0x09, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}},
		{"high bit of a long", []int64{-1 << 63}, []byte{0x08, 0, 0, 0, 0, 0, 0, 0, 0x80}},
	}

	for _, c := range cases {
		if got := WriteBitSet(c.input); !bytes.Equal(got, c.want) {
			t.Errorf("%s: got % x, want % x", c.name, got, c.want)
		}
	}
}

func TestReadBitSetRoundTrip(t *testing.T) {
	for _, input := range [][]int64{{}, {1}, {1<<26 - 1}, {1 << 8}, {0, 1}, {-1 << 63}, {-1, 5}} {
		encoded := append(WriteBitSet(input), 0xAB) // a trailing byte the reader must leave alone

		got, cursor := ReadBitSet(encoded)
		if cursor != len(encoded)-1 {
			t.Errorf("%v: consumed %d of %d bytes", input, cursor, len(encoded)-1)
		}
		if len(got.([]int64)) == 0 && len(input) == 0 {
			continue
		}
		// trailing zero longs are not preserved, which BitSet does not distinguish either
		want := input
		for len(want) > 0 && want[len(want)-1] == 0 {
			want = want[:len(want)-1]
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("round trip of %v gave %v", input, got)
		}
	}
}
