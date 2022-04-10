package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type AcknowledgePlayerDigging struct {
	Location   int64 `proto:"long"`
	Block      int   `proto:"varInt"`
	Status     int   `proto:"varInt"`
	Successful bool  `proto:"bool"`
}

func (apd *AcknowledgePlayerDigging) GetPacketId() int {
	return constants.CBAcknowledgePlayerDigging
}

/**
func (apd *AcknowledgePlayerDigging) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
