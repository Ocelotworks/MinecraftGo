package dataTypes

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

// bitsPerBlockForPalette is the width a block state container is stored at on
// disc, following Strategy.createForBlockStates: the number of bits needed to
// index the palette (Mth.ceillog2), raised to a floor of 4.  A single-entry
// palette gets 0 — such a section carries no data array at all.
//
// The bit count is computed exactly rather than through math.Log, which can
// round a power of two up and silently widen the container by one bit.
func bitsPerBlockForPalette(paletteSize int) int {
	if paletteSize <= 1 {
		return 0
	}

	stateBits := bits.Len32(uint32(paletteSize - 1))
	if stateBits < 4 {
		return 4
	}
	return stateBits
}

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
		if chunk == nil {
			// Unreadable chunk (an unsupported compression scheme, say).
			// Leave the slot nil rather than dereferencing it below.
			cursor += length
			continue
		}
		for i := range chunk.Sections {
			chunk.Sections[i].BitsPerBlock = byte(bitsPerBlockForPalette(len(chunk.Sections[i].BlockStates.Palette)))
		}
		region.Chunks[i] = chunk
		cursor += length
	}

	//chunk, _ := ReadChunk(buf[cursor:])

	//region.Chunks[0] = chunk.(Chunk)
	return region
}
