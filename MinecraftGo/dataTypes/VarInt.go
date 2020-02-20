package dataTypes

import (
	"fmt"
)

func ReadVarInt(buf []byte) (int, int) {
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
