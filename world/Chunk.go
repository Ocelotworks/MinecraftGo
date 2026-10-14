package world

import (
	"fmt"
	"math/bits"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

const (
	// SectionAxisSize is the width of a chunk section on every axis.
	SectionAxisSize = 16
	// BlocksPerSection is Strategy.entryCount() for a block state container:
	// 1 << (bitsPerAxis * 3) with bitsPerAxis = 4.
	BlocksPerSection = SectionAxisSize * SectionAxisSize * SectionAxisSize
	// BiomesPerSection is Strategy.entryCount() for a biome container, which
	// uses bitsPerAxis = 2 — biomes are stored per 4x4x4 cell, not per block.
	BiomesPerSection = 4 * 4 * 4
)

// WorldHeight is the vertical extent of a dimension, mirroring
// LevelHeightAccessor.  Chunk section indices are derived from it rather than
// hard-coded, because only the overworld runs from -64 to 320.
type WorldHeight struct {
	MinY   int
	Height int
}

// Overworld is the height of minecraft:overworld: 384 blocks starting at -64.
var Overworld = WorldHeight{MinY: -64, Height: 384}

func (h WorldHeight) MaxY() int { return h.MinY + h.Height - 1 }

// MinSectionY is LevelHeightAccessor.getMinSectionY(). The shift is an
// arithmetic one so it floors for negative Y, which integer division would not.
func (h WorldHeight) MinSectionY() int { return h.MinY >> 4 }

func (h WorldHeight) MaxSectionY() int { return h.MaxY() >> 4 }

func (h WorldHeight) SectionCount() int { return h.MaxSectionY() - h.MinSectionY() + 1 }

// SectionIndex is LevelHeightAccessor.getSectionIndex: the position of the
// section containing blockY within Chunk.Sections.
func (h WorldHeight) SectionIndex(blockY int) int { return (blockY >> 4) - h.MinSectionY() }

func (h WorldHeight) Contains(blockY int) bool { return blockY >= h.MinY && blockY <= h.MaxY() }

// blockIndex is Strategy.getIndex for a block state container (bitsPerAxis 4):
// blocks are ordered Y, then Z, then X.
func blockIndex(x, y, z int) int { return (y<<4|z)<<4 | x }

// biomeIndex is Strategy.getIndex for a biome container (bitsPerAxis 2).
func biomeIndex(x, y, z int) int { return (y<<2|z)<<2 | x }

// ChunkSection is one 16x16x16 slice of a chunk, laid out like
// LevelChunkSection: a palette plus an index per entry, rather than a full
// Block per position.  Storing 4096 Block values per section instead costs
// roughly 1.5MB per chunk, which is why vanilla does not do it either.
type ChunkSection struct {
	// Y is the absolute section coordinate, i.e. the "Y" tag in the region file.
	Y int

	// NonEmptyBlockCount and FluidCount are the two counters the client reads
	// at the head of every section.
	NonEmptyBlockCount int
	FluidCount         int

	// BlockPalette holds the distinct states in this section; Blocks holds one
	// index into it per position, ordered by blockIndex.  A nil Blocks means
	// every position is BlockPalette[0] — vanilla's ZeroBitStorage.
	BlockPalette []Block
	Blocks       []uint16

	// BiomePalette and Biomes are the same arrangement for the 4x4x4 biome
	// cells, ordered by biomeIndex.
	BiomePalette []string
	Biomes       []uint16
}

// NewChunkSection returns an empty section: all air, all plains, matching the
// container PalettedContainerFactory builds for a missing section.
func NewChunkSection(y int, defaultBiome string) *ChunkSection {
	return &ChunkSection{
		Y:            y,
		BlockPalette: []Block{AirBlock},
		BiomePalette: []string{defaultBiome},
	}
}

// blockAt returns the state at a section-relative position.
func (s *ChunkSection) blockAt(index int) Block {
	if s.Blocks == nil {
		return s.BlockPalette[0]
	}
	return s.BlockPalette[s.Blocks[index]]
}

// setBlockAt writes a state at a section-relative position, extending the
// palette if the state is new (PalettedContainer.set via Palette.idFor).
func (s *ChunkSection) setBlockAt(index int, block Block) {
	id := -1
	for i, entry := range s.BlockPalette {
		if entry.Key() == block.Key() {
			id = i
			break
		}
	}
	if id == -1 {
		s.BlockPalette = append(s.BlockPalette, block)
		id = len(s.BlockPalette) - 1
	}

	if s.Blocks == nil {
		if id == 0 {
			return
		}
		// Promote out of the all-one-value representation.
		s.Blocks = make([]uint16, BlocksPerSection)
	}
	s.Blocks[index] = uint16(id)
}

// RecalcBlockCounts recomputes the two section counters, mirroring
// LevelChunkSection.recalcBlockCounts — it counts per palette entry rather than
// per block, as PalettedContainer.count does.
func (s *ChunkSection) RecalcBlockCounts() {
	counts := make([]int, len(s.BlockPalette))
	if s.Blocks == nil {
		counts[0] = BlocksPerSection
	} else {
		for _, id := range s.Blocks {
			counts[id]++
		}
	}

	s.NonEmptyBlockCount = 0
	s.FluidCount = 0
	for id, count := range counts {
		block := s.BlockPalette[id]
		if block.IsAir() {
			continue
		}
		s.NonEmptyBlockCount += count
		if block.IsFluid() {
			s.FluidCount += count
		}
	}
}

// Chunk is a 16x16 column of sections at a fixed chunk coordinate.
type Chunk struct {
	X, Z     int32
	Height   WorldHeight
	Sections []*ChunkSection
}

// NewChunk returns an empty chunk with one section per level of the dimension.
func NewChunk(x, z int32, height WorldHeight, defaultBiome string) *Chunk {
	sections := make([]*ChunkSection, height.SectionCount())
	for i := range sections {
		sections[i] = NewChunkSection(height.MinSectionY()+i, defaultBiome)
	}
	return &Chunk{X: x, Z: z, Height: height, Sections: sections}
}

// GetBlock returns the state at chunk-relative x/z (0-15) and absolute y.
// Positions outside the build height read as air, as they do in vanilla.
func (c *Chunk) GetBlock(x int, y int, z int) Block {
	if !c.Height.Contains(y) || x < 0 || x > 15 || z < 0 || z > 15 {
		return AirBlock
	}
	return c.Sections[c.Height.SectionIndex(y)].blockAt(blockIndex(x, y&15, z))
}

// SetBlock writes a state at chunk-relative x/z (0-15) and absolute y, ignoring
// positions outside the build height.  Callers must run RecalcBlockCounts (or
// Chunk.RecalcBlockCounts) afterwards before serialising the chunk.
func (c *Chunk) SetBlock(x int, y int, z int, block Block) {
	if !c.Height.Contains(y) || x < 0 || x > 15 || z < 0 || z > 15 {
		return
	}
	// y&15 is the section-relative height: masking floors correctly for
	// negative Y where y%16 would give a negative index.
	c.Sections[c.Height.SectionIndex(y)].setBlockAt(blockIndex(x, y&15, z), block)
}

// RecalcBlockCounts refreshes every section's counters.
func (c *Chunk) RecalcBlockCounts() {
	for _, section := range c.Sections {
		section.RecalcBlockCounts()
	}
}

// ceilLog2 is Mth.ceillog2: the number of bits needed to hold count distinct
// values.  Computed exactly rather than via math.Log, which is only accurate
// here by luck and would silently widen a container by a bit if it rounded up.
func ceilLog2(count int) int {
	if count <= 1 {
		return 0
	}
	return bits.Len32(uint32(count - 1))
}

// blockStateBitsInStorage is the width a block state container is stored at on
// disc: Strategy.createForBlockStates maps 1-4 bits onto a 4-bit configuration,
// 5-8 bits onto themselves, and anything wider onto a global configuration
// whose bitsInStorage is the exact width.  A single-entry palette needs no
// storage at all.
func blockStateBitsInStorage(paletteSize int) int {
	stateBits := ceilLog2(paletteSize)
	if stateBits > 0 && stateBits < 4 {
		return 4
	}
	return stateBits
}

// biomeBitsInStorage is the same for biomes, which have no 4-bit minimum.
func biomeBitsInStorage(paletteSize int) int {
	return ceilLog2(paletteSize)
}

// ParseChunkFromRegionChunk converts a chunk read out of a region file into the
// in-memory form, following PalettedContainer.unpack.
func ParseChunkFromRegionChunk(regionChunk *dataTypes.RegionChunk, height WorldHeight) (*Chunk, error) {
	if regionChunk == nil {
		return nil, fmt.Errorf("nil region chunk")
	}

	chunk := NewChunk(regionChunk.XPos, regionChunk.ZPos, height, "minecraft:plains")

	for i := range regionChunk.Sections {
		regionSection := &regionChunk.Sections[i]
		sectionY := int(regionSection.Y)

		// SerializableChunkData.parse discards sections outside the build
		// height.  Region files carry two extra sections either side of the
		// world holding nothing but light data.
		if sectionY < height.MinSectionY() || sectionY > height.MaxSectionY() {
			continue
		}

		section := chunk.Sections[sectionY-height.MinSectionY()]

		if err := parseBlockStates(section, &regionSection.BlockStates); err != nil {
			return nil, fmt.Errorf("chunk %d,%d section %d block states: %w", chunk.X, chunk.Z, sectionY, err)
		}
		if err := parseBiomes(section, &regionSection.Biomes); err != nil {
			return nil, fmt.Errorf("chunk %d,%d section %d biomes: %w", chunk.X, chunk.Z, sectionY, err)
		}

		section.RecalcBlockCounts()
	}

	return chunk, nil
}

func parseBlockStates(section *ChunkSection, states *dataTypes.ChunkSectionBlockStates) error {
	// A section with no block_states tag is the empty container the factory
	// would build, which the section already is.
	if len(states.Palette) == 0 {
		return nil
	}

	palette := make([]Block, len(states.Palette))
	for i, entry := range states.Palette {
		palette[i] = Block{Type: entry.Name, Properties: entry.Properties}
	}

	entries, err := unpackStorage(states.Data, BlocksPerSection, blockStateBitsInStorage(len(palette)), len(palette))
	if err != nil {
		return err
	}

	section.BlockPalette = palette
	section.Blocks = entries
	return nil
}

func parseBiomes(section *ChunkSection, biomes *dataTypes.ChunkSectionBiome) error {
	if len(biomes.Palette) == 0 {
		return nil
	}

	entries, err := unpackStorage(biomes.Data, BiomesPerSection, biomeBitsInStorage(len(biomes.Palette)), len(biomes.Palette))
	if err != nil {
		return err
	}

	section.BiomePalette = biomes.Palette
	section.Biomes = entries
	return nil
}

// unpackStorage expands a packed long array into one palette index per entry,
// following SimpleBitStorage.unpack.  Values are packed low bits first and
// never straddle a long, so each long holds 64/bits of them with the remaining
// high bits unused.  A nil result means every entry is index 0, which is how a
// zero-bit container is stored.
func unpackStorage(data []int64, entryCount int, bitsInStorage int, paletteSize int) ([]uint16, error) {
	if bitsInStorage == 0 {
		// ZeroBitStorage — a single-entry palette carries no data array.
		return nil, nil
	}

	valuesPerLong := 64 / bitsInStorage
	expectedLongs := (entryCount + valuesPerLong - 1) / valuesPerLong
	if len(data) != expectedLongs {
		// SimpleBitStorage rejects a mismatched array rather than reading
		// past the end of it.
		return nil, fmt.Errorf("invalid length given for storage, got: %d but expected: %d", len(data), expectedLongs)
	}

	mask := uint64(1)<<uint(bitsInStorage) - 1
	entries := make([]uint16, entryCount)

	index := 0
	for _, packed := range data {
		cell := uint64(packed)
		for n := 0; n < valuesPerLong && index < entryCount; n++ {
			id := cell & mask
			if int(id) >= paletteSize {
				return nil, fmt.Errorf("missing palette entry %d, palette has %d entries", id, paletteSize)
			}
			entries[index] = uint16(id)
			index++
			cell >>= uint(bitsInStorage)
		}
	}

	return entries, nil
}
