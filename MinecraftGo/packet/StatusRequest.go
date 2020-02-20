package packet

type StatusRequest struct {
}

func (sr *StatusRequest) Handle(packet []byte, connection *Connection) {
	//sends the client response
}
