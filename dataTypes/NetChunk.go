package dataTypes

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"math"
)

type NetChunkSection struct {
	BlockCount   uint16  //short
	BitsPerBlock byte    //unsigned byte
	Palette      []int   //Palette
	DataArray    []int64 //Array of long
}

func WriteNetChunk(chunk *RegionChunk, blockData map[string]entity.BlockData) []byte {
	output := make([]byte, 0)
	for _, section := range chunk.Sections {
		paletteSize := len(section.BlockStates.Palette)
		dataSize := len(section.BlockStates.Data)
		bitsPerBlock := byte(math.Ceil(math.Log(float64(paletteSize)) / math.Log(2)))
		// Block count
		if dataSize > 0 {
			output = append(output, WriteVarInt(paletteSize/dataSize)...)
		} else {
			output = append(output, WriteVarInt(0)...)
		}

		output = append(output, bitsPerBlock)
		output = append(output, WriteBlockPalette(section, bitsPerBlock, blockData)...)
		blockStates := WriteVarInt(dataSize)
		//output = append(output, WriteVarInt(dataSize)...)
		for _, varInt := range section.BlockStates.Data {
			blockStates = append(blockStates, WriteVarInt(int(varInt))...)
			//output = append(output, WriteVarInt(int(varInt))...)
		}

		//fmt.Println("==== START BLOCK STATES ====")
		//fmt.Println(hex.Dump(blockStates))
		//fmt.Println("==== END BLOCK STATES ==== ")
		output = append(output, blockStates...)

		// fake biome data for now
		output = append(output, WriteVarInt(0)...) // 0 bits per entry
		output = append(output, WriteVarInt(1)...) // Biome ID 1
		output = append(output, WriteVarInt(0)...) // 0 byte data length
	}

	return output
}

func WriteBlockPalette(section ChunkSection, bitsPerBlock byte, blockData map[string]entity.BlockData) []byte {
	if bitsPerBlock == 0 {
		// Single value mode
		return WriteVarInt(GetBlockIDFromName(section.BlockStates.Palette[0].Name, blockData))
	}

	// Minimum bpb is 4
	if bitsPerBlock < 4 {
		bitsPerBlock = 4
	}

	// Higher than 8 we just use the original ID so there is no palette
	if bitsPerBlock > 8 {
		return []byte{}
	}

	output := WriteVarInt(len(section.BlockStates.Palette))

	for _, block := range section.BlockStates.Palette {
		blockId := GetBlockIDFromName(block.Name, blockData)
		output = append(output, WriteVarInt(blockId)...)
	}

	//fmt.Println("==== START PALETTE ====")
	//fmt.Println(hex.Dump(output))
	//fmt.Println("==== END PALETTE ==== ")

	return output
}

func GetBlockIDFromName(name string, blockData map[string]entity.BlockData) int {
	// For now ignore block states
	block, exists := blockData[name]
	if !exists || len(block.States) == 0 {
		fmt.Println("!! block state doesn't exist for ", name)
		return 0
	}
	return block.States[0].ID
}

//
//func WriteChunk(c interface{}) []byte {
//	chunk := c.([]*NetChunkSection)
//	output := make([]byte, 0)
//	fmt.Println("Chunk is section length", len(chunk))
//	for _, chunkSection := range chunk {
//		if chunkSection == nil {
//			fmt.Println("Chunk section is nil")
//			continue
//		}
//		output = append(output, WriteChunkSection(chunkSection)...)
//	}
//	return output
//}
//
//func WriteChunkSection(c interface{}) []byte {
//	output := make([]byte, 0)
//	chunk := c.(*NetChunkSection)
//
//	fmt.Println("Writing chunk section")
//	output = append(output, WriteUnsignedShort(chunk.BlockCount)...) //Block Count (short)
//	output = append(output, chunk.BitsPerBlock)                      //Bits per block (Byte)
//	if len(chunk.Palette) > 0 {
//		output = append(output, WriteChunkPalette(chunk.Palette)...) //Palette
//	}
//	output = append(output, WriteVarInt(len(chunk.DataArray))...)
//
//	for _, long := range chunk.DataArray {
//		b := make([]byte, 8)
//		binary.BigEndian.PutUint64(b, uint64(long))
//		output = append(output, b...)
//	}
//
//	return output
//}
//
//func WriteChunkPalette(p interface{}) []byte {
//	output := make([]byte, 0)
//	palette := p.([]int)
//
//	output = append(output, WriteVarInt(len(palette))...)
//	output = append(output, WriteVarIntArray(palette)...)
//
//	return output
//}

//
// func ReadBlockArray(bitsPerBlock int, dataArray []int64, palette []int) []block.Block{
// 	output := make([]block.Block, 4096)
// 	for i, block := range dataArray {
//
// 	}
// 	return output
// }
//
