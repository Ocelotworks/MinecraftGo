package dataTypes

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/Ocelotworks/MinecraftGo/entity"
)

type Chunk struct {
	Length            int
	CompressionScheme byte
	Biomes            []byte
	Sections          []*NetChunkSection
}

type RegionChunk struct {
	Level ChunkLevel `nbt:"Level"`
}

type ChunkLevel struct {
	XPos       int32          `nbt:"xPos"`
	ZPos       int32          `nbt:"zPos"`
	LastUpdate int64          `nbt:"LastUpdate"`
	Biomes     []byte         `nbt:"Biomes"`
	Entities   []ChunkEntity  `nbt:"Entities"`
	Sections   []ChunkSection `nbt:"Sections"`
	HeightMap  []int32        `nbt:"HeightMap"`
}

type ChunkEntity struct {
	OnGround     bool      `nbt:"OnGround"`
	Air          int16     `nbt:"Air"`
	AttackTime   int16     `nbt:"AttackTime"`
	DeathTime    int16     `nbt:"DeathTime"`
	Fire         int16     `nbt:"Fire"`
	Health       int16     `nbt:"Health"`
	FallDistance float32   `nbt:"FallDistance"`
	ID           string    `nbt:"id"`
	Position     []float64 `nbt:"Pos"`
	Rotation     []float32 `nbt:"Rotation"`
}

type ChunkSection struct {
	Y           int8            `nbt:"Y"`
	BlockLight  []byte          `nbt:"BlockLight"`
	Blocks      []byte          `nbt:"Blocks"`
	Data        []byte          `nbt:"Data"`
	SkyLight    []byte          `nbt:"SkyLight"`
	BlockStates []int64         `nbt:"BlockStates"`
	Palette     []*ChunkPalette `nbt:"Palette"`
}

type ChunkPalette struct {
	Name       string                 `nbt:"Name"`
	Properties map[string]interface{} `nbt:"Properties"`
}

func ReadChunk(buf []byte, blockData map[string]entity.BlockData) (interface{}, int) {
	chunk := Chunk{}

	chunkLength, cursor := ReadInt(buf)
	chunk.Length = chunkLength.(int)
	compressionScheme, length := ReadUnsignedByte(buf[4:])
	chunk.CompressionScheme = compressionScheme.(byte)
	cursor += length

	fmt.Println("Chunk Length is ", chunk.Length)
	fmt.Println("Compression Scheme is ", chunk.CompressionScheme)

	if compressionScheme != byte(2) {
		fmt.Println("!!! Invalid compression scheme!")
		return nil, chunk.Length
	}

	reader := bytes.NewReader(buf[cursor:])

	read, exception := zlib.NewReader(reader)

	if exception != nil {
		fmt.Println("Error decompressing chunk: ", exception)
		return chunk, cursor
	}

	var out bytes.Buffer
	io.Copy(&out, read)
	decompressedBytes := out.Bytes()

	//chunk.Raw = decompressedBytes
	chunkData, length := ReadNBT(decompressedBytes)
	cursor += length

	regionChunk := RegionChunk{}
	NBTStructScan(chunkData, &regionChunk)

	if regionChunk.Level.Biomes == nil {
		fmt.Println("No biomes")
	} else {
		chunk.Biomes = regionChunk.Level.Biomes
	}

	if regionChunk.Level.Sections == nil {
		fmt.Println("Chunk has no sections")
		return chunk, cursor
	}

	chunkSections := regionChunk.Level.Sections
	chunk.Sections = make([]*NetChunkSection, len(chunkSections))
	fmt.Println(len(chunkSections), "chunk sections")
	for x, section := range chunkSections {
		if section.BlockStates == nil {
			fmt.Println("!!! Skipping section as it has no blockstates ", x)
			continue
		}
		bitsPerBlock := byte((len(section.BlockStates) * 64) / 4096)
		chunkSection := NetChunkSection{
			BlockCount:   1,
			BitsPerBlock: bitsPerBlock,
			Palette:      []int{},
			DataArray:    section.BlockStates,
		}

		if section.Palette != nil {
			outputPalette := make([]int, len(section.Palette))
			for order, paletteItem := range section.Palette {
				block := blockData[paletteItem.Name]
				if paletteItem.Properties != nil {
					for _, state := range block.States {
						found := true
						for propertyName, property := range paletteItem.Properties {
							if state.Properties[propertyName] != property.(string) {
								found = false
								break
							}
						}
						if found {
							outputPalette[order] = state.ID
							break
						}
					}
				} else {
					outputPalette[order] = block.States[0].ID
				}
			}

			chunkSection.Palette = outputPalette
		} else {
			fmt.Println("Chunk Section does not have a palette!, ", bitsPerBlock)
		}

		chunk.Sections[x] = &chunkSection
	}

	//chunk.Data = chunkData.([]interface{})

	return nil, 0
}
