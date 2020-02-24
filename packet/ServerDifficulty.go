package packet

type ServerDifficulty struct {
	Difficulty       byte `proto:"unsignedByte"`
	DifficultyLocked bool `proto:"bool"`
}

func (sd *ServerDifficulty) GetPacketId() int {
	return 0x0E
}

func (sd *ServerDifficulty) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
