package packet

import "github.com/Ocelotworks/MinecraftGo/constants"

type HeightMap struct {
	Type int `proto:"varInt"`
	//WORLD_SURFACE
}

// ChunkData is ClientboundLevelChunkWithLightPacket: the chunk itself followed
// by its light data.  The fields below are in wire order.
//
// Several fields are pre-encoded raw bytes because they are variable-length
// structures the struct scanner has no type for:
//
//   - HeightMapData is a map of heightmap type to long array — a VarInt count
//     then that many (VarInt type, VarInt length, longs) entries.  A count of
//     zero sends no heightmaps, which the client accepts.
//   - The four light masks are each a BitSet, sent as a VarInt long count
//     followed by that many longs.
//   - The light arrays are each a VarInt length followed by 2048 bytes, one
//     nibble per block.
type ChunkData struct {
	X                    int    `proto:"int"`
	Z                    int    `proto:"int"`
	HeightMapData        []byte `proto:"raw"`
	DataSize             int    `proto:"varInt"`
	Data                 []byte `proto:"raw"`
	BlockEntityCount     int    `proto:"varInt"`
	BlockEntities        []byte `proto:"raw"`
	SkyLightMask         []byte `proto:"raw"`
	BlockLightMask       []byte `proto:"raw"`
	EmptySkyLightMask    []byte `proto:"raw"`
	EmptyBlockLightMask  []byte `proto:"raw"`
	SkyLightArrayCount   int    `proto:"varInt"`
	SkyLightArrays       []byte `proto:"raw"`
	BlockLightArrayCount int    `proto:"varInt"`
	BlockLightArrays     []byte `proto:"raw"`
}

func (cd *ChunkData) GetPacketId() int {
	return constants.CBLevelChunkWithLight
}
