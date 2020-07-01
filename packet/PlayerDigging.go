package packet

type PlayerDigging struct {
	Status   int   `proto:"varInt"`
	Location int64 `proto:"long"`
	Face     byte  `proto:"unsignedByte"`
}

func (pd *PlayerDigging) GetPacketId() int {
	return 0x1B
}
