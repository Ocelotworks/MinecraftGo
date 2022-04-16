package dataTypes

func ReadRaw(buf []byte) (interface{}, int) {
	return buf, len(buf)
}

func WriteRaw(raw interface{}) []byte {
	return raw.([]byte)
}

func ReadByte(buf []byte) (interface{}, int) {
	return buf[0], 1
}

func WriteByte(raw interface{}) []byte {
	return []byte{raw.(byte)}
}
