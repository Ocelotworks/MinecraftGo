package packet

import "github.com/Ocelotworks/MinecraftGo/dataTypes"

type EntityMetadata struct {
	EntityID int                `proto:"varInt"`
	Metadata dataTypes.Metadata `proto:"entityMetadata"`
}

func (em *EntityMetadata) GetPacketId() int {
	return 0x44
}

func (em *EntityMetadata) Handle(packet []byte, connection *Connection) {
	//Client Only
}
