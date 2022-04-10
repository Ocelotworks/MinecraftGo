package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

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
	return constants.CBJoinGame
}

/**
func (ls *JoinGame) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
