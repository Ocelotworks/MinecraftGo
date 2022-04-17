package dataTypes

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
	"io"
	"math"
)

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
	BlockEntities []ChunkBlockEntity
	//Entities      []ChunkEntity   `nbt:"Entities"`
	Sections   []ChunkSection  `nbt:"sections"`
	Structures ChunkStructures `nbt:"structures"`
	//CarvingMasks  []ChunkCarvingMask `nbt:"CarvingMasks"`

	//Lights List of lists of ?
	//LiquidsToBeTicked List of lists of ?
	//PostProcessing List of lists of ?
	//ToBeTicked List of lists of ?
	//TileEntities List of ?
	//FluidTicks List of ?
	//BlockTicks List of ?
}

// GetBlockAt gets block data for a block at a specific coordinate
func (rc *RegionChunk) GetBlockAt(x int, y int, z int) {
	section := rc.GetSection(y)
	blocksPerEntry := math.Ceil(float64(64 / section.BitsPerBlock))
	blockIndex := float64(x + (z * 16))
	entryNumber := int(math.Ceil(blockIndex / blocksPerEntry))
	bitIndex := int(blockIndex) % int(blocksPerEntry)
	entry := section.BlockStates.Data[entryNumber]
	maskOffset := uint(bitIndex * int(section.BitsPerBlock))
	mask := uint(math.Pow(2, float64(section.BitsPerBlock))-1) << maskOffset
	result := (uint(entry) & mask) >> maskOffset

	fmt.Printf("entry=%d index=%d bpb=%d\n", entryNumber, bitIndex, section.BitsPerBlock)
	fmt.Printf("entry: %64b\n", entry)
	fmt.Printf("mask:  %64b\n", mask)
	//fmt.Printf("upper: %64b\n", upperMask)
	// fmt.Printf("combi: %64b\n", combinedMask)
	fmt.Printf("resul: %64b\n", result)

	block := section.BlockStates.Palette[result]
	fmt.Println("Block name: ", block.Name)

}

func (rc *RegionChunk) GetSection(y int) *ChunkSection {
	for i, section := range rc.Sections {
		// TODO: 16 should be calculated from world height
		if int(int8(section.Y)*16) > y {
			return &rc.Sections[i-1]
		}
	}
	fmt.Println("couldn't find section for ", y)
	return &rc.Sections[0]
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

type ChunkBlockEntity struct {
}

type ChunkSection struct {
	Y            byte                    `nbt:"Y"`
	BlockLight   []byte                  `nbt:"BlockLight"`
	Biomes       ChunkSectionBiome       `nbt:"biomes"`
	BlockStates  ChunkSectionBlockStates `nbt:"block_states"`
	BitsPerBlock byte                    `nbt:"-"`
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
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"BlockStateProperties"`
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

	return &regionChunk, cursor
}
