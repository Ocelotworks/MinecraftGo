package packet

type EntityRotation struct {
	EntityID int  `proto:"varInt"`
	Yaw      byte `proto:"unsignedByte"`
	Pitch    byte `proto:"unsignedByte"`
	OnGround bool `proto:"bool"`
}

func (er *EntityRotation) GetPacketId() int {
	return 0x2B
}

func (er *EntityRotation) Handle(packet []byte, connection *Connection) {
	//Client Only
}
