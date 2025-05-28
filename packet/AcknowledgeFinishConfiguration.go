package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type AcknowledgeFinishConfiguration struct {
}

func (la *AcknowledgeFinishConfiguration) GetPacketId() int {
	return constants.SBAcknowledgeFinishConfiguration
}
