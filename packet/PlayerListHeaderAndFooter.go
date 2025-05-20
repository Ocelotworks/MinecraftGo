package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type PlayerListHeaderAndFooter struct {
	Header string `proto:"string"`
	Footer string `proto:"string"`
}

func (plhaf *PlayerListHeaderAndFooter) GetPacketId() int {
	return constants.CBTabListHeaderAndFooter
}
