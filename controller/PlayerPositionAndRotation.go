package controller

import packetType "github.com/Ocelotworks/MinecraftGo/packet"

type PlayerPositionAndRotation struct {
	CurrentPacket *packetType.PlayerPositionAndRotation
}

func (ppar *PlayerPositionAndRotation) GetPacketStruct() packetType.Packet {
	return &packetType.PlayerPositionAndRotation{}
}

func (ppar *PlayerPositionAndRotation) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	ppar.CurrentPacket = currentPacket.(*packetType.PlayerPositionAndRotation)
}

func (ppar *PlayerPositionAndRotation) Handle(packet []byte, connection *Connection) {
	connection.Minecraft.UpdatePlayerPosition(connection, ppar.CurrentPacket.X, ppar.CurrentPacket.FeetY, ppar.CurrentPacket.Z, ppar.CurrentPacket.Yaw, ppar.CurrentPacket.Pitch)
}
