package packet

import "fmt"

type StatusPing struct {
	Payload int64 `proto:"long"`
}

func (sp *StatusPing) GetPacketId() int {
	return 0x01
}

func (sp *StatusPing) Handle(packet []byte, connection *Connection) {
	//Just send the pong right back
	fmt.Println("Pingy pongu")
	returnPacket := Packet(sp)
	connection.SendPacket(&returnPacket)
}
