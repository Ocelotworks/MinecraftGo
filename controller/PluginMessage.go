package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type PluginMessage struct {
	CurrentPacket *packetType.PluginMessage
}

func (pm *PluginMessage) GetPacketStruct() packetType.Packet {
	return &packetType.PluginMessage{}
}

func (pm *PluginMessage) Init(currentPacket packetType.Packet) {
	pm.CurrentPacket = currentPacket.(*packetType.PluginMessage)
}

func (pm *PluginMessage) Handle(packet []byte, connection *Connection) {
	fmt.Println("We need to handle this ", pm.CurrentPacket.Identifier)
}
