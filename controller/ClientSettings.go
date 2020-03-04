package controller

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"io/ioutil"
	"math"
)

type ClientSettings struct {
	CurrentPacket *packetType.ClientSettings
}

func (cs *ClientSettings) GetPacketStruct() packetType.Packet {
	return &packetType.ClientSettings{}
}

func (cs *ClientSettings) Init(currentPacket packetType.Packet) {
	cs.CurrentPacket = currentPacket.(*packetType.ClientSettings)
}

func (cs *ClientSettings) Handle(packet []byte, connection *Connection) {
	connection.Player.Settings = entity.PlayerSettings{
		Locale:             cs.CurrentPacket.Locale,
		ViewDistance:       cs.CurrentPacket.ViewDistance,
		ChatMode:           cs.CurrentPacket.ChatMode,
		ChatColours:        cs.CurrentPacket.ChatColours,
		DisplayedSkinParts: cs.CurrentPacket.DisplayedSkinParts,
		MainHand:           cs.CurrentPacket.MainHand,
	}

	randomBiomes := make([]int, 1024)
	for i := 0; i < 1024; i++ {
		randomBiomes[i] = 1
	}

	randomHeightMap := make([]int64, 36)
	for i := 0; i < 36; i++ {
		randomHeightMap[i] = math.MaxInt64
	}

	blockFile, exception := ioutil.ReadFile("data/blocks.json")

	blockData := make(map[string]entity.BlockData)
	exception = json.Unmarshal(blockFile, &blockData)

	if exception != nil {
		fmt.Println("Reading block data", exception)
		return
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

	inData, exception := ioutil.ReadFile("C:\\Users\\Peter\\AppData\\Roaming\\.minecraft\\saves\\MCGO Flat Test 2\\region\\r.0.0.mca")

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return
	}

	region := dataTypes.ReadRegionFile(inData, connection.Minecraft.BlockData)

	//fmt.Println(level)

	//palette := level["Sections"].(map[string]interface{})["List-0"].(map[string]interface{})["Compound_1"].(map[string]interface{})["Palette"].(map[string]interface{})

	//byte((len(castBlockStates)*64)/4096),

	for i, chunk := range region.Chunks {
		chunkRaw := dataTypes.WriteChunk(chunk.Sections)
		fmt.Println(len(chunk.Sections))
		if len(chunk.Sections) == 0 {
			fmt.Println("not sending blank chunk")
			continue
		}
		chunkData := packetType.Packet(&packetType.ChunkData{
			X:                i % 32,
			Z:                i / 32,
			FullChunk:        true,
			PrimaryBitMask:   int(math.Pow(2, float64(len(chunk.Sections))) - 1),
			HeightMap:        heightMaps,
			Biomes:           randomBiomes,
			DataSize:         len(chunkRaw),
			Data:             chunkRaw,
			BlockEntityCount: 0,
			BlockEntities:    make([]byte, 0),
		})

		connection.SendPacket(&chunkData)
		//break
	}

	playerSpawn := packetType.Packet(&packetType.SpawnPosition{
		Location: 0,
	})

	connection.SendPacket(&playerSpawn)

	playerPos := packetType.Packet(&packetType.PlayerPositionAndLook{
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
	////	ChunkZ: 2,
	////})
	//
	//connection.SendPacket(&viewPos)

}
