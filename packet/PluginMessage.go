package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PluginMessage struct {
	IsServer   bool
	Identifier string `proto:"string"`
	ByteArray  []byte `proto:"raw"`
}

func (pm *PluginMessage) GetPacketId() int {
	if pm.IsServer {
		return constants.CBPluginMessage
	}
	return constants.SBPluginMessage
}

/**
func (pm *PluginMessage) Handle(packet []byte, connection *Connection) {
	fmt.Println("We need to handle this ", pm.Identifier)
}
*/
