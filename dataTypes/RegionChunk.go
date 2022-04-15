package dataTypes

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
	"io"
)

type Chunk struct {
	Length            int
	CompressionScheme byte
	Biomes            []int32
	Sections          []*NetChunkSection
}

type OuterChunk struct {
	Inner RegionChunk `nbt:""`
}

type RegionChunk struct {
	XPos          int32           `nbt:"xPos"`
	YPos          int32           `nbt:"yPos"`
	ZPos          int32           `nbt:"zPos"`
	DataVersion   int32           `nbt:"DataVersion"`
	LastUpdate    int64           `nbt:"LastUpdate"`
	IsLightOn     byte            `nbt:"isLightOn"`
	InhabitedTime int64           `nbt:"InhabitedTime"`
	Status        string          `nbt:"Status"`
	Heightmaps    ChunkHeightmaps `nbt:"Heightmaps"`
	Biomes        []int32         `nbt:"Biomes"`
	Entities      []ChunkEntity   `nbt:"Entities"`
	Sections      []ChunkSection  `nbt:"sections"`
	Structures    ChunkStructures `nbt:"structures"`
	//CarvingMasks  []ChunkCarvingMask `nbt:"CarvingMasks"`

	//Lights List of lists of ?
	//LiquidsToBeTicked List of lists of ?
	//PostProcessing List of lists of ?
	//ToBeTicked List of lists of ?
	//TileEntities List of ?
	//FluidTicks List of ?
	//BlockTicks List of ?
}

type ChunkCarvingMask struct {
	Air    []byte `nbt:"AIR"`
	Liquid []byte `nbt:"LIQUID"`
}

type ChunkStructures struct {
	References StructureReferences `nbt:"References"`
	Starts     StructureStarts     `nbt:"Starts"`
}

type StructureReferences struct {
	Mineshaft []int64 `nbt:"minecraft:mineshaft"`
}

type StructureStarts struct {
	Mineshaft       StructureStart `nbt:"Mineshaft"`
	PillagerOutpost StructureStart `nbt:"Pillager_Outpost"`
	Stronghold      StructureStart `nbt:"Stronghold"`
	Village         StructureStart `nbt:"Village"`
}

type StructureStart struct {
	Id string `nbt:"id"`
}

type ChunkHeightmaps struct {
	MotionBlocking         []int64 `nbt:"MOTION_BLOCKING"`
	MotionBlockingNoLeaves []int64 `nbt:"MOTION_BLOCKING_NO_LEAVES"`
	OceanFloor             []int64 `nbt:"OCEAN_FLOOR"`
	WorldSurface           []int64 `nbt:"WORLD_SURFACE"`
}

type ChunkEntity struct {
	OnGround     byte      `nbt:"OnGround"`
	Air          int16     `nbt:"Air"`
	AttackTime   int16     `nbt:"AttackTime"`
	DeathTime    int16     `nbt:"DeathTime"`
	Fire         int16     `nbt:"Fire"`
	Health       float32   `nbt:"Health"`
	FallDistance float32   `nbt:"FallDistance"`
	ID           string    `nbt:"id"`
	Position     []float64 `nbt:"Pos"`
	Rotation     []float32 `nbt:"Rotation"`
}

type ChunkSection struct {
	Y           uint8                   `nbt:"Y"`
	BlockLight  []byte                  `nbt:"BlockLight"`
	Biomes      ChunkSectionBiome       `nbt:"biomes"`
	BlockStates ChunkSectionBlockStates `nbt:"block_states"`
}

type ChunkSectionBiome struct {
	Palette []string `nbt:"palette"`
	Data    []int64  `nbt:"data"`
}

type ChunkSectionBlockStates struct {
	Palette []ChunkSectionBlockStatePalette `nbt:"palette"`
	Data    []int64                         `nbt:"data"`
}

type ChunkSectionBlockStatePalette struct {
	Name       string               `nbt:"Name"`
	Properties BlockStateProperties `nbt:"BlockStateProperties"`
}

// BlockStateProperties Should be a map but not supported
type BlockStateProperties struct {
	Axis  string `nbt:"axis"`
	Level string `nbt:"level"`
	List  string `nbt:"lit"`
	Snowy string `nbt:"snowy"`
	Half  string `nbt:"half"`
}

func ReadRegionChunk(buf []byte) (*RegionChunk, int) {
	//chunk := Chunk{}

	chunkLength, cursor := ReadInt(buf)
	//chunk.Length = chunkLength.(int)
	compressionScheme, length := ReadUnsignedByte(buf[4:])
	//chunk.CompressionScheme = compressionScheme.(byte)
	cursor += length

	//fmt.Println("Chunk Length is ", chunk.Length)
	//fmt.Println("Compression Scheme is ", chunk.CompressionScheme)

	if compressionScheme != byte(2) {
		fmt.Println("!!! Invalid compression scheme!")
		return nil, chunkLength.(int)
	}

	reader := bytes.NewReader(buf[cursor:])

	read, exception := zlib.NewReader(reader)

	if exception != nil {
		fmt.Println("Error decompressing chunk: ", exception)
		return nil, cursor
	}

	var out bytes.Buffer
	io.Copy(&out, read)
	decompressedBytes := out.Bytes()

	// chunk.Raw = decompressedBytes
	chunkData, chunkDataLength := nbt.ReadNBT(decompressedBytes)
	cursor += chunkDataLength

	outerChunk := OuterChunk{}
	nbt.NBTStructScan(&outerChunk, &chunkData)

	regionChunk := outerChunk.Inner

	//if regionChunk.Level.Biomes == nil {
	//    fmt.Println("No biomes")
	//} else {
	//    chunk.Biomes = regionChunk.Level.Biomes
	//}
	//
	//if regionChunk.Level.Sections == nil {
	//    fmt.Println("Chunk has no sections")
	//    return chunk, cursor
	//}

	//chunkSections := regionChunk.Level.Sections
	//chunk.Sections = make([]*NetChunkSection, len(chunkSections))
	///fmt.Println(len(chunkSections), "chunk sections")
	//for x, section := range chunkSections {
	//    if section.BlockStates == nil {
	//        fmt.Println("!!! Skipping section as it has no blockstates ", x)
	//        continue
	//    }
	//    bitsPerBlock := byte((len(section.BlockStates) * 64) / 4096)
	//    chunkSection := NetChunkSection{
	//        BlockCount:   uint16(len(section.BlockStates)),
	//        BitsPerBlock: bitsPerBlock,
	//        Palette:      []int{},
	//        DataArray:    section.BlockStates,
	//    }
	//
	//    if section.Palette != nil {
	//        outputPalette := make([]int, len(section.Palette))
	//        for order, paletteItem := range section.Palette {
	//            block := blockData[paletteItem.Name]
	//            // TODO: Handle properties again
	//            //if paletteItem.Properties != nil {
	//            //    for _, state := range block.States {
	//            //        found := true
	//            //        for propertyName, property := range paletteItem.Properties {
	//            //            if state.Properties[propertyName] != property.(string) {
	//            //                found = false
	//            //                break
	//            //            }
	//            //        }
	//            //        if found {
	//            //            outputPalette[order] = state.ID
	//            //            break
	//            //        }
	//            //    }
	//            //} else {
	//            outputPalette[order] = block.States[0].ID
	//            //}
	//        }
	//
	//        chunkSection.Palette = outputPalette
	//    } else {
	//        fmt.Println("Chunk Section does not have a palette!, ", bitsPerBlock)
	//    }
	//
	//    chunk.Sections[x] = &chunkSection
	//}

	// chunk.Data = chunkData.([]interface{})

	return &regionChunk, cursor
}
