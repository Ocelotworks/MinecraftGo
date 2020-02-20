package dataTypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ReadLong(buf []byte) (int64, int) {
	byteBuffer := bytes.NewBuffer(buf[:8])
	output := int64(0)
	exception := binary.Read(byteBuffer, binary.BigEndian, &output)

	if exception != nil {
		fmt.Println("Exception reading long:", exception)
	}

	return output, 9
}
