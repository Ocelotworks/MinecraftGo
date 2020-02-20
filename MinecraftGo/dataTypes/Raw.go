package dataTypes

func WriteRaw(raw interface{}) []byte {
	return raw.([]byte)
}
