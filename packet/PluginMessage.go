package packet

type PluginMessage struct {
	IsServer   bool
	Identifier string `proto:"string"`
	ByteArray  []byte `proto:"raw"`
}

func (pm *PluginMessage) GetPacketId() int {
	if pm.IsServer {
		return 0x0B
	}
	return 0x19
}

/**
func (pm *PluginMessage) Handle(packet []byte, connection *Connection) {
	fmt.Println("We need to handle this ", pm.Identifier)
}
*/
