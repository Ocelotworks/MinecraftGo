package dataTypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ReadLong(buf []byte) (interface{}, int) {
	byteBuffer := bytes.NewBuffer(buf[:8])
	output := int64(0)
	exception := binary.Read(byteBuffer, binary.BigEndian, &output)

	if exception != nil {
		fmt.Println("Exception reading long:", exception)
	}

	return output, 9
}

func WriteLong(long interface{}) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(long.(int64)))
	return b
}

func WriteLongArray(input interface{}) []byte {
	arr := input.([]int64)
	output := make([]byte, 0)
	for _, val := range arr {
		output = append(output, WriteLong(val)...)
	}
	return output
}
