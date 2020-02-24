package dataTypes

import "encoding/binary"

func ReadInt(buf []byte) (interface{}, int) {
	slice := buf[:4]

	return int(binary.BigEndian.Uint32(slice)), 4
}

func WriteInt(input interface{}) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(input.(int)))
	return bytes
}

func WriteIntArray(input interface{}) []byte {
	arr := input.([]int)
	output := make([]byte, 0)
	for _, val := range arr {
		output = append(output, WriteInt(val)...)
	}
	return output
}
