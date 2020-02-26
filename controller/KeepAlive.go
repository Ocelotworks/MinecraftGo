package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type KeepAlive struct {
	CurrentPacket *packetType.KeepAlive
}

func (ka *KeepAlive) GetPacketStruct() packetType.Packet {
	return &packetType.KeepAlive{}
}

func (ka *KeepAlive) Init(currentPacket packetType.Packet) {
	ka.CurrentPacket = currentPacket.(*packetType.KeepAlive)
}

func (ka *KeepAlive) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("KeepAlive", ka.CurrentPacket)
}
