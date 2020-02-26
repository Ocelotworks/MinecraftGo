package packet

type Disconnect struct {
	Reason string `proto:"string"`
}

func (d *Disconnect) GetPacketId() int {
	return 0x01
}

/**
func (d *Disconnect) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
