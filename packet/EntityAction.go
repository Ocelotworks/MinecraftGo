package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type EntityAction struct {
	EntityID  int `proto:"varInt"`
	ActionID  int `proto:"varInt"`
	JumpBoost int `proto:"varInt"`
}

func (ea *EntityAction) GetPacketId() int {
	return constants.SBEntityAction
}
