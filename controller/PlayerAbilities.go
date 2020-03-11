package controller

import packetType "github.com/Ocelotworks/MinecraftGo/packet"

type PlayerAbilities struct {
	CurrentPacket *packetType.ServerPlayerAbilities
}

func (pa *PlayerAbilities) GetPacketStruct() packetType.Packet {
	return &packetType.ServerPlayerAbilities{}
}

func (pa *PlayerAbilities) Init(currentPacket packetType.Packet) {
	pa.CurrentPacket = currentPacket.(*packetType.ServerPlayerAbilities)
}

func (pa *PlayerAbilities) Handle(packet []byte, connection *Connection) {
	//TODO
}
