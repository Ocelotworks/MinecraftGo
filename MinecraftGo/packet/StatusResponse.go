package packet

import "fmt"

type StatusResponse struct {
	Status string `proto:"string"`
}

func (sr *StatusResponse) GetPacketId() int {
	return 0x00
}

func (sr *StatusResponse) Handle(packet []byte, connection *Connection) {
	//We should never get this
	fmt.Println("Status Response??")
}
