package packet

import "github.com/Ocelotworks/MinecraftGo/entity"

type PlayerInfoAddPlayer struct {
	Action  int             `proto:"varInt"`
	Players []entity.Player `proto:"playerArray"`
}

type PlayerInfoRemovePlayer struct {
	Action          int    `proto:"varInt"`
	NumberOfPlayers int    `proto:"varInt"`
	UUID            []byte `proto:"uuid"`
}

func (piap *PlayerInfoAddPlayer) GetPacketId() int {
	return 0x34
}

func (pirp *PlayerInfoRemovePlayer) GetPacketId() int {
	return 0x34
}

/**
func (piap *PlayerInfoAddPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}

func (pirp *PlayerInfoRemovePlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
