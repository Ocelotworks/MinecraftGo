package controller

import (
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type ClientInformation struct {
	CurrentPacket *packetType.ClientInformation
}

func (ci *ClientInformation) GetPacketStruct() packetType.Packet {
	return &packetType.ClientInformation{}
}

func (ci *ClientInformation) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	ci.CurrentPacket = currentPacket.(*packetType.ClientInformation)
}

func (ci *ClientInformation) Handle(packet []byte, connection *Connection) {

	fmt.Println("Client Information ", ci.CurrentPacket)

	// This is the only client information the client is guaranteed to send:
	// ClientHandshakePacketListenerImpl sends it once on acknowledging login,
	// and after that only Options.save resends it, and only when something
	// actually changed.  So the view distance the join reads has to be taken
	// here rather than waiting for a copy in the play state that may never
	// arrive.
	connection.Player.Settings = entity.PlayerSettings{
		Locale:             ci.CurrentPacket.Locale,
		ViewDistance:       ci.CurrentPacket.ViewDistance,
		ChatMode:           ci.CurrentPacket.ChatMode,
		ChatColours:        ci.CurrentPacket.ChatColors,
		DisplayedSkinParts: ci.CurrentPacket.DisplayedSkinParts,
		MainHand:           ci.CurrentPacket.MainHand,
	}
}
