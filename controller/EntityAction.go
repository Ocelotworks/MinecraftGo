package controller

import (
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type EntityAction struct {
	CurrentPacket *packetType.EntityAction
}

func (ea *EntityAction) GetPacketStruct() packetType.Packet {
	return &packetType.EntityAction{}
}

func (ea *EntityAction) Init(currentPacket packetType.Packet) {
	ea.CurrentPacket = currentPacket.(*packetType.EntityAction)
}

func (ea *EntityAction) Handle(packet []byte, connection *Connection) {
	//TODO: handle
	fmt.Println("Entity action ", ea)

	if ea.CurrentPacket.ActionID == 0 || ea.CurrentPacket.ActionID == 1 {
		fmt.Println("Sneaky unsneaky", ea.CurrentPacket.ActionID)
		currentEffect := byte(0x00)
		pose := 0

		if ea.CurrentPacket.ActionID == 0 {
			currentEffect = 0x02 //Crouching
			pose = 5             //Sneaking
		}

		updateMetadata := packetType.Packet(&packetType.EntityMetadata{
			EntityID: connection.Player.EntityID,
			Metadata: &dataTypes.EntityMetadata{
				Effect: &currentEffect,
				Pose:   &pose,
			},
		})

		connection.Minecraft.SendToAllExcept(connection, &updateMetadata)
	}
}
