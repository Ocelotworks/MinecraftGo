package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerPositionAndLook struct {
	TeleportID int     `proto:"varInt"`
	X          float64 `proto:"double"`
	Y          float64 `proto:"double"`
	Z          float64 `proto:"double"`
	VelX       float64 `proto:"double"`
	VelY       float64 `proto:"double"`
	VelZ       float64 `proto:"double"`
	Yaw        float32 `proto:"float"`
	Pitch      float32 `proto:"float"`
	Flags      int     `proto:"int"`
}

func (ppal *PlayerPositionAndLook) GetPacketId() int {
	return constants.CBPlayerPosition
}

/**
func (ppal *PlayerPositionAndLook) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Position And Look ", ppal)

}
*/
