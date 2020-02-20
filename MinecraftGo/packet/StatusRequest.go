package packet

import (
	"encoding/json"
	"fmt"

	"../entity"
)

type StatusRequest struct {
}

func (sr *StatusRequest) GetPacketId() int {
	return 0x01
}
func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response

	status := entity.ServerListPingResponse{
		Version: entity.ServerListPingVersion{},
		Players: entity.ServerListPingPlayers{
			Max:    420,
			Online: 69,
			Sample: []entity.ServerListPingPlayerListItem{{
				Name: "UnacceptableUse",
				ID:   "5d8af060-129e-419c-b3ac-c0dad1c91181",
			}},
		},
		Favicon: "",
		Description: entity.ChatMessageComponent{
			Text: "Hello World!",
		},
	}

	output, exception := json.Marshal(status)

	if exception != nil {
		fmt.Println("Exception encoding server list response:", exception)
		return
	}

	response := Packet(&StatusResponse{Status: string(output)})
	connection.SendPacket(&response)

}
