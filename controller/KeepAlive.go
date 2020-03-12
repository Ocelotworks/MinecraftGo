package controller

import (
	"fmt"
	"time"

	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type KeepAlive struct {
	CurrentPacket *packetType.KeepAlive
}

func (ka *KeepAlive) GetPacketStruct() packetType.Packet {
	return &packetType.KeepAlive{}
}

func (ka *KeepAlive) Init(currentPacket packetType.Packet) {
	ka.CurrentPacket = currentPacket.(*packetType.KeepAlive)
}

func (ka *KeepAlive) Handle(packet []byte, connection *Connection) {
	now := time.Now().Unix()
	ping := int(now - ka.CurrentPacket.ID)

	if ping != connection.Ping {
		updatePing := packetType.Packet(&packetType.PlayerInfoUpdatePing{
			Action:          2,
			NumberOfPlayers: 1,
			UUID:            connection.Player.UUID,
			Ping:            ping,
		})
		go connection.Minecraft.SendToAllInPlay(&updatePing)
		connection.Ping = ping
	}
	fmt.Println("KeepAlive ping", ping)
}
