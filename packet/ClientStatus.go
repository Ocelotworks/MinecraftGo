package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ClientStatus struct {
	Action int `proto:"varInt"`
}

func (cs *ClientStatus) GetPacketId() int {
	return constants.SBClientStatus
}
