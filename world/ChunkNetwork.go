package world

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

// Registry resolves parsed states and biome names onto the numeric IDs the
// client uses, standing in for Strategy.globalMap()/globalPalette().
type Registry interface {
	// BlockStateID is the client's global block state ID for a state.
	BlockStateID(block Block) int
	// BiomeID is the client's ID for a biome in the biome registry sent
	// during configuration.
	BiomeID(name string) int
	// GlobalBlockStateBits and GlobalBiomeBits are Strategy's
	// globalPaletteBitsInMemory for each registry.  When a container falls
	// back to the global palette the client sizes its storage from these and
	// ignores the width we send, so a container written direct must use
	// exactly this many bits or the client reads the wrong number of longs.
	GlobalBlockStateBits() int
	GlobalBiomeBits() int
}

// maxIndirectBlockStateBits is the widest section palette the client will
// accept for block states before falling back to the global palette; likewise
// maxIndirectBiomeBits for biomes.  See Strategy.createForBlockStates and
// Strategy.createForBiomes.
const (
	maxIndirectBlockStateBits = 8
	maxIndirectBiomeBits      = 3
)

// WriteNetworkData serialises every section of the chunk in the order the
// client reads them, matching ClientboundLevelChunkPacketData.extractChunkData.
// The result is the payload of the chunk packet's Data field; the client
// requires it to describe exactly as many sections as the dimension has, with
// nothing left over.
func (c *Chunk) WriteNetworkData(registry Registry) []byte {
	output := make([]byte, 0, len(c.Sections)*2048)
	for _, section := range c.Sections {
		output = append(output, section.WriteNetworkData(registry)...)
	}
	return output
}

// WriteNetworkData serialises one section as LevelChunkSection.write does: the
// two counters, then the block state container, then the biome container.
func (s *ChunkSection) WriteNetworkData(registry Registry) []byte {
	output := make([]byte, 0, 2048)
	output = append(output, dataTypes.WriteShort(int16(s.NonEmptyBlockCount))...)
	output = append(output, dataTypes.WriteShort(int16(s.FluidCount))...)

	blockIDs := make([]int, len(s.BlockPalette))
	for i, block := range s.BlockPalette {
		blockIDs[i] = registry.BlockStateID(block)
	}
	blockBits, blockDirect := containerBits(len(blockIDs), maxIndirectBlockStateBits, registry.GlobalBlockStateBits())
	output = append(output, writeContainer(blockIDs, s.Blocks, BlocksPerSection, blockBits, blockDirect)...)

	biomeIDs := make([]int, len(s.BiomePalette))
	for i, name := range s.BiomePalette {
		biomeIDs[i] = registry.BiomeID(name)
	}
	biomeBits, biomeDirect := containerBits(len(biomeIDs), maxIndirectBiomeBits, registry.GlobalBiomeBits())
	output = append(output, writeContainer(biomeIDs, s.Biomes, BiomesPerSection, biomeBits, biomeDirect)...)

	return output
}

// containerBits picks the width to send a container at, following
// Strategy.getConfigurationForBitCount.  Block states have a 4-bit floor, both
// kinds have a ceiling above which the global palette takes over, and a
// single-entry palette is sent at zero bits with no data at all.
//
// The values between 1 and the floor matter: the client maps a declared width
// of 1, 2 or 3 bits onto a 4-bit storage, so sending the narrower width would
// make it read four times as many longs as we wrote.
func containerBits(paletteSize int, maxIndirectBits int, globalBits int) (bits int, direct bool) {
	entryBits := ceilLog2(paletteSize)
	switch {
	case entryBits == 0:
		return 0, false
	case maxIndirectBits >= 4 && entryBits < 4:
		return 4, false
	case entryBits <= maxIndirectBits:
		return entryBits, false
	default:
		return globalBits, true
	}
}

// writeContainer serialises one PalettedContainer as PalettedContainer.Data.write
// does: the bit width as a byte, the palette, then the packed longs with no
// length prefix — the client knows the length from the width it was given.
//
// entries holds one palette index per position, or nil when every position is
// index 0.  paletteIDs maps those indices onto registry IDs.
func writeContainer(paletteIDs []int, entries []uint16, entryCount int, bitsPerEntry int, direct bool) []byte {
	output := []byte{byte(bitsPerEntry)}

	if bitsPerEntry == 0 {
		// SingleValuePalette writes one VarInt and ZeroBitStorage has no longs.
		return append(output, dataTypes.WriteVarInt(paletteIDs[0])...)
	}

	// values is what actually goes into the bit storage: palette indices for
	// an indirect container, registry IDs for a direct one.
	values := make([]int, entryCount)
	if direct {
		// GlobalPalette.write writes nothing at all.
		for i := range values {
			values[i] = paletteIDs[entryIndex(entries, i)]
		}
	} else {
		output = append(output, dataTypes.WriteVarInt(len(paletteIDs))...)
		for _, id := range paletteIDs {
			output = append(output, dataTypes.WriteVarInt(id)...)
		}
		for i := range values {
			values[i] = int(entryIndex(entries, i))
		}
	}

	// SimpleBitStorage packs values low bits first, never straddling a long.
	valuesPerLong := 64 / bitsPerEntry
	longCount := (entryCount + valuesPerLong - 1) / valuesPerLong
	mask := uint64(1)<<uint(bitsPerEntry) - 1

	for longIndex := 0; longIndex < longCount; longIndex++ {
		cell := uint64(0)
		for n := 0; n < valuesPerLong; n++ {
			i := longIndex*valuesPerLong + n
			if i >= entryCount {
				break
			}
			cell |= (uint64(values[i]) & mask) << uint(n*bitsPerEntry)
		}
		output = append(output, dataTypes.WriteLong(int64(cell))...)
	}

	return output
}

// entryIndex reads a palette index out of an entry array, treating nil as the
// zero-bit case where every position holds index 0.
func entryIndex(entries []uint16, i int) uint16 {
	if entries == nil {
		return 0
	}
	return entries[i]
}
