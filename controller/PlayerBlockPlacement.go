package controller

import (
	"fmt"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"math"
)

type PlayerBlockPlacement struct {
	CurrentPacket *packetType.PlayerBlockPlacement
}

func (pbp *PlayerBlockPlacement) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerBlockPlacement{}
}

func (pbp *PlayerBlockPlacement) Init(currentPacket packetType.Packet) {
	pbp.CurrentPacket = currentPacket.(*packetType.PlayerBlockPlacement)
}

func (pbp *PlayerBlockPlacement) Handle(packet []byte, connection *Connection) {

	if pbp.CurrentPacket.Hand != 0 {
		return
	}

	x := pbp.CurrentPacket.Position >> 38
	y := pbp.CurrentPacket.Position & 0xFFF
	z := (pbp.CurrentPacket.Position >> 12) & 0x3FFFFFF

	fmt.Printf("%b\n", pbp.CurrentPacket.Position)
	fmt.Printf("%b\n", y)

	// Chunk Postion
	cx := x / 16
	cz := z / 16

	// Region Position
	rx := int(math.Floor(float64(cx / 32)))
	rz := int(math.Floor(float64(cz / 32)))

	fmt.Printf("Block placement at (%d,%d,%d) chunk (%d,%d) region (%d, %d)\n", x, y, z, cx, cz, rx, rz)

	return
	chunk := connection.Minecraft.DataStore.Map[rx][rz].GetChunk(int32(cx), int32(cz))
	chunk.GetBlockAt(int(x), int(y), int(z))
}
