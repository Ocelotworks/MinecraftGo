package world

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

type Chunk struct {
	Blocks        [16][16][16]Block
	BlocksInChunk int
}

var sectionHeight = int32(16)
var sectionWidth = int32(16)

func (c *Chunk) SetBlock(x int32, y int32, z int32, block Block) {

	c.Blocks[y%16][z][x] = block
}

func ParseChunkFromRegionChunk(regionChunk *dataTypes.RegionChunk) *Chunk {

	chunk := &Chunk{
		Blocks: [16][16][16]Block{},
	}

	for _, section := range regionChunk.Sections {

		// This chunk section has no blocks
		if section.BlockStates.Data == nil {
			continue
		}

		if len(section.BlockStates.Palette) > 0 {
			for i, block := range section.BlockStates.Palette {
				if block.Name == "minecraft:air" {
					block.Name = "minecraft:stone"
				} else {
					block.Name = "minecraft:dirt"
				}
				section.BlockStates.Palette[i] = block
			}
		}

		// The entire section is made up of a single entry type
		if section.BitsPerBlock < 4 {
			for x := int32(0); x < 16; x++ {
				for y := int32(section.Y); y < (int32(section.Y)*16)+16; y++ {
					for z := int32(0); z < 16; z++ {
						if section.BlockStates.Palette[0].Name != "minecraft:air" {
							chunk.BlocksInChunk++
						}
						chunk.SetBlock(x, y, z, Block{Type: section.BlockStates.Palette[0].Name})
					}
				}
			}
			continue
		}

		// Block states are the direct registry IDs
		if section.BitsPerBlock > 8 {
			for y := int32(0); y < sectionHeight; y++ {
				for z := int32(0); z < sectionWidth; z++ {
					for x := int32(0); x < sectionWidth; x++ {
						chunk.BlocksInChunk++
						chunk.SetBlock(x, y, z, Block{Type: "minecraft:dirt"})
						// TODO: These blocks come directly from the registry ID
					}
				}
			}
			continue
		}
		bitsPerBlock := int32(section.BitsPerBlock)

		individualValueMask := uint((1 << bitsPerBlock) - 1)

		for y := int32(0); y < sectionHeight; y++ {
			for z := int32(0); z < sectionWidth; z++ {
				for x := int32(0); x < sectionWidth; x++ {

					blockNumber := (((y * sectionHeight) + z) * sectionHeight) + x
					startLong := (blockNumber * bitsPerBlock) / 64
					startOffset := (blockNumber * bitsPerBlock) % 64

					data := uint(section.BlockStates.Data[startLong] >> startOffset)

					data &= individualValueMask
					if section.BlockStates.Palette[data].Name != "minecraft:air" {
						chunk.BlocksInChunk++
					}
					//fmt.Printf("Set block x=%d y=%d z=%d to %s \n", x, y, z, section.BlockStates.Palette[data].Name)
					chunk.SetBlock(x, int32(section.Y*16)+y, z, Block{Type: section.BlockStates.Palette[data].Name})
				}
			}
		}
	}

	return chunk
}
