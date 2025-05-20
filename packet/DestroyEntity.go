package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type DestroyEntity struct {
	Count     int   `proto:"varInt"`
	EntityIDs []int `proto:"varIntArray"`
}

func (de *DestroyEntity) GetPacketId() int {
	return constants.CBRemoveEntities
}

/**
func (de *DestroyEntity) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
