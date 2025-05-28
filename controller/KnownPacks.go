package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type KnownPacks struct {
	CurrentPacket *packetType.KnownPacks
}

func (a *KnownPacks) GetPacketStruct() packetType.Packet {
	return &packetType.KnownPacks{}
}

func (a *KnownPacks) Init(currentPacket packetType.Packet) {
	a.CurrentPacket = currentPacket.(*packetType.KnownPacks)
}

func (a *KnownPacks) Handle(packet []byte, connection *Connection) {
	fmt.Println("known packs", a.CurrentPacket)
}
