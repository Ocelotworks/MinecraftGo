package packet

import "fmt"

type PlayerDigging struct {
	Status   int   `proto:"varInt"`
	Location int64 `proto:"long"`
	Face     byte  `proto:"unsignedByte"`
}

func (pd *PlayerDigging) GetPacketId() int {
	return 0x1A
}

func (pd *PlayerDigging) Handle(packet []byte, connection *Connection) {
	//TODO: Handle
	fmt.Println("Player Digging", pd)

	if pd.Status == 2 {
		acknowledge := Packet(&AcknowledgePlayerDigging{
			Location:   pd.Location,
			Block:      0,
			Status:     pd.Status,
			Successful: true,
		})
		connection.SendPacket(&acknowledge)
	}

}
