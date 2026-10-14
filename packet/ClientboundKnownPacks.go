package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ClientboundKnownPacks struct {
	KnownPacks []KnownPack `proto:"prefixedArray"`
}

func (ci *ClientboundKnownPacks) GetPacketId() int {
	return constants.CBSelectKnownPacksConfiguration
}
