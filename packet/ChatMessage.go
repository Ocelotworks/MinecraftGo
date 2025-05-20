package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ChatMessage struct {
	ChatData string `proto:"string"`
	Position byte   `proto:"unsignedByte"`
	Sender   []byte `proto:"uuid"`
}

func (cd *ChatMessage) GetPacketId() int {
	return constants.CBPlayerChatMessage
}

/**
func (cd *ChatMessage) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
