package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityStatus struct {
	EntityID     int  `proto:"int"`
	EntityStatus byte `proto:"byte"`
}

func (es *EntityStatus) GetPacketId() int {
	return constants.CBEntityStatus
}
