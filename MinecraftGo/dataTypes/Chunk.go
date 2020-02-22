package dataTypes

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
)

type Chunk struct {
	Length            int
	CompressionScheme byte
	Data              []interface{}
	Raw               []byte
}

func ReadChunk(buf []byte) (interface{}, int) {
	chunk := Chunk{}

	chunkLength, cursor := ReadInt(buf)
	chunk.Length = chunkLength.(int)
	compressionScheme, length := ReadUnsignedByte(buf[4:])
	chunk.CompressionScheme = compressionScheme.(byte)

	fmt.Println("Chunk Length is ", chunk.Length)
	fmt.Println("Compression Scheme is ", chunk.CompressionScheme)

	reader := bytes.NewReader(buf[cursor:])

	read, exception := zlib.NewReader(reader)

	if exception != nil {
		fmt.Println("Error decompressing chunk: ", exception)
		return chunk, cursor
	}

	fmt.Println(reader.Size(), reader.Len())
	decompressedBytes := make([]byte, reader.Size())
	_, exception = read.Read(decompressedBytes)

	if exception != nil {
		fmt.Println("Error writing chunk into slice: ", exception)
		return chunk, cursor
	}

	fmt.Println("Decompress successful")
	fmt.Println(hex.Dump(decompressedBytes))

	chunk.Raw = decompressedBytes
	chunkData, length := ReadNBT(decompressedBytes)
	cursor += length

	fmt.Println(chunkData)
	chunk.Data = chunkData.([]interface{})

	return chunk, cursor
}
