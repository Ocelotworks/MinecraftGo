package controller

import (
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
		ArrayLength: 1,
		Namespace:   "minecraft",
		ID:          "core",
		Version:     "26.1.2",
	})

	connection.SendPacket(&knownPacks)

}
