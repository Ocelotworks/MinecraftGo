package controller

import (
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type ClientTickEnd struct {
	CurrentPacket *packetType.ClientTickEnd
}

func (tc *ClientTickEnd) GetPacketStruct() packetType.Packet {
	return &packetType.ClientTickEnd{}
}

func (tc *ClientTickEnd) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	tc.CurrentPacket = currentPacket.(*packetType.ClientTickEnd)
}

func (tc *ClientTickEnd) Handle(packet []byte, connection *Connection) {
}
