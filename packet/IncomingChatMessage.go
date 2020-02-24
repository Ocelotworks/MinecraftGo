package packet

import (
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/entity"
)

type IncomingChatMessage struct {
	Message string `proto:"string"`
}

func (icm *IncomingChatMessage) GetPacketId() int {
	return 0x03
}

func (icm *IncomingChatMessage) Handle(packet []byte, connection *Connection) {
	fmt.Println("got chat message")

	chatMessageComponents := []entity.ChatMessageComponent{
		{
			Text: connection.Player.Username,
		},
		{
			Text: icm.Message,
		},
	}

	chatMessage := entity.ChatMessage{
		Translate: "chat.type.text",
		With:      &chatMessageComponents,
	}
	connection.Minecraft.SendMessage(0, chatMessage)

}
