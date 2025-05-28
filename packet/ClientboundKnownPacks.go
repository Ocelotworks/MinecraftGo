package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type KnownPacks struct {
	ArrayLength int    `proto:"varInt"`
	Namespace   string `proto:"string"`
	ID          string `proto:"string"`
	Version     string `proto:"string"`
}

func (ci *KnownPacks) GetPacketId() int {
	return constants.CBKnownPacksConfiguration
}
