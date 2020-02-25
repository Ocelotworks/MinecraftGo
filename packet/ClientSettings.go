package packet

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"

	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
)

type ClientSettings struct {
	Locale             string `proto:"string"`
	ViewDistance       byte   `proto:"unsignedByte"`
	ChatMode           int    `proto:"varInt"`
	ChatColours        bool   `proto:"bool"`
	DisplayedSkinParts byte   `proto:"unsignedByte"`
	MainHand           int    `proto:"varInt"`
}

func (cs *ClientSettings) GetPacketId() int {
	return 0x05
}

func (cs *ClientSettings) Handle(packet []byte, connection *Connection) {
	fmt.Println("Got Client Settings:")
	fmt.Println("Locale: ", cs.Locale)
	fmt.Println("View Distance: ", cs.ViewDistance)
	fmt.Println("Chat Mode: ", cs.ChatMode)
	fmt.Println("Chat Colours: ", cs.ChatColours)
	fmt.Println("Displayed skin: ", cs.DisplayedSkinParts)
	fmt.Println("Main Hand: ", cs.MainHand)

	connection.Player.Settings = entity.PlayerSettings{
		Locale:             cs.Locale,
		ViewDistance:       cs.ViewDistance,
		ChatMode:           cs.ChatMode,
		ChatColours:        cs.ChatColours,
		DisplayedSkinParts: cs.DisplayedSkinParts,
		MainHand:           cs.MainHand,
	}

	//inData, exception := ioutil.ReadFile("world/region/r.1.1.mca")
	//
	//if exception != nil {
	//	fmt.Println("Reading file")
	//	fmt.Println(exception)
	//	return
	//}
	//
	//region := dataTypes.ReadRegionFile(inData)
	//
	//chunk := region.Chunks[0]
	//
	//fmt.Println("Chunky chunk")
	//fmt.Println(chunk)

	randomBiomes := make([]int, 1024)
	for i := 0; i < 1024; i++ {
		randomBiomes[i] = i % 10
	}

	randomHeightMap := make([]int64, 36)
	for i := 0; i < 36; i++ {
		randomHeightMap[i] = int64(i)
	}

	heightMaps := dataTypes.NBTWrite(dataTypes.NBTNamed{
		Type: 10,
		Name: "",
		Data: []interface{}{
			dataTypes.NBTNamed{
				Name: "MOTION_BLOCKING",
				Data: randomHeightMap,
				Type: 12,
			},
		},
	})

	nbtRead, _ := dataTypes.NBTRead(heightMaps, 0)

	fmt.Println(dataTypes.NBTToString(nbtRead, 0))

	fmt.Println(hex.Dump(heightMaps))

	/**
	inData, exception := ioutil.ReadFile("world/region/r.0.0.mca")

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return
	}

	region := dataTypes.ReadRegionFile(inData)

	chunk := region.Chunks[0]

	nbtMap := dataTypes.NBTAsMap(chunk.Data)

	asJson, exception := json.Marshal(nbtMap)

	fmt.Println("As json:")
	fmt.Println(string(asJson))

	*/
	//output := nbtMap.(map[string]interface{})["Unnamed"].(map[string]interface{})["Compound_0"]
	//level := output.(map[string]interface{})["Level"].(map[string]interface{})["Compound_0"].(map[string]interface{})

	//palette := level["Sections"].(map[string]interface{})["List-0"].(map[string]interface{})["Compound_1"].(map[string]interface{})

	//biomes := level["Biomes"].([]uint8)

	for x := 0; x < 7; x++ {
		for y := 0; y < 7; y++ {
			chunkSections := make([]dataTypes.NetChunkSection, 32)
			for i := 0; i < 32; i++ {
				randomBlocks := make([]int64, 4096)
				for z := 0; z < 4096; z++ {
					randomBlocks[z] = rand.Int63n(int64(math.MaxInt64 - (z * 1000)))
				}

				chunkSections[i] = dataTypes.NetChunkSection{
					BlockCount:   1,
					BitsPerBlock: 4,
					Palette:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
					DataArray:    randomBlocks,
				}
			}

			chunkRaw := dataTypes.WriteChunk(chunkSections)

			chunkData := Packet(&ChunkData{
				X:                x,
				Z:                y,
				FullChunk:        true,
				PrimaryBitMask:   0b1111111111111111111111111111111,
				HeightMap:        heightMaps,
				Biomes:           randomBiomes,
				DataSize:         len(chunkRaw),
				Data:             chunkRaw,
				BlockEntityCount: 0,
				BlockEntities:    make([]byte, 0),
			})

			connection.SendPacket(&chunkData)
		}
	}

	playerSpawn := Packet(&SpawnPosition{
		Location: 0,
	})

	connection.SendPacket(&playerSpawn)

	playerPos := Packet(&PlayerPositionAndLook{
		X:          connection.Player.X,
		Y:          connection.Player.Y,
		Z:          connection.Player.Z,
		Yaw:        connection.Player.Yaw,
		Pitch:      connection.Player.Pitch,
		Flags:      0,
		TeleportID: 12345,
	})

	connection.SendPacket(&playerPos)

	connection.Minecraft.PlayerJoin(connection)

	//TODO Player info

	//viewPos := Packet(&UpdateViewPosition{
	//	ChunkX: 1,
	//	ChunkZ: 2,
	//})
	//
	//connection.SendPacket(&viewPos)
}
