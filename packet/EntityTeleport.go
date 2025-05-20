package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityTeleport struct {
	EntityID int     `proto:"varInt"`
	X        float64 `proto:"double"`
	Y        float64 `proto:"double"`
	Z        float64 `proto:"double"`
	Yaw      byte    `proto:"unsignedByte"`
	Pitch    byte    `proto:"unsignedByte"`
	OnGround bool    `proto:"bool"`
}

func (et *EntityTeleport) GetPacketId() int {
	return constants.CBTeleportEntity
}

/**
func (et *EntityTeleport) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
