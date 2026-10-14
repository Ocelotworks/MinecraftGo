package world

import (
	"sort"
	"strings"
)

// Block is one block state: the block's registry name plus the property values
// that pick out this particular state.  It is the equivalent of vanilla's
// BlockState, and the properties matter — dropping them collapses every stair,
// slab, log and fence onto its default state.
type Block struct {
	Type       string
	Properties map[string]string
}

// AirBlock is the state an unset block defaults to, matching the container
// vanilla builds when a section has no stored block states.
var AirBlock = Block{Type: "minecraft:air"}

// airTypes mirrors BlockState.isAir(): the three block types the client counts
// as empty space when it reads a section's non-empty block count.
var airTypes = map[string]bool{
	"minecraft:air":      true,
	"minecraft:cave_air": true,
	"minecraft:void_air": true,
}

// IsAir reports whether this state is air, as BlockState.isAir() does.  The
// empty type counts as air so that a zero-valued Block behaves like one.
func (b Block) IsAir() bool {
	return b.Type == "" || airTypes[b.Type]
}

// IsFluid reports whether this state carries a non-empty FluidState — the
// condition LevelChunkSection.recalcBlockCounts uses for its fluid count.  That
// covers the liquid blocks themselves and anything waterlogged.
func (b Block) IsFluid() bool {
	switch b.Type {
	case "minecraft:water", "minecraft:flowing_water", "minecraft:lava", "minecraft:flowing_lava":
		return true
	}
	return b.Properties["waterlogged"] == "true"
}

// Key returns a stable identity for the state, used to deduplicate palette
// entries the way HashMapPalette does.
func (b Block) Key() string {
	if len(b.Properties) == 0 {
		return b.Type
	}

	names := make([]string, 0, len(b.Properties))
	for name := range b.Properties {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString(b.Type)
	sb.WriteByte('[')
	for i, name := range names {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(name)
		sb.WriteByte('=')
		sb.WriteString(b.Properties[name])
	}
	sb.WriteByte(']')
	return sb.String()
}
