package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ClientboundKnownPacks struct {
	ArrayLength int    `proto:"varInt"`
	Namespace   string `proto:"string"`
	ID          string `proto:"string"`
	Version     string `proto:"string"`
}

func (ci *ClientboundKnownPacks) GetPacketId() int {
	return constants.CBSelectKnownPacksConfiguration
}
