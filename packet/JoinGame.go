package packet

type JoinGame struct {
	EntityID            int    `proto:"int"`
	Gamemode            byte   `proto:"unsignedByte"`
	Dimension           int    `proto:"int"`
	HashedSeed          int64  `proto:"long"`
	MaxPlayers          byte   `proto:"unsignedByte"`
	LevelType           string `proto:"string"`
	ViewDistance        int    `proto:"varInt"`
	ReducedDebugInfo    bool   `proto:"bool"`
	EnableRespawnScreen bool   `proto:"bool"`
}

func (ls *JoinGame) GetPacketId() int {
	return 0x26
}

func (ls *JoinGame) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
