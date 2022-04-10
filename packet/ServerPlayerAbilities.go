package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ServerPlayerAbilities struct {
	Flags byte `proto:"unsignedByte"`
}

func (spb *ServerPlayerAbilities) GetPacketId() int {
	return constants.CBPlayerAbilities
}
