package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ClientTickEnd struct {
}

func (bc *ClientTickEnd) GetPacketId() int {
	return constants.SBClientTickEnd
}
