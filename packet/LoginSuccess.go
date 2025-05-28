package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type LoginSuccess struct {
	UUID             []byte `proto:"uuid"`
	Username         string `proto:"string"`
	PropertiesLength int    `proto:"varInt"`
	Properties       []LoginProperty
}

type LoginProperty struct {
	Key   string `proto:"string"`
	Value string `proto:"string"`
}

func (ls *LoginSuccess) GetPacketId() int {
	return constants.CBLoginSuccess
}

/**
func (ls *LoginSuccess) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
