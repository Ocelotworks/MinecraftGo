package packet

type BlockChange struct {
	Location int64 `proto:"long"`
	BlockID  int   `proto:"varInt"`
}

func (bc *BlockChange) GetPacketId() int {
	return 0x0C
}
