package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ServerboundKnownPacks struct {
	KnownPacks []KnownPack `proto:"prefixedArray"`
}

func (ci *ServerboundKnownPacks) GetPacketId() int {
	return constants.SBSelectKnownPacksConfiguration
}
