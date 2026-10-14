package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type Interact struct {
	EntityID int            `proto:"varInt"`
	Hand     constants.Hand `proto:"varInt"`
}

func (a *Interact) GetPacketId() int {
	return constants.SBInteract
}
