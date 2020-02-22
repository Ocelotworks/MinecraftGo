package dataTypes

type NetChunk struct {
	BlockCount   uint16  //short
	BitsPerBlock byte    //unsigned byte
	Palette      []int   //Palette
	DataArray    []int64 //Array of long
}

func WriteChunk(c interface{}) []byte {
	output := make([]byte, 0)
	chunk := c.(NetChunk)

	output = append(output, WriteUnsignedShort(chunk.BlockCount)...)
	output = append(output, chunk.BitsPerBlock)
	output = append(output, WriteChunkPalette(chunk.Palette)...)
	output = append(output, WriteVarInt(len(chunk.DataArray))...)
	output = append(output, WriteLongArray(chunk.DataArray)...)

	return output
}

func WriteChunkPalette(p interface{}) []byte {
	output := make([]byte, 0)
	palette := p.([]int)

	output = append(output, WriteVarInt(len(palette))...)
	output = append(output, WriteVarIntArray(palette)...)

	return output
}
