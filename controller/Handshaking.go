package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Handshaking struct {
	CurrentPacket *packetType.Handshaking
}

func (h *Handshaking) GetPacketStruct() packetType.Packet {
	return &packetType.Handshaking{}
}

func (h *Handshaking) Init(currentPacket packetType.Packet) {
	h.CurrentPacket = currentPacket.(*packetType.Handshaking)
}

func (h *Handshaking) Handle(packet []byte, connection *Connection) {
	fmt.Printf("Connection to %s:%d Protocol Version %d\n", h.CurrentPacket.ServerAddress, h.CurrentPacket.ServerPort, h.CurrentPacket.ProtocolVersion)
	fmt.Println("Next State", h.CurrentPacket.NextState)
	connection.State = State(h.CurrentPacket.NextState)

}
