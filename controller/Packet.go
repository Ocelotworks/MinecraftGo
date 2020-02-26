package controller

import (
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Packet interface {
	GetPacketStruct() packetType.Packet
	Handle(packet []byte, connection *Connection)
	Init(packetType.Packet)
}
