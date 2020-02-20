package packet

type StatusPing struct {
	Payload int `proto:"long"`
}

func (sp *StatusPing) Handle(packet []byte, connection *Connection) {
	//sends the pong
}
