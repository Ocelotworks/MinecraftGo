package packet

type PlayerPosition struct {
	X        float64 `proto:"double"`
	FeetY    float64 `proto:"double"`
	Z        float64 `proto:"double"`
	OnGround bool    `proto:"bool"`
}

func (pp *PlayerPosition) GetPacketId() int {
	return 0x36
}
