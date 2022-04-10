package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type BlockChange struct {
	Location int64 `proto:"long"`
	BlockID  int   `proto:"varInt"`
}

func (bc *BlockChange) GetPacketId() int {
	return constants.CBBlockChange
}
