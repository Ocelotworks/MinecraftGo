package packet

type SetCompression struct {
	Threshold int `proto:"varInt"`
}

func (sc *SetCompression) GetPacketId() int {
	return 0x03
}

func (sc *SetCompression) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
