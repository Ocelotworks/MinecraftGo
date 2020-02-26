package controller

import packetType "github.com/Ocelotworks/MinecraftGo/packet"

type PlayerRotation struct {
	CurrentPacket *packetType.PlayerRotation
}

func (pr *PlayerRotation) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerRotation{}
}

func (pr *PlayerRotation) Init(currentPacket packetType.Packet) {
	pr.CurrentPacket = currentPacket.(*packetType.PlayerRotation)
}

func (pr *PlayerRotation) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, 0, 0, 0, pr.CurrentPacket.Yaw, pr.CurrentPacket.Pitch)
}
