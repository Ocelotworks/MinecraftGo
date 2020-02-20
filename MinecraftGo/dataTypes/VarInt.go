package dataTypes

import (
	"encoding/hex"
	"fmt"
)

func ReadVarInt(buf []byte) (interface{}, int) {
	numRead := 0
	result := 0

	var read byte

	for numRead, read = range buf {
		value := int32(read & 0b01111111)
		result |= int(value << (7 * numRead))

		fmt.Println(hex.EncodeToString([]byte{read}))

		if numRead > 5 {
			fmt.Println("Numread overflow")
			return -1, numRead
		}
		if (read & 0b10000000) == 0 {
			break
		}
	}

	fmt.Println("Read finished: ", result, numRead+1)

	return result, numRead + 1
}

func WriteVarInt(value interface{}) []byte {
	intValue := value.(int)
	output := make([]byte, 0)
	for {
		temp := byte(intValue & 0b01111111)
		intValue >>= 7
		if intValue != 0 {
			temp |= 0b10000000
		}
		output = append(output, temp)
		if intValue == 0 {
			return output
		}
	}
}
