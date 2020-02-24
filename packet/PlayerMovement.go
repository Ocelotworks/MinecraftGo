package packet

import "fmt"

type PlayerMovement struct {
	OnGround bool `proto:"bool"`
}

func (pm *PlayerMovement) GetPacketId() int {
	return 0x14
}

func (pm *PlayerMovement) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Movement", pm)
}
