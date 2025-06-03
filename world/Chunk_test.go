package world

import (
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseChunkFromRegionChunk(t *testing.T) {
	inData, err := os.ReadFile("C:\\Users\\unacc\\IdeaProjects\\MinecraftGo\\data\\nbt-test\\region.mca")

	assert.NoError(t, err)

	region := dataTypes.ReadRegionFile(inData)

	chunk := ParseChunkFromRegionChunk(region.Chunks[0])

	assert.Equal(t, "minecraft:andesite", chunk.Blocks[15][5][15].Type)
	assert.Equal(t, "minecraft:andesite", chunk.Blocks[0][5][0].Type)
	assert.Equal(t, "minecraft:acacia_log", chunk.Blocks[0][0][0].Type)
	assert.Equal(t, "minecraft:basalt", chunk.Blocks[8][16][8].Type)

	//for x := int32(0); x < 16; x++ {
	//    for y := int32(0); y < 16*16; y++ {
	//        for z := int32(0); z < 16; z++ {
	//            if chunk.Blocks[x][y][z].Type == "minecraft:andesite" {
	//                fmt.Println("Found andesite at ", x, y, z)
	//            }
	//            if chunk.Blocks[x][y][z].Type == "minecraft:basalt" {
	//                fmt.Println("Found basalt at ", x, y, z)
	//            }
	//        }
	//    }
	//}
}
