package packet

type Animation struct {
	Hand int `proto:"varInt"`
}

func (a *Animation) GetPacketId() int {
	return 0x2A
}

/**
func (a *Animation) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Animation", a)

	action := byte(0)

	if a.Hand == 1 {
		action = 3
	}

	animation := Packet(&EntityAnimation{
		EntityID:  connection.Player.EntityID,
		Animation: action,
	})

	connection.Minecraft.SendToAllExcept(connection, &animation)
}
*/
