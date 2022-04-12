package nbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteByte(t *testing.T) {
	byteVal := NewByte(0x10)
	assert.Equal(t, []byte{0x10}, byteVal.Write())
}

func TestReadByte(t *testing.T) {
	input := []byte{0x10}
	byteVal := Byte{}
	length := byteVal.Read(input)
	assert.Equal(t, 1, length)
	assert.Equal(t, byte(0x10), byteVal.Data)
}

func TestReadWriteByte(t *testing.T) {
	tests := []byte{
		0x00, 0x01, 0x02, 0xFF,
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Byte Test %d", i), func(t *testing.T) {
			input := NewByte(test)
			write := input.Write()
			output := Byte{}
			output.Read(write)
			assert.Equal(t, input, output)
		})
	}
}
