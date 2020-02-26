package packet

type EntityHeadLook struct {
	EntityID int  `proto:"varInt"`
	Yaw      byte `proto:"unsignedByte"`
}

func (ehl *EntityHeadLook) GetPacketId() int {
	return 0x3C
}

/**
func (ehl *EntityHeadLook) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
