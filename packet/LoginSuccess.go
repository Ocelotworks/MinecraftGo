package packet

type LoginSuccess struct {
	UUID     string `proto:"string"`
	Username string `proto:"string"`
}

func (ls *LoginSuccess) GetPacketId() int {
	return 0x02
}

/**
func (ls *LoginSuccess) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
