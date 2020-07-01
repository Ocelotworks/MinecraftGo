package packet

type ChatMessage struct {
	ChatData string `proto:"string"`
	Position byte   `proto:"unsignedByte"`
	Sender   []byte `proto:"uuid"`
}

func (cd *ChatMessage) GetPacketId() int {
	return 0x0E
}

/**
func (cd *ChatMessage) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
