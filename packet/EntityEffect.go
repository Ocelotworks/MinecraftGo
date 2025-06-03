package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityEffect struct {
	EntityId  int  `proto:"varInt"`
	EffectId  int  `proto:"varInt"`
	Amplifier int  `proto:"varInt"`
	Duration  int  `proto:"varInt"`
	Flags     byte `proto:"byte"`
}

func (bc *EntityEffect) GetPacketId() int {
	return constants.CBEntityEffect
}
