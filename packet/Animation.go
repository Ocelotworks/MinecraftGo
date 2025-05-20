package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type Animation struct {
	Hand int `proto:"varInt"`
}

func (a *Animation) GetPacketId() int {
	return constants.SBSwingArm
}
