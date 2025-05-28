package dataTypes

func ReadPrefixedOptionalByteArray(buf []byte) (interface{}, int) {
	present, n := ReadBoolean(buf)
	if !(present.(bool)) {
		return nil, n
	}

	return buf[n:], n
}

func WritePrefixedOptionalByteArray(input interface{}) []byte {
	byteArray := input.([]byte)

	if byteArray == nil || len(byteArray) == 0 {
		return WriteBoolean(false)
	}

	var b []byte
	b = WriteBoolean(true)
	b = append(b, byteArray...)
	return b
}
