package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type Disconnect struct {
	Reason string `proto:"string"`
}

func (d *Disconnect) GetPacketId() int {
	return constants.CBDisconnect
}

/**
func (d *Disconnect) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
