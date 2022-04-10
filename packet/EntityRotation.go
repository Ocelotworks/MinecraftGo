package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityRotation struct {
	EntityID int  `proto:"varInt"`
	Yaw      byte `proto:"unsignedByte"`
	Pitch    byte `proto:"unsignedByte"`
	OnGround bool `proto:"bool"`
}

func (er *EntityRotation) GetPacketId() int {
	return constants.CBEntityRotation
}

/**
func (er *EntityRotation) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
