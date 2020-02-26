package packet

type PlayerPosition struct {
	X        float64 `proto:"double"`
	FeetY    float64 `proto:"double"`
	Z        float64 `proto:"double"`
	OnGround bool    `proto:"bool"`
}

func (pp *PlayerPosition) GetPacketId() int {
	return 0x36
}

/**
func (pp *PlayerPosition) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, pp.X, pp.FeetY, pp.Z, 0, 0)
}
*/
