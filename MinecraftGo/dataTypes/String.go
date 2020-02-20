package dataTypes

import (
	"bytes"
)

func ReadString(buf []byte) (string, int) {
	stringLength, stringStart := ReadVarInt(buf)

	output := string(bytes.Runes(buf[stringStart:]))[:stringLength]

	return output, stringStart + len([]byte(output))

}
