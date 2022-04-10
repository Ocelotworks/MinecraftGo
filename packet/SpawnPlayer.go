package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

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
	return constants.CBSpawnPlayer
}

/**
func (sp *SpawnPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
