package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type LoginAcknowledged struct {
}

func (la *LoginAcknowledged) GetPacketId() int {
	return constants.SBLoginAcknowledged
}
