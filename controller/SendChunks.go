package controller

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/Ocelotworks/MinecraftGo/world"
)

// maxSendViewDistance caps how far around a player chunks are sent regardless
// of what the client asks for.  Each chunk carries a full set of sky light
// arrays, so an unbounded view distance is a lot of traffic on a connection
// with compression disabled.
const maxSendViewDistance = 8

// lightSectionCount is the number of sections the light engine covers: one per
// chunk section plus one either side of the world, matching
// LevelLightEngine.getLightSectionCount().
func lightSectionCount(height world.WorldHeight) int {
	return height.SectionCount() + 2
}

// fullSkyLightArray is one section's worth of sky light at full brightness.
// Light is stored as a nibble per block, so 4096 blocks fit in 2048 bytes and
// 0xFF sets both nibbles to 15.
func fullSkyLightArray() []byte {
	data := make([]byte, 2048)
	for i := range data {
		data[i] = 0xFF
	}
	return data
}

// fullMask returns the longs of a BitSet with the low count bits set, which is
// what dataTypes.WriteBitSet takes.
func fullMask(count int) []int64 {
	longs := make([]int64, (count+63)/64)
	for i := 0; i < count; i++ {
		longs[i/64] |= 1 << uint(i%64)
	}
	return longs
}

// chunkLight is the light half of a chunk packet.  It is the same for every
// chunk: every section gets full sky light and no block light at all.  Without
// this the client renders the world at light level zero.
//
// The masks say which sections carry data and which are known to be empty; a
// section named in neither is left at whatever the client has.  There is one
// array per bit set in the corresponding mask.
type chunkLight struct {
	allSections    []byte
	noSections     []byte
	sectionCount   int
	skyLightArrays []byte
}

func newChunkLight(height world.WorldHeight) chunkLight {
	sections := lightSectionCount(height)

	skyLightArrays := make([]byte, 0, sections*(2048+2))
	for i := 0; i < sections; i++ {
		skyLightArrays = append(skyLightArrays, dataTypes.WriteVarInt(2048)...)
		skyLightArrays = append(skyLightArrays, fullSkyLightArray()...)
	}

	return chunkLight{
		allSections:    dataTypes.WriteBitSet(fullMask(sections)),
		noSections:     dataTypes.WriteBitSet([]int64{}),
		sectionCount:   sections,
		skyLightArrays: skyLightArrays,
	}
}

// newChunkPacket builds ClientboundLevelChunkWithLightPacket for a chunk whose
// sections have already been serialised into chunkData.
func newChunkPacket(x int, z int, chunkData []byte, light chunkLight) *packetType.ChunkData {
	return &packetType.ChunkData{
		X: x,
		Z: z,
		// An empty heightmap map: the client applies whatever it is
		// given and primes the rest itself.
		HeightMapData:        dataTypes.WriteVarInt(0),
		DataSize:             len(chunkData),
		Data:                 chunkData,
		BlockEntityCount:     0,
		BlockEntities:        []byte{},
		SkyLightMask:         light.allSections,
		BlockLightMask:       light.noSections,
		EmptySkyLightMask:    light.noSections,
		EmptyBlockLightMask:  light.allSections,
		SkyLightArrayCount:   light.sectionCount,
		SkyLightArrays:       light.skyLightArrays,
		BlockLightArrayCount: 0,
		BlockLightArrays:     []byte{},
	}
}

// SendChunks sends the chunks around a player, as
// ClientboundLevelChunkWithLightPacket does.
func (mc *Minecraft) SendChunks(connection *Connection) {
	viewDistance := int(connection.Player.Settings.ViewDistance)
	if viewDistance <= 0 || viewDistance > maxSendViewDistance {
		viewDistance = maxSendViewDistance
	}

	centreX := int32(mc.BlockCoordToChunkCoord(connection.Player.X))
	centreZ := int32(mc.BlockCoordToChunkCoord(connection.Player.Z))

	light := newChunkLight(mc.ChunkManager.Height)

	for _, chunk := range mc.ChunkManager.ChunksAround(centreX, centreZ, viewDistance) {
		packet := packetType.Packet(newChunkPacket(int(chunk.X), int(chunk.Z), chunk.WriteNetworkData(mc.WorldRegistry), light))
		connection.SendPacket(&packet)
	}
}
