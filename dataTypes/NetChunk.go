package dataTypes

import (
	"encoding/hex"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"log"
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
			fmt.Println("dataSize > 0, + ", int16((dataSize*64)/paletteSize))
			output = append(output, WriteShort(int16((dataSize*64)/paletteSize))...)
		} else {
			fmt.Println("dataSize <= 0, ", WriteShort(int16(0)))
			// fmt.Println("Empty! chunk")
			output = append(output, WriteShort(int16(0))...)
		}

		// Block Palette
		output = append(output, WriteBlockPalette(section, bitsPerBlock, blockData)...)

		// Block Data
		for _, varInt := range section.BlockStates.Data {
			output = append(output, WriteLong(varInt)...)
		}

		output = append(output, WriteBiomePalette(section, bitsPerBlock, blockData)...)

		for _, varInt := range section.Biomes.Data {
			output = append(output, WriteLong(varInt)...)
		}

		fmt.Println("Data size ", 0)
		fmt.Println("pallette size", 0)
		fmt.Println("Section", hex.Dump(output))
	}

	return output
}

func WriteBiomePalette(section ChunkSection, bitsPerBlock byte, blockData map[string]entity.BlockData) []byte {

	if bitsPerBlock == 0 {
		if len(section.Biomes.Palette) == 0 {
			return append(WriteByte(byte(0)), WriteVarInt(0)...)
		}

		return append(WriteByte(byte(0)), WriteVarInt(0)...) // TODO: actual biome data
	}

	if bitsPerBlock < 1 {
		bitsPerBlock = 1
	}

	if bitsPerBlock > 3 {
		return []byte{6}
	}

	output := append([]byte{bitsPerBlock}, WriteVarInt(len(section.Biomes.Palette))...)
	for _, _ = range section.Biomes.Palette {
		blockId := 0 // TODO ..
		//fmt.Printf("name=%s id=%d index=%d\n", block.Name, blockId, i)
		output = append(output, WriteVarInt(blockId)...)
	}

	return output

}

func WriteBlockPalette(section ChunkSection, bitsPerBlock byte, blockData map[string]entity.BlockData) []byte {
	// Single value mode
	if bitsPerBlock == 0 {

		// All air
		if len(section.BlockStates.Palette) == 0 {
			if len(section.BlockStates.Data) != 0 {
				panic("palette cannot be 0 whilst block data has data")
			}

			fmt.Println("palette and block data = 0, +", WriteVarInt(0))
			return append(WriteByte(byte(0)), WriteVarInt(0)...)
		}

		fmt.Println("All single block type +", 0, WriteVarInt(GetBlockID(section.BlockStates.Palette[0], blockData)))
		// All a single other block type
		return append(WriteByte(byte(0)), WriteVarInt(GetBlockID(section.BlockStates.Palette[0], blockData))...)
	}

	// Minimum bpb is 4
	if bitsPerBlock < 4 {
		bitsPerBlock = 4
	}

	// Higher than 8 we just use the original ID so there is no palette
	if bitsPerBlock > 8 {
		fmt.Println("Bits per block are over 8 so we just return 15")
		return []byte{15}
	}

	// [bitsPerBlock] [paletteLength] [palette...]

	fmt.Println("Palette length", WriteVarInt(len(section.BlockStates.Palette)))
	output := append([]byte{bitsPerBlock}, WriteVarInt(len(section.BlockStates.Palette))...)
	for _, block := range section.BlockStates.Palette {
		blockId := GetBlockID(block, blockData)
		//fmt.Printf("name=%s id=%d index=%d\n", block.Name, blockId, i)
		fmt.Println("block id", WriteVarInt(blockId))
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
		log.Panicf("!! block state doesn't exist for %s", target.Name)
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
