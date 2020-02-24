package packet

type DestroyEntity struct {
	Count     int   `proto:"varInt"`
	EntityIDs []int `proto:"varIntArray"`
}

func (de *DestroyEntity) GetPacketId() int {
	return 0x38
}

func (de *DestroyEntity) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
