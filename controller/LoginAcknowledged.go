package controller

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
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

func (lpr *LoginAcknowledged) Handle(packet []byte, connection *Connection) {

	connection.State = CONFIGURATION
	knownPacks := packetType.Packet(&packetType.KnownPacks{
		ArrayLength: 1,
		Namespace:   "minecraft",
		ID:          "core",
		Version:     "1.21.5",
	})

	connection.SendPacket(&knownPacks)

	for registryName, registryEntries := range lpr.Minecraft.Registries {

		registryDataPacket := packetType.RegistryData{
			RegistryID:  registryName,
			EntryLength: len(registryEntries),
			NBTBytes:    make([]byte, 0),
		}

		for entryName, entryData := range registryEntries {
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteString(entryName)...)
			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, dataTypes.WriteBoolean(true)...)
			networkWrappedEntry := nbt.NetworkWrapperCompound(entryData)
			actualNbtBytes := networkWrappedEntry.Write()

			registryDataPacket.NBTBytes = append(registryDataPacket.NBTBytes, actualNbtBytes...)
		}

		completePacket := packetType.Packet(&registryDataPacket)

		connection.SendPacket(&completePacket)
	}

	finishConfiguration := packetType.Packet(&packetType.FinishConfiguration{})
	connection.SendPacket(&finishConfiguration)

}
