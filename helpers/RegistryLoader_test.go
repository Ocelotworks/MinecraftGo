package helpers

import (
	"testing"

	"github.com/Ocelotworks/MinecraftGo/helpers/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAllRegistries(t *testing.T) {
	testutil.ChdirToRepoRoot(t)

	registries := LoadAllRegistries()

	// entries are keyed by namespaced name, as they are sent to the client
	dimensions := registries["minecraft:dimension_type"]
	require.NotEmpty(t, dimensions)
	for _, name := range []string{"minecraft:overworld", "minecraft:overworld_caves", "minecraft:the_end", "minecraft:the_nether"} {
		assert.Contains(t, dimensions, name)
	}

	// nested directories name the registry, not the entry
	assert.NotEmpty(t, registries["minecraft:worldgen/biome"])

	// tag files are not registry entries
	for name := range registries {
		assert.NotContains(t, name, "tags", "tag file loaded as a registry: %s", name)
	}
}

func TestLoadAllTags(t *testing.T) {
	testutil.ChdirToRepoRoot(t)

	tags := LoadAllTags()

	// the first path segment is the registry; the rest, nested or not, is the tag
	assert.Contains(t, tags["minecraft:block"], "infiniburn_overworld")
	assert.Contains(t, tags["minecraft:item"], "enchantable/head_armor")
	assert.Contains(t, tags["minecraft:timeline"], "in_overworld")
}

func TestLoadRegistryIDs(t *testing.T) {
	testutil.ChdirToRepoRoot(t)

	ids := LoadRegistryIDs()

	// ids are the client's protocol ids, so air is the first block and stone the first item
	assert.Equal(t, 0, ids["minecraft:block"]["minecraft:air"])
	assert.Equal(t, 1, ids["minecraft:item"]["minecraft:stone"])
	assert.NotEmpty(t, ids["minecraft:entity_type"])
}
