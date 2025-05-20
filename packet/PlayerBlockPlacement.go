package packet

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

type UseItemOn struct {
	Hand           int                 `proto:"varInt"`
	Position       *dataTypes.Position `proto:"position"`
	Face           int                 `proto:"varInt"`
	CursorPosX     float32             `proto:"float"`
	CursorPosY     float32             `proto:"float"`
	CursorPosZ     float32             `proto:"float"`
	InsideBlock    bool                `proto:"bool"`
	WorldBorderHit bool                `proto:"bool"`
	Sequence       int                 `proto:"varInt"`
}

func (pbp *UseItemOn) GetPacketId() int {
	return constants.SBUseItemOn
}
