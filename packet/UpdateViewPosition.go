package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type UpdateViewPosition struct {
	ChunkX int `proto:"varInt"`
	ChunkZ int `proto:"varInt"`
}

func (uvp *UpdateViewPosition) GetPacketId() int {
	return constants.CBSetCenterChunk
}
