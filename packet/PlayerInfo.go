package packet

import "github.com/Ocelotworks/MinecraftGo/entity"

type PlayerInfoAddPlayer struct {
	Action  int             `proto:"varInt"`
	Players []entity.Player `proto:"playerArray"`
}

func (piap *PlayerInfoAddPlayer) GetPacketId() int {
	return 0x34
}

func (piap *PlayerInfoAddPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
