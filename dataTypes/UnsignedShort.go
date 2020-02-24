package dataTypes

import (
	"encoding/binary"
)

func ReadUnsignedShort(buf []byte) (interface{}, int) {
	return binary.BigEndian.Uint16(buf), 2
}

func WriteUnsignedShort(short interface{}) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, short.(uint16))
	return bytes
}
