package dataTypes

func ReadUnsignedByte(buf []byte) (interface{}, int) {
	return buf[0], 2
}

func WriteUnsignedByte(input interface{}) []byte {
	return []byte{input.(byte)}
}
