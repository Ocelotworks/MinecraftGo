package dataTypes

import (
	"encoding/binary"
	"math"
)

func ReadFloat(buf []byte) (interface{}, int) {
	in, count := ReadInt(buf)
	return math.Float32frombits(uint32(in.(int))), count
}

func WriteFloat(input interface{}) []byte {
	return WriteInt(int(math.Float32bits(input.(float32))))
}

func ReadDouble(buf []byte) (interface{}, int) {
	slice := buf[:8]
	return math.Float64frombits(binary.BigEndian.Uint64(slice)), 9
}

func WriteDouble(input interface{}) []byte {
	return WriteInt(int(math.Float64bits(input.(float64))))
}
