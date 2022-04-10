package dataTypes

func ReadBitSet(buf []byte) (interface{}, int) {
	size, cursor := ReadVarInt(buf)

	output := make([]int64, size.(int))

	for i := 0; i < size.(int); i++ {
		longValue, cursorDelta := ReadLong(buf[cursor:])
		output[i] = longValue.(int64)
		cursor += cursorDelta
	}

	return output, cursor
}

func WriteBitSet(input interface{}) []byte {
	castInput := input.([]int64)
	// Write the long length
	output := WriteVarInt(len(castInput))

	for i := 0; i < len(castInput); i++ {
		output = append(output, WriteLong(castInput[i])...)
	}

	return output
}
