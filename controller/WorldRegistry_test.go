package controller

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/entity"
	"github.com/Ocelotworks/MinecraftGo/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBlockData(t *testing.T) map[string]entity.BlockData {
	t.Helper()

	raw, err := os.ReadFile("../data/blocks.json")
	require.NoError(t, err)

	blocks := make([]entity.BlockData, 0)
	require.NoError(t, json.Unmarshal(raw, &blocks))

	byName := make(map[string]entity.BlockData, len(blocks))
	for _, block := range blocks {
		byName["minecraft:"+block.Name] = block
	}
	return byName
}

// A block's states are numbered from minStateId with the last property varying
// fastest, so a state's ID depends on getting that ordering right.
func TestBlockStateID(t *testing.T) {
	blockData := testBlockData(t)
	registry := NewWorldRegistry(blockData, []string{"minecraft:plains", "minecraft:desert"})

	cases := []struct {
		block world.Block
		want  int
	}{
		{world.Block{Type: "minecraft:air"}, blockData["minecraft:air"].DefaultState},
		{world.Block{Type: "minecraft:stone"}, blockData["minecraft:stone"].DefaultState},
		{world.Block{Type: "minecraft:bedrock"}, blockData["minecraft:bedrock"].DefaultState},
		{world.Block{Type: "minecraft:grass_block", Properties: map[string]string{"snowy": "false"}}, blockData["minecraft:grass_block"].DefaultState},
		{world.Block{Type: "minecraft:oak_log", Properties: map[string]string{"axis": "y"}}, blockData["minecraft:oak_log"].DefaultState},
		{world.Block{Type: "minecraft:water", Properties: map[string]string{"level": "0"}}, blockData["minecraft:water"].DefaultState},
		{world.Block{
			Type: "minecraft:oak_stairs",
			Properties: map[string]string{
				"facing": "north", "half": "bottom", "shape": "straight", "waterlogged": "false",
			},
		}, blockData["minecraft:oak_stairs"].DefaultState},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, registry.BlockStateID(c.block), c.block.Key())
	}

	// Every state of a block must land inside that block's own ID range.
	stairs := blockData["minecraft:oak_stairs"]
	id := registry.BlockStateID(world.Block{
		Type: "minecraft:oak_stairs",
		Properties: map[string]string{
			"facing": "east", "half": "top", "shape": "outer_right", "waterlogged": "true",
		},
	})
	assert.GreaterOrEqual(t, id, stairs.MinStateId)
	assert.LessOrEqual(t, id, stairs.MaxStateId)

	// An unknown block falls back to air rather than to an arbitrary ID.
	assert.Equal(t, 0, registry.BlockStateID(world.Block{Type: "minecraft:not_a_block"}))
}

// Biome IDs are the index in sorted order, matching the order KnownPacks sends
// the registry entries in.
func TestBiomeID(t *testing.T) {
	registry := NewWorldRegistry(testBlockData(t), []string{"minecraft:plains", "minecraft:badlands", "minecraft:desert"})

	assert.Equal(t, 0, registry.BiomeID("minecraft:badlands"))
	assert.Equal(t, 1, registry.BiomeID("minecraft:desert"))
	assert.Equal(t, 2, registry.BiomeID("minecraft:plains"))
	assert.Equal(t, 2, registry.GlobalBiomeBits())
}

func TestGlobalBlockStateBits(t *testing.T) {
	registry := NewWorldRegistry(testBlockData(t), nil)

	// The client sizes a direct container from this, so it has to be wide
	// enough for the highest state ID in the block data.
	blockData := testBlockData(t)
	maxState := 0
	for _, block := range blockData {
		if block.MaxStateId > maxState {
			maxState = block.MaxStateId
		}
	}

	assert.GreaterOrEqual(t, 1<<uint(registry.GlobalBlockStateBits()), maxState+1)
	assert.Less(t, 1<<uint(registry.GlobalBlockStateBits()-1), maxState+1)
}
