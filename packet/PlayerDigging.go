package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

// TODO: This has changed

type PlayerAction struct {
	Status   int   `proto:"varInt"`
	Location int64 `proto:"long"`
	Face     byte  `proto:"unsignedByte"`
	Sequence int   `proto:"varInt"`
}

func (pd *PlayerAction) GetPacketId() int {
	return constants.SBPlayerAction
}
