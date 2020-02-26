package controller

import (
	"encoding/json"
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type StatusRequest struct {
	CurrentPacket *packetType.StatusRequest
}

func (sr *StatusRequest) GetPacketStruct() packetType.Packet {
	return &packetType.StatusRequest{}
}

func (sr *StatusRequest) Init(currentPacket packetType.Packet) {
	sr.CurrentPacket = currentPacket.(*packetType.StatusRequest)
}

func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response
	fmt.Println("Status Request")
	status := entity.ServerListPingResponse{
		Version: entity.ServerListPingVersion{
			Name:     "1.15.2",
			Protocol: 578,
		},
		Players: entity.ServerListPingPlayers{
			Max:    connection.Minecraft.MaxPlayers,
			Online: connection.Minecraft.ConnectedPlayers,
			Sample: []entity.ServerListPingPlayerListItem{{
				Name: "UnacceptableUse",
				ID:   "5d8af060-129e-419c-b3ac-c0dad1c91181",
			}},
		},
		Description: connection.Minecraft.ServerName,
	}

	output, exception := json.Marshal(status)

	if exception != nil {
		fmt.Println("Exception encoding server list response:", exception)
		return
	}

	response := packetType.Packet(&packetType.StatusResponse{Status: string(output)})
	connection.SendPacket(&response)
}
