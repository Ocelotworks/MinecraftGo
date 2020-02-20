package dataTypes

import (
	"encoding/binary"
)

func ReadUnsignedShort(buf []byte) (uint16, int) {
	return binary.BigEndian.Uint16(buf), 3
}
