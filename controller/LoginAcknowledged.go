package controller

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type LoginAcknowledged struct {
	CurrentPacket *packetType.LoginAcknowledged
	Minecraft     *Minecraft
}

func (lpr *LoginAcknowledged) GetPacketStruct() packetType.Packet {
	return &packetType.LoginAcknowledged{}
}

func (lpr *LoginAcknowledged) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	lpr.CurrentPacket = currentPacket.(*packetType.LoginAcknowledged)
	lpr.Minecraft = minecraft
}

var clientRegistries = []string{
	"minecraft:worldgen/biome",
	"minecraft:chat_type",
	"minecraft:trim_pattern",
	"minecraft:trim_material",
	"minecraft:wolf_variant",
	"minecraft:wolf_sound_variant",
	"minecraft:pig_variant",
	"minecraft:frog_variant",
	"minecraft:cat_variant",
	"minecraft:cow_variant",
	"minecraft:chicken_variant",
	"minecraft:painting_variant",
	"minecraft:dimension_type",
	"minecraft:damage_type",
	"minecraft:banner_pattern",
	//"minecraft:enchantment",
	"minecraft:jukebox_song",
	"minecraft:instrument",
	"minecraft:test_environment",
	"minecraft:test_instance",
}

func (lpr *LoginAcknowledged) Handle(packet []byte, connection *Connection) {

	connection.State = CONFIGURATION
	knownPacks := packetType.Packet(&packetType.KnownPacks{
		ArrayLength: 1,
		Namespace:   "minecraft",
		ID:          "core",
		Version:     "1.21.5",
	})

	connection.SendPacket(&knownPacks)

	for _, registryName := range clientRegistries {

		registryEntries := lpr.Minecraft.Registries[registryName]

		registryDataPacket := packetType.RegistryData{
			RegistryID:  registryName,
			EntryLength: len(registryEntries),
			NBTBytes:    make([]byte, 0),
		}

		for entryName, _ := range registryEntries {
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteString(entryName)...)
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteBoolean(false)...)
			//networkWrappedEntry := nbt.NetworkWrapperCompound(entryData)
			//actualNbtBytes := networkWrappedEntry.Write()
			//
			//registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, actualNbtBytes...)
		}

		completePacket := packetType.Packet(&registryDataPacket)

		connection.SendPacket(&completePacket)
	}

	finishConfiguration := packetType.Packet(&packetType.FinishConfiguration{})
	connection.SendPacket(&finishConfiguration)

}
