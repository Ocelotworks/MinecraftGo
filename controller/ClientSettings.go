package controller

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"io/ioutil"
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

	heldItemChange := packetType.Packet(&packetType.HeldItemChange{
		Slot:     0,
		IsServer: true,
	})
	connection.SendPacket(&heldItemChange)

	// If I ignore this it will go away
	//tagBytes := dataTypes.WriteString("minecraft:block")
	//blockTagBytes := make([]byte, 0)
	//blockTagCount := 0
	//for _, block := range connection.Minecraft.BlockData {
	//    for _, state := range block.States {
	//        blockTagBytes = append(blockTagBytes, dataTypes.WriteVarInt(state.ID)...)
	//        blockTagCount++
	//    }
	//}
	//
	//tagBytes = append(tagBytes, dataTypes.WriteVarInt(blockTagCount)...)
	//tagBytes = append(tagBytes, blockTagBytes...)
	//tagBytes = append(tagBytes, dataTypes.WriteString("minecraft:item")...)
	//tagBytes = append(tagBytes, dataTypes.WriteVarInt(0)...)
	//tagBytes = append(tagBytes, dataTypes.WriteString("minecraft:fluid")...)
	//tagBytes = append(tagBytes, dataTypes.WriteVarInt(0)...)
	//tagBytes = append(tagBytes, dataTypes.WriteString("minecraft:entity_type")...)
	//tagBytes = append(tagBytes, dataTypes.WriteVarInt(0)...)
	//tagBytes = append(tagBytes, dataTypes.WriteString("minecraft:game_event")...)
	//tagBytes = append(tagBytes, dataTypes.WriteVarInt(0)...)
	//
	//tags := packetType.Packet(&packetType.Tags{
	//    TagCount: 5,
	//    Tags:     tagBytes,
	//})
	//
	//connection.SendPacket(&tags)

	entityStatus := packetType.Packet(&packetType.EntityStatus{
		EntityID:     connection.Player.EntityID,
		EntityStatus: constants.PlayerOp4,
	})

	connection.SendPacket(&entityStatus)

	entityStatus2 := packetType.Packet(&packetType.EntityStatus{
		EntityID:     connection.Player.EntityID,
		EntityStatus: constants.PlayerReducedDebugInfoDisabled,
	})

	connection.SendPacket(&entityStatus2)

	viewPos := packetType.Packet(&packetType.UpdateViewPosition{
		ChunkX: 6,
		ChunkZ: 6,
	})

	connection.SendPacket(&viewPos)

	playerSpawn := packetType.Packet(&packetType.SpawnPosition{
		Location: 0,
		Angle:    0,
	})

	connection.SendPacket(&playerSpawn)

	randomHeightMap := make([]int64, 36)
	for i := range randomHeightMap {
		randomHeightMap[i] = 0x0100804020100804
	}

	inData, exception := ioutil.ReadFile("data/worlds/mcgo/region/r.0.0.mca") //ioutil.ReadFile("data/worlds/MCGO_FlatTest/region/r.0.0.mca")

	if exception != nil {
		fmt.Println("Reading file")
		fmt.Println(exception)
		return
	}

	region := dataTypes.ReadRegionFile(inData)

	for _, chunk := range region.Chunks {
		if chunk == nil {
			continue
		}
		if len(chunk.Sections) == 0 {
			fmt.Println("not sending blank chunk")
			continue
		}
		fmt.Println("Sending chunk", chunk.XPos, chunk.YPos, chunk.ZPos)
		chunkRaw := dataTypes.WriteNetChunk(chunk, connection.Minecraft.BlockData)

		lightMaskLength := (len(chunk.Sections) + 2) / 8
		lightMask := dataTypes.WriteVarInt(lightMaskLength)
		for i := 0; i < lightMaskLength; i++ {
			lightMask = append(lightMask, dataTypes.WriteLong(int64(0))...)
		}

		chunkData := packetType.Packet(&packetType.ChunkData{
			X: int(chunk.XPos),
			Z: int(chunk.ZPos),
			HeightMap: packetType.HeightMapOuter{Inner: packetType.HeightMap{
				MotionBlocking: randomHeightMap,
			}},
			DataSize:             len(chunkRaw),
			Data:                 chunkRaw,
			BlockEntityCount:     0,
			BlockEntities:        []byte{},
			TrustEdges:           true,
			SkyLightMask:         lightMask,
			BlockLightMask:       lightMask,
			EmptySkyLightMask:    lightMask,
			EmptyBlockLightMask:  lightMask,
			SkyLightArrayCount:   2048,
			SkyLightArrays:       make([]byte, 2048),
			BlockLightArrayCount: 2048,
			BlockLightArrays:     make([]byte, 2048),
		})

		connection.SendPacket(&chunkData)
	}

	playerPos := packetType.Packet(&packetType.PlayerPositionAndLook{
		X:               connection.Player.X,
		Y:               connection.Player.Y,
		Z:               connection.Player.Z,
		Yaw:             connection.Player.Yaw,
		Pitch:           connection.Player.Pitch,
		Flags:           0,
		TeleportID:      12345,
		DismountVehicle: true,
	})

	connection.SendPacket(&playerPos)

	connection.Minecraft.PlayerJoin(connection)

	//TODO Player info

}
