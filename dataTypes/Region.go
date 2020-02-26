package dataTypes

import (
	"encoding/binary"
	"fmt"
)

type Region struct {
	ChunkOffsets []ChunkOffset
	Chunks       []Chunk
}

type ChunkOffset struct {
	Offset    []byte
	Size      byte
	Timestamp int
}

func ReadRegionFile(buf []byte) Region {
	region := Region{
		ChunkOffsets: make([]ChunkOffset, 1024),
	}

	cursor := 0
	for i := 0; i < 1024; i++ {
		chunkOffset := ChunkOffset{}
		chunkOffset.Offset = buf[cursor : cursor+3]
		chunkOffset.Size = buf[cursor+4]
		fmt.Println("Size ", chunkOffset.Size)
		cursor += 4
		region.ChunkOffsets[i] = chunkOffset
	}

	for i := 0; i < 1024; i++ {
		region.ChunkOffsets[i].Timestamp = int(binary.BigEndian.Uint32(buf[cursor : cursor+4]))
		cursor += 4
		fmt.Println(region.ChunkOffsets[i])
	}

	fmt.Println("Cursor: ", cursor)

	region.Chunks = make([]Chunk, 1024)
	for i := 0; i < 1024; i++ {
		chunk, length := ReadChunk(buf[cursor:])
		region.Chunks[i] = chunk.(Chunk)
		cursor += length
	}

	return region
}
