package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerCommand struct {
	EntityID  int `proto:"varInt"`
	ActionID  int `proto:"varInt"`
	JumpBoost int `proto:"varInt"`
}

func (ea *PlayerCommand) GetPacketId() int {
	return constants.SBPlayerCommand
}
