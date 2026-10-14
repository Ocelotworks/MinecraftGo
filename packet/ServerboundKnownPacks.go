package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ServerboundKnownPacks struct {
	ArrayLength int    `proto:"varInt"`
	Namespace   string `proto:"string"`
	ID          string `proto:"string"`
	Version     string `proto:"string"`
}

func (ci *ServerboundKnownPacks) GetPacketId() int {
	return constants.SBSelectKnownPacksConfiguration
}
