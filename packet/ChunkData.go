package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type ChunkData struct {
	X                    int     `proto:"int"`
	Z                    int     `proto:"int"`
	HeightMap            []byte  `proto:"raw"`
	DataSize             int     `proto:"varInt"`
	Data                 []byte  `proto:"raw"`
	BlockEntityCount     int     `proto:"varInt"`
	BlockEntities        []byte  `proto:"raw"`
	TrustEdges           bool    `proto:"bool"`
	SkyLightMask         []int64 `proto:"bitset"`
	BlockLightMask       []int64 `proto:"bitset"`
	EmptySkyLightMask    []int64 `proto:"bitset"`
	EmptyBlockLightMask  []int64 `proto:"bitset"`
	SkyLightArrayCount   int     `proto:"varInt"`
	SkyLightArrays       []byte  `proto:"raw"`
	BlockLightArrayCount int     `proto:"varInt"`
	BlockLightArrays     []byte  `proto:"raw"`
}

func (cd *ChunkData) GetPacketId() int {
	return constants.CBChunkDataAndUpdateLight
}

/**
func (cd *ChunkData) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
