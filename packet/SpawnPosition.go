package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type SpawnPosition struct {
	DimensionType    int    `proto:"varInt"`
	DimensionName    string `proto:"string"`
	HashedSeed       int64  `proto:"long"`
	GameMode         byte   `proto:"unsignedByte"`
	PreviousGameMode byte   `proto:"byte"`
	IsDebug          bool   `proto:"bool"`
	IsFlat           bool   `proto:"bool"`
	HasDeathLocation bool   `proto:"bool"`
	// todo: death location
	PortalCooldown int  `proto:"varInt"`
	SeaLevel       int  `proto:"varInt"`
	DataKept       byte `proto:"byte"`
}

func (sp *SpawnPosition) GetPacketId() int {
	return constants.CBDefaultSpawnPosition
}

/**
func (sp *SpawnPosition) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
