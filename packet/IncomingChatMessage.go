package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type IncomingChatMessage struct {
	Message string `proto:"string"`
}

func (icm *IncomingChatMessage) GetPacketId() int {
	return constants.SBChatMessage
}

/**
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
*/
