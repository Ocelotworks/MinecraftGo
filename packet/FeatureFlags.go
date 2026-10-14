package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type FeatureFlags struct {
	FeatureFlags []string `proto:"stringArray"`
}

func (ff *FeatureFlags) GetPacketId() int {
	return constants.CBUpdateEnabledFeaturesConfiguration
}
