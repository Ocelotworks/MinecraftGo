package packet

type HeldItemChange struct {
	Slot     byte `proto:"unsignedByte"`
	IsServer bool
}

func (hic *HeldItemChange) GetPacketId() int {
	if hic.IsServer {
		return 0x23
	}
	return 0x40
}

/**
func (hic *HeldItemChange) Handle(packet []byte, connection *Connection) {
	fmt.Println("Held Item Change ", hic)
}
*/
