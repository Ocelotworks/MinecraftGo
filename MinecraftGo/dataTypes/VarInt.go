package dataTypes

import (
	"fmt"
)

func ReadVarInt(buf []byte) (interface{}, int) {
	numRead := 0
	result := 0

	for numRead, read := range buf {
		value := int32(read & 0b01111111)
		result |= int(value << (7 * numRead))

		if numRead > 5 {
			fmt.Println("Numread overflow")
			return -1, numRead
		}
		if (read & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead + 1
}

func WriteVarInt(value interface{}) []byte {
	output := make([]byte, 0)
	for {
		temp := byte(value.(int) & 0b01111111)
		temp >>= 7
		if value != 0 {
			temp |= 0b10000000
		}
		output = append(output, temp)
		if value == 0 {
			return output
		}
	}
}
