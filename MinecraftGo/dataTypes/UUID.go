package dataTypes

func ReadUUID(buf []byte) (interface{}, int) {
	return buf[:16], 17
}

func WriteUUID(input interface{}) []byte {
	//For now lets just treat it as an array of bytes
	return input.([]byte)
}
