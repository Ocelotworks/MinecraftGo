package controller

import (
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type PlayerPosition struct {
	CurrentPacket *packetType.PlayerPosition
}

func (pp *PlayerPosition) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerPosition{}
}

func (pp *PlayerPosition) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	pp.CurrentPacket = currentPacket.(*packetType.PlayerPosition)
}

func (pp *PlayerPosition) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, pp.CurrentPacket.X, pp.CurrentPacket.FeetY, pp.CurrentPacket.Z, 0, 0)
}
