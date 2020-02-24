package dataTypes

import (
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestFixedValues(t *testing.T) {
	//test, _ := ReadVarInt([]byte{0x00})
	//assert.Equal(t, test, 0)
	//
	//test, _ = ReadVarInt([]byte{0x01})
	//assert.Equal(t, test, 1)
	//
	//test, _ = ReadVarInt([]byte{0x02})
	//assert.Equal(t, test, 2)
	//
	//test, _ = ReadVarInt([]byte{0x7f})
	//assert.Equal(t, test, 127)
	//
	//test, _ = ReadVarInt([]byte{0x80, 0x01})
	//assert.Equal(t, test, 128)
	//
	//test, _ = ReadVarInt([]byte{0xff, 0x01})
	//assert.Equal(t, test, 255)
	//
	//test, _ = ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x07})
	//assert.Equal(t, test, 2147483647)
	//
	//test, _ = ReadVarInt([]byte{0xff, 0xff, 0xff, 0xff, 0x0f})
	//assert.Equal(t, test, -1)
	//
	//test, _ = ReadVarInt([]byte{0x80, 0x80, 0x80, 0x80, 0x08})
	//assert.Equal(t, test, -2147483648)
}

func TestReadWrite(t *testing.T) {
	for i := 0; i <= 2048; i++ {
		//test, _ := ReadVarInt(WriteVarInt(i))
		//	assert.Equal(t, i, test)
	}
}
