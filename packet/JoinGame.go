package packet

type JoinGame struct {
	EntityID            int      `proto:"int"`
	Gamemode            byte     `proto:"unsignedByte"`
	PreviousGamemode    byte     `proto:"unsignedByte"`
	WorldNames          []string `proto:"stringArray"`
	DimensionCodec      []byte   `proto:"raw"`
	Dimension           int      `proto:"int"`
	HashedSeed          int64    `proto:"long"`
	MaxPlayers          byte     `proto:"unsignedByte"`
	ViewDistance        int      `proto:"varInt"`
	ReducedDebugInfo    bool     `proto:"bool"`
	EnableRespawnScreen bool     `proto:"bool"`
	IsDebug             bool     `proto:"bool"`
	IsFlat              bool     `proto:"bool"`
}

func (ls *JoinGame) GetPacketId() int {
	return 0x25
}

/**
func (ls *JoinGame) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
