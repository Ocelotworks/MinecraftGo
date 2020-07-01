package packet

type ServerPlayerAbilities struct {
	Flags byte `proto:"unsignedByte"`
}

func (spb *ServerPlayerAbilities) GetPacketId() int {
	return 0x1A
}
