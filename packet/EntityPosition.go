package packet

type EntityPosition struct {
	EntityID int   `proto:"varInt"`
	DeltaX   int16 `proto:"short"`
	DeltaY   int16 `proto:"short"`
	DeltaZ   int16 `proto:"short"`
	OnGround bool  `proto:"bool"`
}

func (ep *EntityPosition) GetPacketId() int {
	return 0x28
}

/**
func (ep *EntityPosition) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
