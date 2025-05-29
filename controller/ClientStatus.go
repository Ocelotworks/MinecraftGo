package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type ClientStatus struct {
	CurrentPacket *packetType.ClientStatus
}

func (a *ClientStatus) GetPacketStruct() packetType.Packet {
	return &packetType.ClientStatus{}
}

func (a *ClientStatus) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	a.CurrentPacket = currentPacket.(*packetType.ClientStatus)
}

func (a *ClientStatus) Handle(packet []byte, connection *Connection) {
	fmt.Println("client status", a.CurrentPacket.Action)
}
