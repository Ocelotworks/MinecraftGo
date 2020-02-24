package packet

import "fmt"

type KeepAlive struct {
	ID int64 `proto:"long"`
}

func (ka *KeepAlive) GetPacketId() int {
	return 0x21 //Client
}

func (ka *KeepAlive) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("KeepAlive", ka)
}
