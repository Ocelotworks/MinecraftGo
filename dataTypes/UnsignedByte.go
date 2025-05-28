package dataTypes

func ReadUnsignedByte(buf []byte) (interface{}, int) {
	return buf[0], 1
}

func WriteUnsignedByte(input interface{}) []byte {
	return []byte{input.(byte)}
}

func ReadVarIntByteArray(buf []byte) (interface{}, int) {
	arrayLength, cursor := ReadVarInt(buf)
	return buf[cursor : arrayLength.(int)+cursor], arrayLength.(int) + cursor
}

func WriteVarIntByteArray(input interface{}) []byte {
	byteArray := input.([]byte)
	b := WriteVarInt(len(byteArray))
	b = append(b, byteArray...)
	return b
}
