package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

//TODO
type PlayerAbilities struct {
}

func (pa *PlayerAbilities) GetPacketId() int {
	return constants.CBPlayerAbilities
}

/**
func (pa *PlayerAbilities) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
