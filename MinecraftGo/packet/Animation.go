package packet

import "fmt"

type Animation struct {
	Hand int `proto:"varInt"`
}

func (a *Animation) GetPacketId() int {
	return 0x1A
}

func (a *Animation) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Animation", a)
}
