package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type FinishConfiguration struct {
}

func (et *FinishConfiguration) GetPacketId() int {
	return constants.CBFinishConfiguration
}
