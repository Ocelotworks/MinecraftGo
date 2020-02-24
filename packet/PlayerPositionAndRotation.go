package packet

type PlayerPositionAndRotation struct {
	X        float64 `proto:"double"`
	FeetY    float64 `proto:"double"`
	Z        float64 `proto:"double"`
	Yaw      float32 `proto:"float"`
	Pitch    float32 `proto:"float"`
	OnGround bool    `proto:"bool"`
}

func (ppar *PlayerPositionAndRotation) GetPacketId() int {
	return 0x12
}

func (ppar *PlayerPositionAndRotation) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, ppar.X, ppar.FeetY, ppar.Z, ppar.Yaw, ppar.Pitch)
}
