package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityPositionAndRotation struct {
	EntityID int   `proto:"varInt"`
	DeltaX   int16 `proto:"short"`
	DeltaY   int16 `proto:"short"`
	DeltaZ   int16 `proto:"short"`
	Yaw      byte  `proto:"unsignedByte"`
	Pitch    byte  `proto:"unsignedByte"`
	OnGround bool  `proto:"bool"`
}

func (epar *EntityPositionAndRotation) GetPacketId() int {
	return constants.CBEntityPositionAndRotation
}

/**
func (epar *EntityPositionAndRotation) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
