package packet

type Animation struct {
	Hand int `proto:"varInt"`
}

func (a *Animation) GetPacketId() int {
	return 0x2A
}
