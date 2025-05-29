package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type StatusPing struct {
	CurrentPacket *packetType.StatusPing
}

func (sp *StatusPing) GetPacketStruct() packetType.Packet {
	return &packetType.StatusPing{}
}

func (sp *StatusPing) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	sp.CurrentPacket = currentPacket.(*packetType.StatusPing)
}

func (sp *StatusPing) Handle(packet []byte, connection *Connection) {
	//Just send the pong right back
	fmt.Println("Pingy pongu")
	returnPacket := packetType.Packet(sp.CurrentPacket)
	connection.SendPacket(&returnPacket)
}
