package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerMovement struct {
	OnGround bool `proto:"bool"`
}

func (pm *PlayerMovement) GetPacketId() int {
	return constants.SBPlayerMovement
}

/**
func (pm *PlayerMovement) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Movement", pm)
}
*/
