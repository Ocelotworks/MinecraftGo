package controller

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
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

	enabledFeatures := packetType.Packet(&packetType.FeatureFlags{
		FeatureFlags: []string{"minecraft:vanilla"},
	})
	connection.SendPacket(&enabledFeatures)

	knownPacks := packetType.Packet(&packetType.ClientboundKnownPacks{
		KnownPacks: []packetType.KnownPack{
			{Namespace: "minecraft", ID: "core", Version: constants.GameVersion},
		},
	})

	connection.SendPacket(&knownPacks)

}
