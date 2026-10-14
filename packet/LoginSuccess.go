package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type LoginSuccess struct {
	UUID       []byte          `proto:"uuid"`
	Username   string          `proto:"string"`
	Properties []LoginProperty `proto:"prefixedArray"`
	SessionID  []byte          `proto:"uuid"`
}

type LoginProperty struct {
	Key   string `proto:"string"`
	Value string `proto:"string"`
}

func (ls *LoginSuccess) GetPacketId() int {
	return constants.CBLoginFinished
}

/**
func (ls *LoginSuccess) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
