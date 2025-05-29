package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"

	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/gofrs/uuid"
)

type StatusRequest struct {
	CurrentPacket *packetType.StatusRequest
}

func (sr *StatusRequest) GetPacketStruct() packetType.Packet {
	return &packetType.StatusRequest{}
}

func (sr *StatusRequest) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	sr.CurrentPacket = currentPacket.(*packetType.StatusRequest)
}

func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response
	fmt.Println("Status Request")

	players := make([]entity.ServerListPingPlayerListItem, 0)

	for _, c := range connection.Minecraft.Connections {
		if c.Player != nil && c.Player.EntityID > 0 {
			playerUUID, exception := uuid.FromBytes(c.Player.UUID)

			if exception != nil {
				fmt.Println(exception)
				continue
			}

			players = append(players, entity.ServerListPingPlayerListItem{
				Name: c.Player.Username,
				ID:   playerUUID.String(),
			})
		}
	}

	status := entity.ServerListPingResponse{
		Version: entity.ServerListPingVersion{
			Name:     constants.GameVersion,
			Protocol: constants.ProtocolVersion,
		},
		Players: entity.ServerListPingPlayers{
			Max:    connection.Minecraft.MaxPlayers,
			Online: connection.Minecraft.ConnectedPlayers,
			Sample: players,
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
