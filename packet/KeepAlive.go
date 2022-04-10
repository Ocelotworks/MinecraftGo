package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type KeepAlive struct {
	ID int64 `proto:"long"`
}

func (ka *KeepAlive) GetPacketId() int {
	return constants.CBKeepAlive //Client
}

/**
func (ka *KeepAlive) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("KeepAlive", ka)
}
*/
