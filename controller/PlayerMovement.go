package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type PlayerMovement struct {
	CurrentPacket *packetType.PlayerMovement
}

func (pm *PlayerMovement) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerMovement{}
}

func (pm *PlayerMovement) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	pm.CurrentPacket = currentPacket.(*packetType.PlayerMovement)
}

func (pm *PlayerMovement) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Movement", pm)
}
