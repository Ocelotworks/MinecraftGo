package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityAnimation struct {
	EntityID  int  `proto:"varInt"`
	Animation byte `proto:"unsignedByte"`
}

func (ea *EntityAnimation) GetPacketId() int {
	return constants.CBEntityAnimation
}

/**
func (ea *EntityAnimation) Handle(packet []byte, connection *Connection) {
	//Client only packet
}
*/
