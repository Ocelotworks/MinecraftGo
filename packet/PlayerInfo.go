package packet

import (
	"github.com/Ocelotworks/MinecraftGo/entity"
)

type PlayerInfoAddPlayer struct {
	Action  int             `proto:"varInt"`
	Players []entity.Player `proto:"playerArray"`
}

type PlayerInfoRemovePlayer struct {
	Action          int    `proto:"varInt"`
	NumberOfPlayers int    `proto:"varInt"`
	UUID            []byte `proto:"uuid"`
}

type PlayerInfoUpdatePing struct {
	Action          int    `proto:"varInt"`
	NumberOfPlayers int    `proto:"varInt"`
	UUID            []byte `proto:"uuid"`
	Ping            int    `proto:"varInt"`
}

// TODO: this is all completely different now
func (piap *PlayerInfoAddPlayer) GetPacketId() int {
	return 97
}

func (pirp *PlayerInfoRemovePlayer) GetPacketId() int {
	return 96
}

func (piup *PlayerInfoUpdatePing) GetPacketId() int {
	return 95
}

/**
func (piap *PlayerInfoAddPlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}

func (pirp *PlayerInfoRemovePlayer) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
