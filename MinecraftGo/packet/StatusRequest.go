package packet

import (
	"encoding/json"
	"fmt"
	"time"

	"../entity"
)

type StatusRequest struct {
}

func (sr *StatusRequest) GetPacketId() int {
	return 0x01
}
func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response
	fmt.Println("Status Request")
	red := entity.Red
	status := entity.ServerListPingResponse{
		Version: entity.ServerListPingVersion{
			Name:     "1.15.2",
			Protocol: 578,
		},
		Players: entity.ServerListPingPlayers{
			Max:    420,
			Online: 69,
			Sample: []entity.ServerListPingPlayerListItem{{
				Name: "UnacceptableUse",
				ID:   "5d8af060-129e-419c-b3ac-c0dad1c91181",
			}},
		},
		Description: entity.ChatMessageComponent{
			Text:   time.Now().Format(time.ANSIC),
			Colour: &red,
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
