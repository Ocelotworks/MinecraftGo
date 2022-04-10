package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityPosition struct {
	EntityID int   `proto:"varInt"`
	DeltaX   int16 `proto:"short"`
	DeltaY   int16 `proto:"short"`
	DeltaZ   int16 `proto:"short"`
	OnGround bool  `proto:"bool"`
}

func (ep *EntityPosition) GetPacketId() int {
	return constants.CBEntityPosition
}

/**
func (ep *EntityPosition) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
