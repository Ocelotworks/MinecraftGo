package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type PlayerDigging struct {
	CurrentPacket *packetType.PlayerDigging
}

func (pd *PlayerDigging) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerDigging{}
}

func (pd *PlayerDigging) Init(currentPacket packetType.Packet) {
	pd.CurrentPacket = currentPacket.(*packetType.PlayerDigging)
}

func (pd *PlayerDigging) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Digging", pd)

	if pd.CurrentPacket.Status == 2 {
		acknowledge := packetType.Packet(&packetType.AcknowledgePlayerDigging{
			Location:   pd.CurrentPacket.Location,
			Block:      0,
			Status:     pd.CurrentPacket.Status,
			Successful: true,
		})
		connection.SendPacket(&acknowledge)

		blockChangePacket := packetType.Packet(&packetType.BlockChange{
			Location: pd.CurrentPacket.Location,
			BlockID:  0,
		})

		connection.Minecraft.SendToAllInPlay(&blockChangePacket)
	}

}
