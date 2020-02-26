package packet

//TODO
type PlayerAbilities struct {
}

func (pa *PlayerAbilities) GetPacketId() int {
	return 0x32
}

/**
func (pa *PlayerAbilities) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
