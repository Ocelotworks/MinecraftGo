package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityHeadLook struct {
	EntityID int  `proto:"varInt"`
	Yaw      byte `proto:"unsignedByte"`
}

func (ehl *EntityHeadLook) GetPacketId() int {
	return constants.CBSetHeadRotation
}

/**
func (ehl *EntityHeadLook) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
