package packet

type Packet interface {
	GetPacketId() int
	Handle(packet []byte, connection *Connection)
}
