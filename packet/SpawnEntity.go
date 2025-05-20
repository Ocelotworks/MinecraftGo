package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type SpawnEntity struct {
	EntityID int                  `proto:"varInt"`
	UUID     []byte               `proto:"uuid"`
	Type     constants.EntityType `proto:"varInt"`
	X        float64              `proto:"double"`
	Y        float64              `proto:"double"`
	Z        float64              `proto:"double"`
	Pitch    byte                 `proto:"unsignedByte"`
	Yaw      byte                 `proto:"unsignedByte"`
	HeadYaw  byte                 `proto:"unsignedByte"`
	Data     int                  `proto:"varInt"`
	VelX     int16                `proto:"short"`
	VelY     int16                `proto:"short"`
	VelZ     int16                `proto:"short"`
}

func (sp *SpawnEntity) GetPacketId() int {
	return constants.CBSpawnEntity
}

/**
func (sp *SpawnEntity) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
