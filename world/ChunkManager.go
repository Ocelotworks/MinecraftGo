package world

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"os"
)

type ChunkManager struct {
	Chunks map[int32]map[int32]*Chunk
}

func NewChunkManager() *ChunkManager {
	return &ChunkManager{
		Chunks: make(map[int32]map[int32]*Chunk),
	}

}

func (cm *ChunkManager) loadRegionFile(x int, y int) *dataTypes.RegionMetadata {
	inData, exception := os.ReadFile(fmt.Sprintf("C:\\Users\\unacc\\IdeaProjects\\MinecraftGo\\data\\worlds\\MCGO_FlatTest\\region\\r.%d.%d.mca", x, y))

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return nil
	}

	region := dataTypes.ReadRegionFile(inData)
	return &region
}

func (cm *ChunkManager) LoadRegion(x int, y int) {

	regionMetadata := cm.loadRegionFile(x, y)
	cm.Chunks = make(map[int32]map[int32]*Chunk)

	for _, chunk := range regionMetadata.Chunks {
		if chunk == nil {
			continue
		}
		_, ok := cm.Chunks[chunk.XPos]
		if !ok {
			cm.Chunks[chunk.XPos] = make(map[int32]*Chunk)
		}

		cm.Chunks[chunk.XPos][chunk.ZPos] = ParseChunkFromRegionChunk(chunk)

	}

}
