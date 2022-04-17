package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerBlockPlacement struct {
	Hand        int     `proto:"varInt"`
	Position    int64   `proto:"long"`
	Face        int     `proto:"varInt"`
	CursorPosX  float32 `proto:"float"`
	CursorPosY  float32 `proto:"float"`
	CursorPosZ  float32 `proto:"float"`
	InsideBlock bool    `proto:"bool"`
}

func (pbp *PlayerBlockPlacement) GetPacketId() int {
	return constants.SBPlayerBlockPlacement
}
