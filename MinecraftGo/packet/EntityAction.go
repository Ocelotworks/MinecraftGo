package packet

import "fmt"

import "../dataTypes"

type EntityAction struct {
	EntityID  int `proto:"varInt"`
	ActionID  int `proto:"varInt"`
	JumpBoost int `proto:"varInt"`
}

func (ea *EntityAction) GetPacketId() int {
	return 0x1B
}

func (ea *EntityAction) Handle(packet []byte, connection *Connection) {
	//TODO: handle
	fmt.Println("Entity action ", ea)

	if ea.ActionID == 0 || ea.ActionID == 1 {
		fmt.Println("Sneaky unsneaky", ea.ActionID)
		currentEffect := byte(0x00)
		pose := 0

		if ea.ActionID == 0 {
			currentEffect = 0x02 //Crouching
			pose = 5             //Sneaking
		}

		updateMetadata := Packet(&EntityMetadata{
			EntityID: connection.Player.EntityID,
			Metadata: &dataTypes.EntityMetadata{
				Effect: &currentEffect,
				Pose:   &pose,
			},
		})

		connection.Minecraft.SendToAllExcept(connection, &updateMetadata)
	}

}
