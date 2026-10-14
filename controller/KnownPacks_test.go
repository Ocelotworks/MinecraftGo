package controller

import (
	"bytes"
	"testing"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
	"github.com/Ocelotworks/MinecraftGo/helpers/testutil"
)

// The client's registry data refers to tags in registries it owns (an enchantment
// names item and entity_type tags, a dimension_type names block tags), so they must
// be in the UpdateTags payload or the client fails to load its registries.
func TestUpdateTagsIncludesStaticRegistryTags(t *testing.T) {
	mc := &Minecraft{
		Registries: map[string]map[string]nbt.NBTValue{},
		RegistryIDs: map[string]map[string]int{
			"minecraft:block": {"minecraft:netherrack": 100, "minecraft:magma_block": 200, "minecraft:bedrock": 300},
		},
		Tags: map[string]map[string][]string{
			"minecraft:block": {
				"infiniburn_overworld": {"minecraft:netherrack", "minecraft:magma_block"},
				"infiniburn_end":       {"#minecraft:infiniburn_overworld", "minecraft:bedrock"},
			},
			// has tag files but no protocol ids: a dynamic registry the client is not sent
			"minecraft:villager_trade": {"some_tag": {"minecraft:x"}},
		},
	}
	payload := (&KnownPacks{Minecraft: mc}).buildUpdateTagsPayload()

	header := append(dataTypes.WriteVarInt(1), dataTypes.WriteString("minecraft:block")...)
	if !bytes.HasPrefix(payload, header) {
		t.Fatalf("payload should hold only the block registry, got %x", payload)
	}
	if bytes.Contains(payload, dataTypes.WriteString("minecraft:villager_trade")) {
		t.Errorf("registry without protocol ids was sent")
	}

	// infiniburn_end expands the referenced overworld tag: netherrack, magma_block, bedrock.
	end := append(dataTypes.WriteString("minecraft:infiniburn_end"), dataTypes.WriteVarInt(3)...)
	end = append(end, dataTypes.WriteVarInt(100)...)
	end = append(end, dataTypes.WriteVarInt(200)...)
	end = append(end, dataTypes.WriteVarInt(300)...)
	if !bytes.Contains(payload, end) {
		t.Errorf("infiniburn_end not resolved to block ids 100,200,300")
	}
}

// Builds the payload from the real data files, which catches tag cycles and
// registries whose tags cannot be resolved to ids.
func TestUpdateTagsWithRealData(t *testing.T) {
	testutil.ChdirToRepoRoot(t) // data/ is relative to the repo root

	mc := &Minecraft{Registries: map[string]map[string]nbt.NBTValue{}}
	mc.LoadRegistries()

	payload := (&KnownPacks{Minecraft: mc}).buildUpdateTagsPayload()

	for _, want := range []string{
		"minecraft:enchantable/head_armor", // item, needed by enchantment
		"minecraft:sensitive_to_smite",     // entity_type, needed by enchantment
		"minecraft:soul_speed_blocks",      // block, needed by enchantment
		"minecraft:infiniburn_overworld",   // block, needed by dimension_type
		"minecraft:in_overworld",           // timeline, needed by dimension_type
	} {
		if !bytes.Contains(payload, dataTypes.WriteString(want)) {
			t.Errorf("tag %s missing from payload", want)
		}
	}
	t.Logf("payload is %d bytes", len(payload))
}

// A registry that is listed but empty is sent to the client as no entries, which it
// takes as the server really having none, and fails loading anything that needs one.
func TestClientRegistriesAreLoaded(t *testing.T) {
	testutil.ChdirToRepoRoot(t)

	mc := &Minecraft{Registries: map[string]map[string]nbt.NBTValue{}}
	mc.LoadRegistries()

	for _, name := range clientRegistries {
		if len(mc.Registries[name]) == 0 {
			t.Errorf("%s is in clientRegistries but no entries were loaded for it", name)
		}
	}
}
