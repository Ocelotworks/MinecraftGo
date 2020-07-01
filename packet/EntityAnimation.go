package packet

type EntityAnimation struct {
	EntityID  int  `proto:"varInt"`
	Animation byte `proto:"unsignedByte"`
}

func (ea *EntityAnimation) GetPacketId() int {
	return 0x05
}

/**
func (ea *EntityAnimation) Handle(packet []byte, connection *Connection) {
	//Client only packet
}
*/
