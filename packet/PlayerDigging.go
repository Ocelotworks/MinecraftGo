package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerDigging struct {
	Status   int   `proto:"varInt"`
	Location int64 `proto:"long"`
	Face     byte  `proto:"unsignedByte"`
}

func (pd *PlayerDigging) GetPacketId() int {
	return constants.SBPlayerDigging
}
