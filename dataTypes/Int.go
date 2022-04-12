package dataTypes

import "encoding/binary"

func ReadInt(buf []byte) (interface{}, int) {
	return int(binary.BigEndian.Uint32(buf[:4])), 4
}

func WriteInt(input interface{}) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(input.(int)))
	return bytes
}

func WriteIntArray(input interface{}) []byte {
	arr := input.([]int)
	output := WriteInt(len(arr))
	for _, val := range arr {
		output = append(output, WriteInt(val)...)
	}
	return output
}

func ReadIntArray(buf []byte) (interface{}, int) {
	length, cursor := ReadInt(buf)
	arr := make([]int, length.(int))
	for i := 0; i < length.(int); i++ {
		arrayItem, cursorDelta := ReadInt(buf[cursor:])
		cursor += cursorDelta
		arr[i] = arrayItem.(int)
	}

	return arr, cursor
}
