package controller

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type KnownPacks struct {
	CurrentPacket *packetType.ServerboundKnownPacks
	Minecraft     *Minecraft
}

// clientRegistries is the client's RegistryDataLoader.SYNCHRONIZED_REGISTRIES for
// the current GameVersion: every registry the client expects to be told about.  A
// registry missing here makes the client fail with "Missing element" the first time
// something looks up one of its entries (e.g. block_transformer while building item
// components), so this has to be re-checked against the client on each version bump.
var clientRegistries = []string{
	"minecraft:worldgen/biome",
	"minecraft:chat_type",
	"minecraft:trim_pattern",
	"minecraft:trim_material",
	"minecraft:wolf_variant",
	"minecraft:wolf_sound_variant",
	"minecraft:pig_variant",
	"minecraft:pig_sound_variant",
	"minecraft:frog_variant",
	"minecraft:cat_variant",
	"minecraft:cat_sound_variant",
	"minecraft:cow_variant",
	"minecraft:cow_sound_variant",
	"minecraft:chicken_variant",
	"minecraft:chicken_sound_variant",
	"minecraft:zombie_nautilus_variant",
	"minecraft:painting_variant",
	"minecraft:sulfur_cube_archetype",
	"minecraft:dimension_type",
	"minecraft:damage_type",
	"minecraft:banner_pattern",
	"minecraft:enchantment",
	"minecraft:jukebox_song",
	"minecraft:instrument",
	"minecraft:test_environment",
	"minecraft:test_instance",
	"minecraft:dialog",
	"minecraft:world_clock",
	"minecraft:timeline",
	"minecraft:decorated_pot_pattern",
	"minecraft:block_transformer",
	"minecraft:worldgen/block_state_provider",
}

// resolveTagToIDs converts a list of tag values into a flat list of 0-based registry
// entry IDs.  Values may be:
//   - a namespaced entry name ("minecraft:day") → looked up in entryIDMap
//   - a tag reference ("#minecraft:universal") → recursively expanded from allTags
//
// allTags keys are bare tag names ("universal", "in_overworld", …).
// entryIDMap accepts both bare ("day") and namespaced ("minecraft:day") forms.
func resolveTagToIDs(values []string, allTags map[string][]string, entryIDMap map[string]int) []int {
	result := make([]int, 0, len(values))
	for _, v := range values {
		if strings.HasPrefix(v, "#") {
			// Tag reference — strip "#" then optional "minecraft:" to get bare name
			ref := strings.TrimPrefix(strings.TrimPrefix(v, "#"), "minecraft:")
			if refValues, ok := allTags[ref]; ok {
				result = append(result, resolveTagToIDs(refValues, allTags, entryIDMap)...)
			} else {
				fmt.Printf("UpdateTags: unresolved tag reference %q\n", v)
			}
		} else {
			if id, ok := entryIDMap[v]; ok {
				result = append(result, id)
			} else {
				fmt.Printf("UpdateTags: entry %q not found in registry\n", v)
			}
		}
	}
	return result
}

// tagRegistry is one registry's worth of tag bindings for the UpdateTags packet.
type tagRegistry struct {
	name       string
	tags       map[string][]string // tag name -> values, which may reference other tags as "#tag"
	entryIDMap map[string]int      // entry name -> registry ID, bare and namespaced forms
}

