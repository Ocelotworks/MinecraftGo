package controller

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/Ocelotworks/MinecraftGo/entity"
	"github.com/Ocelotworks/MinecraftGo/world"
)

// WorldRegistry maps parsed block states and biome names onto the numeric IDs
// the client expects in chunk packets.  It is the server-side stand-in for the
// client's global palettes (vanilla: Strategy.globalMap()).
type WorldRegistry struct {
	// blockStates is keyed by world.Block.Key(), i.e. name plus sorted
	// properties, so that states of the same block are told apart.
	blockStates map[string]int
	// defaultStates falls back to a block's default state when a specific
	// combination of properties is not recognised.
	defaultStates map[string]int

	biomes map[string]int

	globalBlockStateBits int
	globalBiomeBits      int
}

// NewWorldRegistry builds the ID lookups from the loaded block data and the
// names in the biome registry that was sent to the client during configuration.
//
// Biome IDs are the index of each entry in the same sorted order KnownPacks
// uses when it sends RegistryData — the client numbers entries by the order it
// receives them, so the two must agree.
func NewWorldRegistry(blockData map[string]entity.BlockData, biomeNames []string) *WorldRegistry {
	registry := &WorldRegistry{
		blockStates:   make(map[string]int),
		defaultStates: make(map[string]int),
		biomes:        make(map[string]int),
	}

	maxStateID := 0
	for name, block := range blockData {
		registry.defaultStates[name] = block.DefaultState
		if block.MaxStateId > maxStateID {
			maxStateID = block.MaxStateId
		}
		for key, id := range expandBlockStates(name, block) {
			registry.blockStates[key] = id
		}
	}

	sorted := append([]string(nil), biomeNames...)
	sort.Strings(sorted)
	for id, name := range sorted {
		registry.biomes[name] = id
	}

	registry.globalBlockStateBits = ceilLog2(maxStateID + 1)
	registry.globalBiomeBits = ceilLog2(len(sorted))

	return registry
}

// expandBlockStates enumerates every state of a block and returns the state ID
// for each, keyed the same way world.Block.Key() keys them.
//
// A block's states are numbered from minStateId as a mixed-radix count over its
// properties, with the last property varying fastest — the order vanilla's
// StateDefinition produces.
func expandBlockStates(name string, block entity.BlockData) map[string]int {
	if len(block.States) == 0 {
		return map[string]int{name: block.MinStateId}
	}

	// The game sorts a block's properties by name before counting through
	// them, so the radix order is alphabetical whatever order the data file
	// happens to list them in.
	properties := append([]entity.BlockState(nil), block.States...)
	sort.Slice(properties, func(i, j int) bool { return properties[i].Name < properties[j].Name })

	// values[i] is the list of possible values for property i, in the order
	// that decides their numbering.
	values := make([][]string, len(properties))
	total := 1
	for i, state := range properties {
		if state.Type == "bool" && len(state.Values) == 0 {
			// BooleanProperty enumerates true before false.
			values[i] = []string{"true", "false"}
		} else {
			values[i] = state.Values
		}
		if len(values[i]) == 0 {
			return map[string]int{name: block.MinStateId}
		}
		total *= len(values[i])
	}

	states := make(map[string]int, total)
	for offset := 0; offset < total; offset++ {
		combination := make(map[string]string, len(properties))

		// Walk the properties right to left, peeling off the fastest-varying
		// one first.
		remaining := offset
		for i := len(properties) - 1; i >= 0; i-- {
			size := len(values[i])
			combination[properties[i].Name] = values[i][remaining%size]
			remaining /= size
		}

		key := world.Block{Type: name, Properties: combination}.Key()
		states[key] = block.MinStateId + offset
	}

	return states
}

// BlockStateID returns the client's global block state ID for a state, falling
// back to the block's default state when the exact property combination is
// unknown and to air when the block itself is.
func (r *WorldRegistry) BlockStateID(block world.Block) int {
	if id, ok := r.blockStates[block.Key()]; ok {
		return id
	}
	if id, ok := r.defaultStates[block.Type]; ok {
		return id
	}
	if block.Type != "" {
		fmt.Println("Unknown block state:", block.Key())
	}
	return 0
}

// BiomeID returns the ID of a biome in the registry sent to the client.
func (r *WorldRegistry) BiomeID(name string) int {
	if id, ok := r.biomes[name]; ok {
		return id
	}
	return 0
}

func (r *WorldRegistry) GlobalBlockStateBits() int { return r.globalBlockStateBits }

func (r *WorldRegistry) GlobalBiomeBits() int { return r.globalBiomeBits }

// ceilLog2 is Mth.ceillog2: the bits needed to index count distinct values.
func ceilLog2(count int) int {
	if count <= 1 {
		return 0
	}
	return bits.Len32(uint32(count - 1))
}
