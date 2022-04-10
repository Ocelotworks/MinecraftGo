package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type SpawnPosition struct {
	Location int64 `proto:"long"`
}

func (sp *SpawnPosition) GetPacketId() int {
	return constants.CBSpawnPosition
}

/**
func (sp *SpawnPosition) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
