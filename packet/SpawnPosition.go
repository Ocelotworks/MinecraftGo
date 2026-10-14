package packet

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

// SpawnPosition is ClientboundSetDefaultSpawnPositionPacket: the level's
// respawn data, which is a GlobalPos — the dimension the spawn is in and the
// block it is at — followed by the angle to face on respawning.
//
// This is the compass target and the fallback respawn point, not the position
// the joining player is placed at; that is PlayerPositionAndLook.
type SpawnPosition struct {
	DimensionName string              `proto:"string"`
	Location      *dataTypes.Position `proto:"position"`
	Yaw           float32             `proto:"float"`
	Pitch         float32             `proto:"float"`
}

func (sp *SpawnPosition) GetPacketId() int {
	return constants.CBSetDefaultSpawnPosition
}
