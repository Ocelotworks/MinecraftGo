package controller

import (
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/constants"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Interact struct {
	CurrentPacket *packetType.Interact
}

func (a *Interact) GetPacketStruct() packetType.Packet {
	return &packetType.Interact{}
}

func (a *Interact) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	a.CurrentPacket = currentPacket.(*packetType.Interact)
}

func (a *Interact) Handle(packet []byte, connection *Connection) {
	fmt.Println("Interact", a)

	animation := packetType.Packet(&packetType.SwingAnimation{
		EntityID:      a.CurrentPacket.EntityID,
		Hand:          a.CurrentPacket.Hand,
		AnimationType: constants.SwingAnimationWhack,
		Duration:      60,
	})

	connection.Minecraft.SendToAllExcept(connection, &animation)
}
