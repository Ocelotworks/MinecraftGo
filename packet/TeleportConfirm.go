package packet

type TeleportConfirm struct {
	TeleportID int `proto:"varInt"`
}

func (tc *TeleportConfirm) GetPacketId() int {
	return 0x00
}

/**
func (tc *TeleportConfirm) Handle(packet []byte, connection *Connection) {
	fmt.Println("Teleport confirm ", tc)
}
*/
