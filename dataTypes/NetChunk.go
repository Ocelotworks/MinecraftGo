package dataTypes

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"reflect"
)

// WriteNetChunk converts a RegionChunk to data compatible with the packet.ChunkData packet
func WriteNetChunk(chunk *RegionChunk, blockData map[string]entity.BlockData) []byte {
	//fmt.Printf("sections=%d\n", len(chunk.Sections))
	output := make([]byte, 0)
	for _, section := range chunk.Sections {
		//preLen := len(output)
		paletteSize := len(section.BlockStates.Palette)
		dataSize := len(section.BlockStates.Data)
		bitsPerBlock := section.BitsPerBlock

		// Block count
		if dataSize > 0 {
			output = append(output, WriteShort(int16((dataSize*64)/paletteSize))...)
		} else {
			// fmt.Println("Empty! chunk")
			output = append(output, WriteShort(int16(0))...)
		}

		output = append(output, WriteBlockPalette(section, bitsPerBlock, blockData)...)
		output = append(output, WriteVarInt(dataSize)...)
		for _, varInt := range section.BlockStates.Data {
			output = append(output, WriteLong(varInt)...)
		}

		// fake biome data for now
		output = append(output, 0)                 // bits per entry
		output = append(output, WriteVarInt(0)...) // Biome
		output = append(output, WriteVarInt(0)...) // data length
		//postLen := len(output)
		//fmt.Printf("bpb=%d paletteSize=%d dataSize=%d preLen=%d postLen=%d sectionSize=%d\n", bitsPerBlock, paletteSize, dataSize, preLen, postLen, postLen-preLen)
		//output = append(output, WriteLongArray(make([]int64, 6))...)
	}

	//fmt.Println(hex.Dump(output))

	return output
}

func WriteBlockPalette(section ChunkSection, bitsPerBlock byte, blockData map[string]entity.BlockData) []byte {
	if bitsPerBlock == 0 {
		// Single value mode
		if len(section.BlockStates.Palette) == 0 {
			return append(WriteShort(int16(0)), WriteVarInt(0)...)
		}
		return append(WriteShort(int16(0)), WriteVarInt(GetBlockID(section.BlockStates.Palette[0], blockData))...)
	}

	// Minimum bpb is 4
	if bitsPerBlock < 4 {
		bitsPerBlock = 4
	}

	// Higher than 8 we just use the original ID so there is no palette
	if bitsPerBlock > 8 {
		return WriteShort(int16(15))
	}

	// [bitsPerBlock] [paletteLength] [palette...]

	output := append([]byte{bitsPerBlock}, WriteVarInt(len(section.BlockStates.Palette))...)

	for _, block := range section.BlockStates.Palette {
		blockId := GetBlockID(block, blockData)
		//fmt.Printf("name=%s id=%d index=%d\n", block.Name, blockId, i)
		output = append(output, WriteVarInt(blockId)...)
	}

	//fmt.Printf("bitsPerBlock=%d paletteLength=%d\n", bitsPerBlock, len(section.BlockStates.Palette))

	//fmt.Println("==== START PALETTE ====")
	//fmt.Println(hex.Dump(output))
	//fmt.Println("==== END PALETTE ==== ")

	return output
}

func GetBlockID(target ChunkSectionBlockStatePalette, blockData map[string]entity.BlockData) int {
	// For now ignore block states
	block, exists := blockData[target.Name]
	if !exists || len(block.States) == 0 {
		fmt.Println("!! block state doesn't exist for ", target.Name)
		return 0
	}
	for _, blockState := range block.States {
		if reflect.DeepEqual(blockState.Properties, target.Properties) {
			return blockState.ID
		}
	}
	fmt.Println("Couldn't find applicable block state ", target.Properties)
	return block.States[0].ID
}
