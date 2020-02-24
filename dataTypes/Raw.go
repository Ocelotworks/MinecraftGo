package dataTypes

func ReadRaw(buf []byte) (interface{}, int) {
	return buf, len(buf)
}

func WriteRaw(raw interface{}) []byte {
	return raw.([]byte)
}
