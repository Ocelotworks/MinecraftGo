package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerPositionAndLook struct {
	X               float64 `proto:"double"`
	Y               float64 `proto:"double"`
	Z               float64 `proto:"double"`
	Yaw             float32 `proto:"float"`
	Pitch           float32 `proto:"float"`
	Flags           byte    `proto:"unsignedByte"`
	TeleportID      int     `proto:"varInt"`
	DismountVehicle bool    `proto:"bool"`
}

func (ppal *PlayerPositionAndLook) GetPacketId() int {
	return constants.CBSyncPlayerPosition
}

/**
func (ppal *PlayerPositionAndLook) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Position And Look ", ppal)

}
*/
