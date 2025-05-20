package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ServerDifficulty struct {
	Difficulty       byte `proto:"unsignedByte"`
	DifficultyLocked bool `proto:"bool"`
}

func (sd *ServerDifficulty) GetPacketId() int {
	return constants.CBChangeDifficulty
}

/**
func (sd *ServerDifficulty) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
