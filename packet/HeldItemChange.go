package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type HeldItemChange struct {
	Slot     byte `proto:"unsignedByte"`
	IsServer bool
}

func (hic *HeldItemChange) GetPacketId() int {
	if hic.IsServer {
		return constants.CBSetHeldItem
	}
	return constants.SBSetHeldItem
}

/**
func (hic *HeldItemChange) Handle(packet []byte, connection *Connection) {
	fmt.Println("Held Item Change ", hic)
}
*/
