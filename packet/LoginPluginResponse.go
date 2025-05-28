package packet

type LoginPluginResponse struct {
	MessageId int `proto:"varInt"`
	Data      []byte
}

func (ls *LoginPluginResponse) GetPacketId() int {
	return 0x02
}

/**
func (ls *LoginSuccess) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
