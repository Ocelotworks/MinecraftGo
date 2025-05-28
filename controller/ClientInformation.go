package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type ClientInformation struct {
	CurrentPacket *packetType.ClientInformation
}

func (ci *ClientInformation) GetPacketStruct() packetType.Packet {
	return &packetType.ClientInformation{}
}

func (ci *ClientInformation) Init(currentPacket packetType.Packet) {
	ci.CurrentPacket = currentPacket.(*packetType.ClientInformation)
}

func (ci *ClientInformation) Handle(packet []byte, connection *Connection) {

	fmt.Println("Client Information ", ci.CurrentPacket)

}
