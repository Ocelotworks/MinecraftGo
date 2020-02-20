package packet

type Packet interface {
	Handle(packet []byte, connection *Connection)
}
