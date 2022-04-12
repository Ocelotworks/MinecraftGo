package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteByteArray(t *testing.T) {
	byteVal := NewByteArray([]byte{0x01, 0x02, 0x03, 0x04})
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x04, 0x01, 0x02, 0x03, 0x04}, byteVal.Write())
}

func TestReadByteArray(t *testing.T) {
	input := []byte{0x00, 0x00, 0x00, 0x04, 0x01, 0x02, 0x03, 0x04}
	byteVal := ByteArray{}
	length := byteVal.Read(input)
	assert.Equal(t, 8, length)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, byteVal.Data)
}

func TestReadWriteByteArray(t *testing.T) {
	tests := [][]byte{
		{0x00, 0x00, 0x01, 0x02},
		{},
		{0x00},
		{0xFF, 0x00, 0xFF},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("ByteArray Test %d", i), func(t *testing.T) {
			input := NewByteArray(test)
			write := input.Write()
			output := ByteArray{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
