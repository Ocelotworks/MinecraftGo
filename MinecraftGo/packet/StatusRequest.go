package packet

type StatusRequest struct {
}

func (sr *StatusRequest) GetPacketId() int {
	return 0x01
}
func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response
}
