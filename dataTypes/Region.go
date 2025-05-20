package dataTypes

import (
	"encoding/binary"
	"fmt"
	"math"
)

type RegionMetadata struct {
	ChunkOffsets []*ChunkOffset
	Chunks       []*RegionChunk
}

func (rm *RegionMetadata) GetChunk(cx int32, cz int32) *RegionChunk {
	for _, chunk := range rm.Chunks {
		if chunk != nil && chunk.XPos == cx && chunk.ZPos == cz {
			return chunk
		}
	}
	fmt.Println("Unable to find chunk for ", cx, cz)
	return nil
}

type ChunkOffset struct {
	Offset    int
	Size      int
	Timestamp int
}

type Region struct {
	Chunks []RegionChunk `nbt:"Chunk*"`
}

func ReadRegionFile(buf []byte) RegionMetadata {
	region := RegionMetadata{
		ChunkOffsets: make([]*ChunkOffset, 1024),
	}

	cursor := 0
	for i := 0; i < 1024; i++ {
		chunkOffset := ChunkOffset{}
		entry := int(binary.BigEndian.Uint32(buf[cursor : cursor+4]))
		chunkOffset.Offset = ((entry >> 8) & 0xFFFFFF) * 4096
		chunkOffset.Size = (entry & 0xFF) * 4096
		//fmt.Println("Size ", chunkOffset.Size)
		cursor += 4
		region.ChunkOffsets[i] = &chunkOffset
	}

	for i := 0; i < 1024; i++ {
		region.ChunkOffsets[i].Timestamp = int(binary.BigEndian.Uint32(buf[cursor : cursor+4]))
		cursor += 4
		//fmt.Println(region.ChunkOffsets[i])
	}

	// fmt.Println("Cursor: ", cursor)

	region.Chunks = make([]*RegionChunk, 1024)
	for i := 0; i < 1024; i++ {
		offset := region.ChunkOffsets[i].Offset
		size := region.ChunkOffsets[i].Size
		if size == 0 {
			continue
		}
		//fmt.Println("Offset is ", offset)

		chunk, length := ReadRegionChunk(buf[offset:])
		for i := range chunk.Sections {
			chunk.Sections[i].BitsPerBlock = byte(math.Ceil(math.Log(float64(len(chunk.Sections[i].BlockStates.Palette))) / math.Log(2)))
			if chunk.Sections[i].BitsPerBlock < 4 {
				chunk.Sections[i].BitsPerBlock = 4
			}
		}
		region.Chunks[i] = chunk
		cursor += length
	}

	//chunk, _ := ReadChunk(buf[cursor:])

	//region.Chunks[0] = chunk.(Chunk)
	return region
}
