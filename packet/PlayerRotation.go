package packet

type PlayerRotation struct {
	Yaw      float32 `proto:"float"`
	Pitch    float32 `proto:"float"`
	OnGround bool    `proto:"bool"`
}

func (pr *PlayerRotation) GetPacketId() int {
	return 0x14
}

/**
func (pr *PlayerRotation) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, 0, 0, 0, pr.Yaw, pr.Pitch)
}
*/
