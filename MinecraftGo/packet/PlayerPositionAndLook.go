package packet

import "fmt"

type PlayerPositionAndLook struct {
	X          float64 `proto:"double"`
	Y          float64 `proto:"double"`
	Z          float64 `proto:"double"`
	Yaw        float32 `proto:"float"`
	Pitch      float32 `proto:"float"`
	Flags      byte    `proto:"unsignedByte"`
	TeleportID int     `proto:"varInt"`
}

func (ppal *PlayerPositionAndLook) GetPacketId() int {
	return 0x36
}

func (ppal *PlayerPositionAndLook) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Position And Look ", ppal)

}
