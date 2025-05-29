package controller

import (
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

// This is this way because having them as part of the packet creates an import cycle between controller and packet
// TODO: find a way to structure this that doesn't cause that

type Packet interface {
	GetPacketStruct() packetType.Packet
	Handle(packet []byte, connection *Connection)
	Init(packetType.Packet, *Minecraft)
}
