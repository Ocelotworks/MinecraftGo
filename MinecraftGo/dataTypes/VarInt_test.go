package dataTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	test, _ := ReadVarInt([]byte{0xc2, 0x04})
	assert.Equal(t, test, 0)

	test, _ = ReadVarInt([]byte{0x01})
	assert.Equal(t, test, 1)

	test, _ = ReadVarInt([]byte{0x02})
	assert.Equal(t, test, 2)

	test, _ = ReadVarInt([]byte{0x7f})
	assert.Equal(t, test, 127)

	test, _ = ReadVarInt([]byte{0x80, 0x01})
	assert.Equal(t, test, 128)

	test, _ = ReadVarInt([]byte{0xff, 0x01})
	assert.Equal(t, test, 255)

	test, _ = ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x07})
	assert.Equal(t, test, 2147483647)

	test, _ = ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x0f})
	assert.Equal(t, test, -1)

	test, _ = ReadVarInt([]byte{0x80, 0x80, 0x80, 0x80, 0x08})
	assert.Equal(t, test, -2147483648)
}
