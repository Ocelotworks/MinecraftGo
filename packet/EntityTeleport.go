package packet

type EntityTeleport struct {
	EntityID int     `proto:"varInt"`
	X        float64 `proto:"double"`
	Y        float64 `proto:"double"`
	Z        float64 `proto:"double"`
	Yaw      byte    `proto:"unsignedByte"`
	Pitch    byte    `proto:"unsignedByte"`
	OnGround bool    `proto:"bool"`
}

func (et *EntityTeleport) GetPacketId() int {
	return 0x56
}

/**
func (et *EntityTeleport) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
