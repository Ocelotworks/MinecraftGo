package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Animation struct {
	CurrentPacket *packetType.Animation
}

func (a *Animation) GetPacketStruct() packetType.Packet {
	return &packetType.Animation{}
}

func (a *Animation) Init(currentPacket packetType.Packet) {
	a.CurrentPacket = currentPacket.(*packetType.Animation)
}

func (a *Animation) Handle(packet []byte, connection *Connection) {
	fmt.Println("Animation", a)

	action := byte(0)

	if a.CurrentPacket.Hand == 1 {
		action = 3
	}

	animation := packetType.Packet(&packetType.EntityAnimation{
		EntityID:  connection.Player.EntityID,
		Animation: action,
	})

	connection.Minecraft.SendToAllExcept(connection, &animation)
}
