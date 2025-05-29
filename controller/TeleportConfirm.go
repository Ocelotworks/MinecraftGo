package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type TeleportConfirm struct {
	CurrentPacket *packetType.TeleportConfirm
}

func (tc *TeleportConfirm) GetPacketStruct() packetType.Packet {
	return &packetType.TeleportConfirm{}
}

func (tc *TeleportConfirm) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	tc.CurrentPacket = currentPacket.(*packetType.TeleportConfirm)
}

func (tc *TeleportConfirm) Handle(packet []byte, connection *Connection) {
	fmt.Println("Teleport confirm ", tc.CurrentPacket)
}
