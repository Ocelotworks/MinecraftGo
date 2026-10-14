package world

import (
	"encoding/binary"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testRegistry hands out IDs in first-seen order, which is enough to check the
// wire layout without depending on the real block data.
type testRegistry struct {
	blocks          map[string]int
	biomes          map[string]int
	globalBlockBits int
	globalBiomeBits int
}

func newTestRegistry() *testRegistry {
	return &testRegistry{
		blocks:          make(map[string]int),
		biomes:          make(map[string]int),
		globalBlockBits: 15,
		globalBiomeBits: 7,
	}
}

func (r *testRegistry) BlockStateID(block Block) int {
	key := block.Key()
	if id, ok := r.blocks[key]; ok {
		return id
	}
	id := len(r.blocks) + 1
	r.blocks[key] = id
	return id
}

func (r *testRegistry) BiomeID(name string) int {
	if id, ok := r.biomes[name]; ok {
		return id
	}
	id := len(r.biomes) + 1
	r.biomes[name] = id
	return id
}

func (r *testRegistry) GlobalBlockStateBits() int { return r.globalBlockBits }
func (r *testRegistry) GlobalBiomeBits() int      { return r.globalBiomeBits }

// clientReader decodes a chunk payload the way the vanilla client does, so that
// a mistake in the number of longs written shows up as a decode failure rather
// than as silently wrong blocks.
type clientReader struct {
	buf []byte
	pos int
}

func (r *clientReader) varInt(t *testing.T) int {
	t.Helper()
	require.Less(t, r.pos, len(r.buf), "ran off the end of the payload")
	value, n := dataTypes.ReadVarInt(r.buf[r.pos:])
	r.pos += n
	return value.(int)
}

func (r *clientReader) short(t *testing.T) int {
	t.Helper()
	require.LessOrEqual(t, r.pos+2, len(r.buf))
	value := int(int16(binary.BigEndian.Uint16(r.buf[r.pos:])))
	r.pos += 2
	return value
}

func (r *clientReader) byteValue(t *testing.T) int {
	t.Helper()
	require.Less(t, r.pos, len(r.buf))
	value := int(r.buf[r.pos])
	r.pos++
	return value
}

func (r *clientReader) longs(t *testing.T, count int) []uint64 {
	t.Helper()
	require.LessOrEqual(t, r.pos+count*8, len(r.buf), "not enough longs in the payload")
	out := make([]uint64, count)
	for i := range out {
		out[i] = binary.BigEndian.Uint64(r.buf[r.pos:])
		r.pos += 8
	}
	return out
}

// container decodes one PalettedContainer, returning the registry ID at every
// position.  bitsInMemory mirrors Strategy.getConfigurationForBitCount: the
// client sizes its storage from the width it would use in memory, not from the
// width on the wire, so a container declared at 2 bits is read at 4.
func (r *clientReader) container(t *testing.T, entryCount int, maxIndirectBits int, globalBits int) []int {
	t.Helper()

	declaredBits := r.byteValue(t)

	if declaredBits == 0 {
		single := r.varInt(t)
		values := make([]int, entryCount)
		for i := range values {
			values[i] = single
		}
		return values
	}

	var palette []int
	bitsInMemory := declaredBits
	if declaredBits <= maxIndirectBits {
		if maxIndirectBits >= 4 && bitsInMemory < 4 {
			bitsInMemory = 4
		}
		size := r.varInt(t)
		palette = make([]int, size)
		for i := range palette {
			palette[i] = r.varInt(t)
		}
	} else {
		bitsInMemory = globalBits
	}

	valuesPerLong := 64 / bitsInMemory
	cells := r.longs(t, (entryCount+valuesPerLong-1)/valuesPerLong)
	mask := uint64(1)<<uint(bitsInMemory) - 1

	values := make([]int, 0, entryCount)
	for _, cell := range cells {
		for n := 0; n < valuesPerLong && len(values) < entryCount; n++ {
			raw := int(cell & mask)
			cell >>= uint(bitsInMemory)
			if palette != nil {
				require.Less(t, raw, len(palette), "palette index out of range")
				raw = palette[raw]
			}
			values = append(values, raw)
		}
	}

	return values
}

// readChunk decodes a whole chunk payload and checks nothing is left over,
// which is the condition the client enforces when it writes the buffer.
func readChunk(t *testing.T, payload []byte, sectionCount int, registry *testRegistry) [][]int {
	t.Helper()

	reader := &clientReader{buf: payload}
	sections := make([][]int, sectionCount)

	for i := 0; i < sectionCount; i++ {
		reader.short(t) // non-empty block count
		reader.short(t) // fluid count
		sections[i] = reader.container(t, BlocksPerSection, maxIndirectBlockStateBits, registry.GlobalBlockStateBits())
		reader.container(t, BiomesPerSection, maxIndirectBiomeBits, registry.GlobalBiomeBits())
	}

	assert.Equal(t, len(payload), reader.pos, "payload has trailing bytes the client would not read")
	return sections
}

func TestWriteNetworkDataRoundTrip(t *testing.T) {
	chunk := loadTestChunk(t)
	registry := newTestRegistry()

	sections := readChunk(t, chunk.WriteNetworkData(registry), len(chunk.Sections), registry)

	for sectionIndex, section := range chunk.Sections {
		for i := 0; i < BlocksPerSection; i++ {
			expected := registry.BlockStateID(section.blockAt(i))
			require.Equal(t, expected, sections[sectionIndex][i],
				"section %d entry %d", sectionIndex, i)
		}
	}
}

// A section holding one state must go out at zero bits with a single-value
// palette and no long array at all.
func TestWriteNetworkDataSingleValueSection(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")
	registry := newTestRegistry()

	payload := chunk.Sections[0].WriteNetworkData(registry)

	// 2 counters + (bits + palette VarInt) for blocks + the same for biomes.
	assert.Len(t, payload, 2+2+1+1+1+1)
}

// Palette widths of 1 to 3 bits have to be sent as 4: the client maps them onto
// a 4-bit storage and would read four times as many longs as we wrote.
func TestWriteNetworkDataUsesFourBitFloor(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")
	chunk.SetBlock(0, 0, 0, Block{Type: "minecraft:stone"})
	chunk.RecalcBlockCounts()

	registry := newTestRegistry()
	section := chunk.Sections[Overworld.SectionIndex(0)]
	payload := section.WriteNetworkData(registry)

	assert.Equal(t, byte(4), payload[4], "expected a 4-bit block state container")

	sections := readChunk(t, payload, 1, registry)
	assert.Equal(t, registry.BlockStateID(Block{Type: "minecraft:stone"}), sections[0][0])
	assert.Equal(t, registry.BlockStateID(AirBlock), sections[0][1])
}

// Beyond 256 states a section falls back to the global palette, which sends no
// palette and stores registry IDs at the client's global width.
func TestWriteNetworkDataFallsBackToGlobalPalette(t *testing.T) {
	chunk := NewChunk(0, 0, Overworld, "minecraft:plains")
	for i := 0; i < 300; i++ {
		chunk.SetBlock(i%16, i/16, 0, Block{
			Type:       "minecraft:stone",
			Properties: map[string]string{"variant": string(rune('a' + i%26)), "n": string(rune('a' + i/26))},
		})
	}
	chunk.RecalcBlockCounts()

	registry := newTestRegistry()
	section := chunk.Sections[Overworld.SectionIndex(0)]
	require.Greater(t, len(section.BlockPalette), 256)

	payload := section.WriteNetworkData(registry)
	assert.Equal(t, byte(registry.GlobalBlockStateBits()), payload[4])

	sections := readChunk(t, payload, 1, registry)
	for i := 0; i < BlocksPerSection; i++ {
		require.Equal(t, registry.BlockStateID(section.blockAt(i)), sections[0][i], "entry %d", i)
	}
}

func TestWriteNetworkDataBiomes(t *testing.T) {
	chunk := loadTestChunk(t)
	registry := newTestRegistry()

	reader := &clientReader{buf: chunk.Sections[0].WriteNetworkData(registry)}
	reader.short(t)
	reader.short(t)
	reader.container(t, BlocksPerSection, maxIndirectBlockStateBits, registry.GlobalBlockStateBits())
	biomes := reader.container(t, BiomesPerSection, maxIndirectBiomeBits, registry.GlobalBiomeBits())

	require.Len(t, biomes, BiomesPerSection)
	for i := 0; i < BiomesPerSection; i++ {
		expected := registry.BiomeID(chunk.Sections[0].BiomePalette[entryIndex(chunk.Sections[0].Biomes, i)])
		require.Equal(t, expected, biomes[i], "biome cell %d", i)
	}
}
