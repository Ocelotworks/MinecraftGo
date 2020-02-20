package packet

type HeldItemChange struct {
	Slot byte `proto:"byte"`
}

func (hic *HeldItemChange) GetPacketId() int {
	return 0x40
}

func (hic *HeldItemChange) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
