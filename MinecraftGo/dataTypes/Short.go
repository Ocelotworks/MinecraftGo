package dataTypes

import "encoding/binary"

func ReadShort(buf []byte) (interface{}, int) {
	return int16(binary.BigEndian.Uint16(buf)), 2
}

func WriteShort(short interface{}) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(short.(int16)))
	return bytes
}
