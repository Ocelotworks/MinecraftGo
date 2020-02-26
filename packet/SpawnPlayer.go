package packet

type SpawnPlayer struct {
	EntityID int     `proto:"varInt"`
	UUID     []byte  `proto:"uuid"`
	X        float64 `proto:"double"`
	Y        float64 `proto:"double"`
	Z        float64 `proto:"double"`
	Yaw      byte    `proto:"unsignedByte"`
	Pitch    byte    `proto:"unsignedByte"`
}

func (sp *SpawnPlayer) GetPacketId() int {
	return 0x05
}

/**
func (sp *SpawnPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
