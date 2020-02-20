package packet

type StatusPing struct {
	Payload int `proto:"long"`
}

func (sp *StatusPing) GetPacketId() int {
	return 0x00
}

func (sp *StatusPing) Handle(packet []byte, connection *Connection) {
	//Just send the pong right back
	returnPacket := Packet(sp)
	connection.SendPacket(&returnPacket)
}
