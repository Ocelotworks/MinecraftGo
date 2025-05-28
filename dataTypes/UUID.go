package dataTypes

func ReadUUID(buf []byte) (interface{}, int) {
	return buf[:16], 17
}

func WriteUUID(input interface{}) []byte {
	//For now lets just treat it as an array of bytes
	return input.([]byte)
}

func ReadPrefixedOptionalUUID(buf []byte) (interface{}, int) {
	present, n := ReadBoolean(buf)
	if !(present.(bool)) {
		return nil, n
	}

	return ReadUUID(buf[n:])
}

func WritePrefixedOptionalUUID(input interface{}) []byte {
	uuid := input.([]byte)

	if uuid == nil || len(uuid) == 0 {
		return WriteBoolean(false)
	}

	var b []byte
	b = WriteBoolean(true)
	b = append(b, uuid...)
	return b
}
