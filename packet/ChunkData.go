package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type HeightMapOuter struct {
	Inner HeightMap `nbt:""`
}

type HeightMap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	//WORLD_SURFACE
}

type ChunkData struct {
	X                    int            `proto:"int"`
	Z                    int            `proto:"int"`
	HeightMap            HeightMapOuter `proto:"nbt"`
	DataSize             int            `proto:"varInt"`
	Data                 []byte         `proto:"raw"`
	BlockEntityCount     int            `proto:"varInt"`
	BlockEntities        []byte         `proto:"raw"`
	TrustEdges           bool           `proto:"bool"`
	SkyLightMask         []byte         `proto:"raw"`
	BlockLightMask       []byte         `proto:"raw"`
	EmptySkyLightMask    []byte         `proto:"raw"`
	EmptyBlockLightMask  []byte         `proto:"raw"`
	SkyLightArrayCount   int            `proto:"varInt"`
	SkyLightArrays       []byte         `proto:"raw"`
	BlockLightArrayCount int            `proto:"varInt"`
	BlockLightArrays     []byte         `proto:"raw"`
}

func (cd *ChunkData) GetPacketId() int {
	return constants.CBChunkDataAndUpdateLight
}

/**
func (cd *ChunkData) Handle(packet []byte, connection *Connection) {
	//Client Only Packet
}
*/
