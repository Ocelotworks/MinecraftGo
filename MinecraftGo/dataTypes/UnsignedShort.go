package dataTypes

import (
	"encoding/binary"
)

func ReadUnsignedShort(buf []byte) (interface{}, int) {
	return binary.BigEndian.Uint16(buf), 3
}
