package packet

type ChatMessage struct {
	ChatData string `proto:"string"`
	Position byte   `proto:"unsignedByte"`
}

func (cd *ChatMessage) GetPacketId() int {
	return 0x0F
}

/**
func (cd *ChatMessage) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
