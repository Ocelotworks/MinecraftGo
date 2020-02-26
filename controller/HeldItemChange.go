package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type HeldItemChange struct {
	CurrentPacket *packetType.HeldItemChange
}

func (hic *HeldItemChange) GetPacketStruct() packetType.Packet {
	return &packetType.HeldItemChange{}
}

func (hic *HeldItemChange) Init(currentPacket packetType.Packet) {
	hic.CurrentPacket = currentPacket.(*packetType.HeldItemChange)
}

func (hic *HeldItemChange) Handle(packet []byte, connection *Connection) {
	fmt.Println("Held Item Change ", hic.CurrentPacket)
}
