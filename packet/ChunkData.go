package packet

type ChunkData struct {
	X                int    `proto:"int"`
	Z                int    `proto:"int"`
	FullChunk        bool   `proto:"bool"`
	PrimaryBitMask   int    `proto:"varInt"`
	HeightMap        []byte `proto:"raw"`
	Biomes           []int  `proto:"intArray"`
	DataSize         int    `proto:"varInt"`
	Data             []byte `proto:"raw"`
	BlockEntityCount int    `proto:"varInt"`
	BlockEntities    []byte `proto:"raw"`
}

func (cd *ChunkData) GetPacketId() int {
	return 0x21
}

/**
func (cd *ChunkData) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
