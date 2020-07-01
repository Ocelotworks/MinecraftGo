package packet

type SpawnPosition struct {
	Location int64 `proto:"long"`
}

func (sp *SpawnPosition) GetPacketId() int {
	return 0x42
}

/**
func (sp *SpawnPosition) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
