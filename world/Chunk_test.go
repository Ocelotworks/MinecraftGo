package world

import (
	"os"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestChunk(t *testing.T) *Chunk {
	t.Helper()

	inData, err := os.ReadFile("../data/worlds/MCGO_FlatTest/region/r.0.0.mca")
	require.NoError(t, err)

	region := dataTypes.ReadRegionFile(inData)

	chunk, err := ParseChunkFromRegionChunk(region.Chunks[0], Overworld)
	require.NoError(t, err)

	return chunk
}

func TestParseChunkFromRegionChunk(t *testing.T) {
	chunk := loadTestChunk(t)

	// The superflat test world is bedrock, two dirt, then grass, then air.
	assert.Equal(t, "minecraft:bedrock", chunk.GetBlock(0, -64, 0).Type)
	assert.Equal(t, "minecraft:dirt", chunk.GetBlock(0, -63, 0).Type)
	assert.Equal(t, "minecraft:dirt", chunk.GetBlock(0, -62, 0).Type)
	assert.Equal(t, "minecraft:grass_block", chunk.GetBlock(0, -61, 0).Type)
	assert.Equal(t, "minecraft:air", chunk.GetBlock(0, -60, 0).Type)

	// Block state properties survive parsing.
	assert.Equal(t, "false", chunk.GetBlock(0, -61, 0).Properties["snowy"])
}

func TestParseChunkCoversWholeWorldHeight(t *testing.T) {
	chunk := loadTestChunk(t)

	// The client reads exactly one section per level of the dimension, so
	// every section has to be present even where the region file has none.
	assert.Len(t, chunk.Sections, 24)
	assert.Equal(t, -4, chunk.Sections[0].Y)
	assert.Equal(t, 19, chunk.Sections[23].Y)

	// Sections the region file omitted, or stored as a single-entry palette,
	// still read as air rather than as an unset block.
	for y := 0; y < 320; y += 16 {
		assert.True(t, chunk.GetBlock(0, y, 0).IsAir(), "expected air at y=%d", y)
	}
}

// Section Y is a signed byte in the region file; a section at or above Y=8
// overflows if it is multiplied by 16 before being widened.
func TestParseChunkHandlesHighSections(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")

	for _, y := range []int{-64, 0, 128, 200, 300, 319} {
		chunk.SetBlock(1, y, 2, Block{Type: "minecraft:stone"})
	}

	for _, y := range []int{-64, 0, 128, 200, 300, 319} {
		assert.Equal(t, "minecraft:stone", chunk.GetBlock(1, y, 2).Type, "at y=%d", y)
	}

	// A block written high up must not have landed anywhere else.
	assert.True(t, chunk.GetBlock(1, 129, 2).IsAir())
	assert.True(t, chunk.GetBlock(1, -48, 2).IsAir())
}

func TestSetBlockOutsideWorldHeight(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")

	// Below the world and above it: both must be ignored, not wrap into a
	// section or panic on a negative index.
	assert.NotPanics(t, func() {
		chunk.SetBlock(0, -70, 0, Block{Type: "minecraft:stone"})
		chunk.SetBlock(0, -65, 0, Block{Type: "minecraft:stone"})
		chunk.SetBlock(0, 320, 0, Block{Type: "minecraft:stone"})
	})
	assert.True(t, chunk.GetBlock(0, -64, 0).IsAir())
}

func TestBlockCounts(t *testing.T) {
	chunk := loadTestChunk(t)

	// Four solid layers of 16x16 in the bottom section, nothing above.
	assert.Equal(t, 4*16*16, chunk.Sections[0].NonEmptyBlockCount)
	assert.Equal(t, 0, chunk.Sections[0].FluidCount)
	for _, section := range chunk.Sections[1:] {
		assert.Equal(t, 0, section.NonEmptyBlockCount)
	}
}

func TestFluidCount(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")

	chunk.SetBlock(0, 0, 0, Block{Type: "minecraft:water", Properties: map[string]string{"level": "0"}})
	chunk.SetBlock(1, 0, 0, Block{Type: "minecraft:oak_slab", Properties: map[string]string{"waterlogged": "true"}})
	chunk.SetBlock(2, 0, 0, Block{Type: "minecraft:stone"})
	chunk.SetBlock(3, 0, 0, Block{Type: "minecraft:cave_air"})
	chunk.RecalcBlockCounts()

	section := chunk.Sections[Overworld.SectionIndex(0)]
	assert.Equal(t, 3, section.NonEmptyBlockCount)
	assert.Equal(t, 2, section.FluidCount)
}

func TestCeilLog2(t *testing.T) {
	// Mth.ceillog2: the bits needed to index that many distinct values.
	for value, expected := range map[int]int{0: 0, 1: 0, 2: 1, 3: 2, 4: 2, 5: 3, 8: 3, 9: 4, 16: 4, 17: 5, 256: 8, 257: 9} {
		assert.Equal(t, expected, ceilLog2(value), "ceilLog2(%d)", value)
	}
}

func TestUnpackStorageRejectsWrongLength(t *testing.T) {
	// SimpleBitStorage refuses a data array that is not exactly the length its
	// bit width implies, rather than reading off the end of it.
	_, err := unpackStorage(make([]int64, 100), BlocksPerSection, 4, 4)
	assert.Error(t, err)

	_, err = unpackStorage(make([]int64, 256), BlocksPerSection, 4, 4)
	assert.NoError(t, err)
}

func TestUnpackStorageRejectsMissingPaletteEntry(t *testing.T) {
	data := make([]int64, 256)
	data[0] = 0x7 // index 7, but the palette only has 2 entries

	_, err := unpackStorage(data, BlocksPerSection, 4, 2)
	assert.Error(t, err)
}

func TestUnpackStorageOrdersValuesLowBitsFirst(t *testing.T) {
	data := make([]int64, 256)
	// Three 4-bit values in the first long: 1, 2, 3.
	data[0] = 1 | 2<<4 | 3<<8

	entries, err := unpackStorage(data, BlocksPerSection, 4, 4)
	require.NoError(t, err)

	assert.Equal(t, []uint16{1, 2, 3, 0}, entries[:4])
}
