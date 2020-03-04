package dataTypes

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"io"
	"strconv"
)

type Chunk struct {
	Length            int
	CompressionScheme byte
	Biomes            []byte
	Sections          []*NetChunkSection
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

	level := chunkData.([]interface{})[0].(NBTNamed).Data.([]interface{})[0].(NBTNamed).Data

	asMap := NBTAsMap(level.([]interface{})).(map[string]interface{})

	if asMap["Biomes"] == nil {
		fmt.Println("No biomes")
	} else {
		chunk.Biomes = asMap["Biomes"].([]byte)
	}

	if asMap["Sections"] == nil {
		fmt.Println("Chunk has no sections")
		return chunk, cursor
	}

	chunkSections := asMap["Sections"].(map[string]interface{})["List-0"].(map[string]interface{})
	chunk.Sections = make([]*NetChunkSection, len(chunkSections))
	fmt.Println(len(chunkSections), "chunk sections")
	for x, section := range chunkSections {
		rawChunkSection := section.(map[string]interface{})
		if rawChunkSection["BlockStates"] == nil {
			fmt.Println("!!! Skipping section as it has no blockstates ", x)
			continue
		}

		blockStates := rawChunkSection["BlockStates"].([]int64)
		bitsPerBlock := byte((len(blockStates) * 64) / 4096)
		chunkSection := NetChunkSection{
			BlockCount:   1,
			BitsPerBlock: bitsPerBlock,
			Palette:      []int{},
			DataArray:    blockStates,
		}

		if rawChunkSection["Palette"] != nil {
			rawPalette := rawChunkSection["Palette"].(map[string]interface{})["List-0"].(map[string]interface{})

			outputPalette := make([]int, len(rawPalette))
			i := 0
			for x, rawPaletteItem := range rawPalette {
				//God this is horrible but it will do for now
				orderStr := bytes.SplitN([]byte(x), []byte("_"), 2)[1]
				order, _ := strconv.Atoi(string(orderStr))

				fmt.Println(orderStr)

				paletteItem := rawPaletteItem.(map[string]interface{})
				block := blockData[paletteItem["Name"].(string)]
				if paletteItem["Properties"] != nil {
					blockProperties := paletteItem["Properties"].(map[string]interface{})["Compound_0"].(map[string]interface{})
					for _, state := range block.States {
						found := true
						for propertyName, property := range blockProperties {
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
				i++
			}

			chunkSection.Palette = outputPalette
		} else {
			fmt.Println("Chunk Section does not have a palette!, ", bitsPerBlock)
		}

		orderStr := bytes.SplitN([]byte(x), []byte("_"), 2)[1]
		order, _ := strconv.Atoi(string(orderStr))
		chunk.Sections[order] = &chunkSection
	}

	//chunk.Data = chunkData.([]interface{})

	return chunk, cursor
}
