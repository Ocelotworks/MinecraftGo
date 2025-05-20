package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"math"
)

type PlayerBlockPlacement struct {
	CurrentPacket *packetType.UseItemOn
}

func (pbp *PlayerBlockPlacement) GetPacketStruct() packetType.Packet {
	return &packetType.UseItemOn{}
}

func (pbp *PlayerBlockPlacement) Init(currentPacket packetType.Packet) {
	pbp.CurrentPacket = currentPacket.(*packetType.UseItemOn)
}

func (pbp *PlayerBlockPlacement) Handle(packet []byte, connection *Connection) {

	if pbp.CurrentPacket.Hand != 0 {
		return
	}

	// Chunk Position
	cx := pbp.CurrentPacket.Position.X / 16
	cz := pbp.CurrentPacket.Position.Z / 16

	// Region Position
	rx := int(math.Floor(float64(cx / 32)))
	rz := int(math.Floor(float64(cz / 32)))

	fmt.Printf("Block placement at %d chunk (%d,%d) region (%d, %d)\n", pbp.CurrentPacket.Position, cx, cz, rx, rz)

	chunk := connection.Minecraft.DataStore.Map[rx][rz].GetChunk(int32(cx), int32(cz))
	chunk.GetBlockAtPos(pbp.CurrentPacket.Position)

}
