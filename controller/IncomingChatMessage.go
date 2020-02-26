package controller

import (
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type IncomingChatMessage struct {
	CurrentPacket *packetType.IncomingChatMessage
}

func (icm *IncomingChatMessage) GetPacketStruct() packetType.Packet {
	return &packetType.IncomingChatMessage{}
}

func (icm *IncomingChatMessage) Init(currentPacket packetType.Packet) {
	icm.CurrentPacket = currentPacket.(*packetType.IncomingChatMessage)
}

func (icm *IncomingChatMessage) Handle(packet []byte, connection *Connection) {
	chatMessageComponents := []entity.ChatMessageComponent{
		{
			Text: connection.Player.Username,
		},
		{
			Text: icm.CurrentPacket.Message,
		},
	}

	chatMessage := entity.ChatMessage{
		Translate: "chat.type.text",
		With:      &chatMessageComponents,
	}
	connection.Minecraft.SendMessage(0, chatMessage)
}
