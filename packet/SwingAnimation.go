package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type SwingAnimation struct {
	EntityID      int                          `proto:"varInt"`
	Hand          constants.Hand               `proto:"varInt"`
	AnimationType constants.SwingAnimationType `proto:"varInt"`
	Duration      int                          `proto:"varInt"`
}

func (ea *SwingAnimation) GetPacketId() int {
	return constants.CBSwingAnimation
}
