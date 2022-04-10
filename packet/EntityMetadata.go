package packet

import (
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
)

type EntityMetadata struct {
	EntityID int                `proto:"varInt"`
	Metadata dataTypes.Metadata `proto:"entityMetadata"`
}

func (em *EntityMetadata) GetPacketId() int {
	return constants.CBEntityMetadata
}

/**
func (em *EntityMetadata) Handle(packet []byte, connection *Connection) {
	//Client Only
}
*/
