package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type TeleportConfirm struct {
	TeleportID int `proto:"varInt"`
}

func (tc *TeleportConfirm) GetPacketId() int {
	return constants.SBTeleportConfirm
}

/**
func (tc *TeleportConfirm) Handle(packet []byte, connection *Connection) {
	fmt.Println("Teleport confirm ", tc)
}
*/
