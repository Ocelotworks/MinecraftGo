package packet

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
)

type RegistryData struct {
	RegistryID  string `proto:"string"`
	EntryLength int    `proto:"varInt"`
	NBTBytes    []byte `proto:"raw"`
}

func (apd *RegistryData) GetPacketId() int {
	return constants.CBRegistryDataConfiguration
}
