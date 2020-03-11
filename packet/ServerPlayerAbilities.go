package packet

type ServerPlayerAbilities struct {
	Flags        byte    `proto:"unsignedByte"`
	FlyingSpeed  float32 `proto:"float"`
	WalkingSpeed float32 `proto:"float"`
}

func (spb *ServerPlayerAbilities) GetPacketId() int {
	return 0x19
}
