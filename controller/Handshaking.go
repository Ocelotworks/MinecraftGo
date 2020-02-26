package controller

import (
	"fmt"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type Handshaking struct {
	CurrentPacket *packetType.Handshaking
}

func (h *Handshaking) GetPacketStruct() packetType.Packet {
	return &packetType.Handshaking{}
}

func (h *Handshaking) Init(currentPacket packetType.Packet) {
	h.CurrentPacket = currentPacket.(*packetType.Handshaking)
}

func (h *Handshaking) Handle(packet []byte, connection *Connection) {
	//fmt.Println(hex.Dump(packet))
	//fmt.Println("We are now handling a handshake packet")
	//
	//protocolVersion, cursor := dataTypes.ReadVarInt(packet)
	//h.ProtocolVersion = protocolVersion
	//fmt.Printf("Protocol Version %d (end %d)\n", protocolVersion, cursor)
	//
	//somethingElse, end := dataTypes.ReadVarInt(packet[cursor:])
	//cursor += end
	//fmt.Println("Mysterious value", somethingElse)
	//
	//serverAddress, end := dataTypes.ReadString(packet[cursor:])
	//cursor += end
	//h.ServerAddress = serverAddress
	//fmt.Printf("Server Address lives at %d: '%s'\n", cursor, h.ServerAddress)
	//
	//fmt.Println(packet[cursor:])
	//serverPort, end := dataTypes.ReadUnsignedShort(packet[cursor:])
	//cursor += end
	//fmt.Println("Server port", serverPort)
	//
	//nextState, end := dataTypes.ReadVarInt(packet[cursor:])
	//h.NextState = nextState
	//
	//fmt.Println(hex.Dump(packet))
	fmt.Printf("Connection to %s:%d Protocol Version %d\n", h.CurrentPacket.ServerAddress, h.CurrentPacket.ServerPort, h.CurrentPacket.ProtocolVersion)
	fmt.Println("Next State", h.CurrentPacket.NextState)
	connection.State = State(h.CurrentPacket.NextState)

}