// buildUpdateTagsPayload constructs the raw binary payload for CBUpdateTagsConfiguration.
//
// Two kinds of registry get tags, as in vanilla:
//   - the ones in clientRegistries, which we sent RegistryData for, numbered by the
//     sorted order of those entries;
//   - the static registries (block, item, entity_type, ...), which the client owns,
//     numbered by the protocol ids in data/registry_ids.json.
//
// The client's registry data refers to tags in both (an enchantment names item and
// entity_type tags, a dimension_type names block tags) and fails to load if one is
// missing, so the static ones cannot be skipped.
func (a *KnownPacks) buildUpdateTagsPayload() []byte {
	relevant := make([]tagRegistry, 0)

	for _, regName := range clientRegistries {
		tags, ok := a.Minecraft.Tags[regName]
		if !ok || len(tags) == 0 {
			continue
		}

		// Build ID map from the same sorted order used in RegistryData packets.
		// Accept both bare ("day") and namespaced ("minecraft:day") entry name forms.
		registryEntries := a.Minecraft.Registries[regName]
		sortedNames := make([]string, 0, len(registryEntries))
		for name := range registryEntries {
			sortedNames = append(sortedNames, name)
		}
		sort.Strings(sortedNames)

		entryIDMap := make(map[string]int, len(sortedNames)*2)
		for i, name := range sortedNames {
			entryIDMap[name] = i
			entryIDMap["minecraft:"+name] = i
		}

		relevant = append(relevant, tagRegistry{regName, tags, entryIDMap})
	}

	// Static registries: any registry we have both tag files and protocol ids for.
	// Registries with tag files but no ids are dynamic ones the client is not sent.
	staticNames := make([]string, 0)
	for regName := range a.Minecraft.Tags {
		if _, isStatic := a.Minecraft.RegistryIDs[regName]; isStatic {
			staticNames = append(staticNames, regName)
		}
	}
	sort.Strings(staticNames)

	for _, regName := range staticNames {
		tags := a.Minecraft.Tags[regName]
		ids := a.Minecraft.RegistryIDs[regName]

		entryIDMap := make(map[string]int, len(ids)*2)
		for name, id := range ids {
			entryIDMap[name] = id
			entryIDMap[strings.TrimPrefix(name, "minecraft:")] = id
		}

		relevant = append(relevant, tagRegistry{regName, tags, entryIDMap})
	}

	buf := make([]byte, 0)
	buf = append(buf, dataTypes.WriteVarInt(len(relevant))...)

	for _, rt := range relevant {
		buf = append(buf, dataTypes.WriteString(rt.name)...)
		buf = append(buf, dataTypes.WriteVarInt(len(rt.tags))...)

		for tagName, tagValues := range rt.tags {
			ids := resolveTagToIDs(tagValues, rt.tags, rt.entryIDMap)
			// Wire format uses fully-qualified tag identifier: "minecraft:in_overworld"
			buf = append(buf, dataTypes.WriteString("minecraft:"+tagName)...)
			buf = append(buf, dataTypes.WriteVarInt(len(ids))...)
			for _, id := range ids {
				buf = append(buf, dataTypes.WriteVarInt(id)...)
			}
		}
	}
	return buf
}

func (a *KnownPacks) GetPacketStruct() packetType.Packet {
	return &packetType.ServerboundKnownPacks{}
}

func (a *KnownPacks) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	a.CurrentPacket = currentPacket.(*packetType.ServerboundKnownPacks)
	a.Minecraft = minecraft
}

func (a *KnownPacks) Handle(packet []byte, connection *Connection) {
	for _, registryName := range clientRegistries {

		registryEntries := a.Minecraft.Registries[registryName]

		// Sort entry names so IDs are deterministic (required for UpdateTags references).
		sortedNames := make([]string, 0, len(registryEntries))
		for entryName := range registryEntries {
			sortedNames = append(sortedNames, entryName)
		}
		sort.Strings(sortedNames)

		registryDataPacket := packetType.RegistryData{
			RegistryID:  registryName,
			EntryLength: len(registryEntries),
			NBTBytes:    make([]byte, 0),
		}

		for _, entryName := range sortedNames {
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteString(entryName)...)
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteBoolean(false)...)
		}

		completePacket := packetType.Packet(&registryDataPacket)
		connection.SendPacket(&completePacket)
	}

	// Send tag bindings before FinishConfiguration so the client can resolve its
	// built-in tags (e.g. minecraft:timeline#in_overworld) against the entries we sent.
	updateTagsPacket := packetType.Packet(&packetType.UpdateTagsConfiguration{
		Data: a.buildUpdateTagsPayload(),
	})
	connection.SendPacket(&updateTagsPacket)

	finishConfiguration := packetType.Packet(&packetType.FinishConfiguration{})
	connection.SendPacket(&finishConfiguration)
}
