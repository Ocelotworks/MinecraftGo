package packet

type EntityAction struct {
	EntityID  int `proto:"varInt"`
	ActionID  int `proto:"varInt"`
	JumpBoost int `proto:"varInt"`
}

func (ea *EntityAction) GetPacketId() int {
	return 0x1C
}
