package dataTypes

import "encoding/binary"

type NetChunkSection struct {
	BlockCount   uint16  //short
	BitsPerBlock byte    //unsigned byte
	Palette      []int   //Palette
	DataArray    []int64 //Array of long
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
	output = append(output, WriteVarInt(len(chunk.DataArray))...)

	for _, long := range chunk.DataArray {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(long))
		output = append(output, b...)
	}

	return output
}

func WriteChunkPalette(p interface{}) []byte {
	output := make([]byte, 0)
	palette := p.([]int)

	output = append(output, WriteVarInt(len(palette))...)
	output = append(output, WriteVarIntArray(palette)...)

	return output
}
