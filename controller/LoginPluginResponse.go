package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type LoginPluginResponse struct {
	CurrentPacket *packetType.LoginPluginResponse
}

func (lpr *LoginPluginResponse) GetPacketStruct() packetType.Packet {
	return &packetType.IncomingChatMessage{}
}

func (lpr *LoginPluginResponse) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	lpr.CurrentPacket = currentPacket.(*packetType.LoginPluginResponse)
}

func (lpr *LoginPluginResponse) Handle(packet []byte, connection *Connection) {
	fmt.Println("Login Plugin Response", lpr.CurrentPacket.MessageId, lpr.CurrentPacket.Data)
}
