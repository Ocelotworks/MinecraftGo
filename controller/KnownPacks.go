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
	"minecraft:painting_variant",
	"minecraft:dimension_type",
	"minecraft:damage_type",
	"minecraft:banner_pattern",
	//"minecraft:enchantment",
	"minecraft:jukebox_song",
	"minecraft:instrument",
	"minecraft:test_environment",
	"minecraft:test_instance",
	"minecraft:zombie_nautilus_variant",
	"minecraft:dialog",
	"minecraft:timeline",
	"minecraft:world_clock",
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

// buildUpdateTagsPayload constructs the raw binary payload for CBUpdateTagsConfiguration.
//
// Only registries in clientRegistries are included — those are the ones we sent
// RegistryData for, so we own their tag bindings.  Static registries (block, item,
// entity_type, …) are managed entirely by the client and must not be overridden here;
// their tag files may also contain deeply-nested cross-references that would cause a
// stack overflow if we tried to resolve them against an empty entry map.
func (a *KnownPacks) buildUpdateTagsPayload() []byte {
	// Collect only the registries we actually sent RegistryData for that also have tags.
	type regTags struct {
		name string
		tags map[string][]string
	}
	relevant := make([]regTags, 0)
	for _, regName := range clientRegistries {
		if tags, ok := a.Minecraft.Tags[regName]; ok && len(tags) > 0 {
			relevant = append(relevant, regTags{regName, tags})
		}
	}

	buf := make([]byte, 0)
	buf = append(buf, dataTypes.WriteVarInt(len(relevant))...)

	for _, rt := range relevant {
		// Build ID map from the same sorted order used in RegistryData packets.
		// Accept both bare ("day") and namespaced ("minecraft:day") entry name forms.
		registryEntries := a.Minecraft.Registries[rt.name]
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

		buf = append(buf, dataTypes.WriteString(rt.name)...)
		buf = append(buf, dataTypes.WriteVarInt(len(rt.tags))...)

		for tagName, tagValues := range rt.tags {
			ids := resolveTagToIDs(tagValues, rt.tags, entryIDMap)
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
	return &packetType.ClientboundKnownPacks{}
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
