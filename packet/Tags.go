package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type Tags struct {
	TagCount int    `proto:"varInt"`
	Tags     []byte `proto:"raw"`
}

func (t *Tags) GetPacketId() int {
	return constants.CBTags
}
