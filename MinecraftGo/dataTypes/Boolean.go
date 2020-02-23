package dataTypes

func ReadBoolean(buf []byte) (interface{}, int) {
	return buf[0] == 1, 1
}

func WriteBoolean(input interface{}) []byte {
	output := 0x00
	if input.(bool) {
		output = 0x01
	}
	return []byte{byte(output)}
}
