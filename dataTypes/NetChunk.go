package dataTypes

type NetChunkSection struct {
	BlockCount   uint16 //short
	BitsPerBlock byte   //unsigned byte
	Palette      []int  //Palette
	DataArray    []byte //Array of long
}

func WriteChunk(c interface{}) []byte {
	chunk := c.([]NetChunkSection)
	output := make([]byte, 0)
	for _, chunkSection := range chunk {
		output = append(output, WriteChunkSection(chunkSection)...)
	}
	return output
}

func WriteChunkSection(c interface{}) []byte {
	output := make([]byte, 0)
	chunk := c.(NetChunkSection)

	output = append(output, WriteUnsignedShort(chunk.BlockCount)...) //Block Count (short)
	output = append(output, chunk.BitsPerBlock)                      //Bits per block (Byte)
	if len(chunk.Palette) > 0 {
		output = append(output, WriteChunkPalette(chunk.Palette)...) //Palette
	}
	output = append(output, WriteVarInt(len(chunk.DataArray)/8)...)
	output = append(output, chunk.DataArray...)

	return output
}

func WriteChunkPalette(p interface{}) []byte {
	output := make([]byte, 0)
	palette := p.([]int)

	output = append(output, WriteVarInt(len(palette))...)
	output = append(output, WriteVarIntArray(palette)...)

	return output
}
