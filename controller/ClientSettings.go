package controller

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
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

	if connection.Joined {
		fmt.Println("Ignoring ClientSettings from already joined client")
		return
	}

	heldItemChange := packetType.Packet(&packetType.HeldItemChange{
		Slot:     0,
		IsServer: true,
	})
	connection.SendPacket(&heldItemChange)

	entityStatus := packetType.Packet(&packetType.EntityStatus{
		EntityID:     connection.Player.EntityID,
		EntityStatus: constants.PlayerOp4,
	})

	connection.SendPacket(&entityStatus)

	viewPos := packetType.Packet(&packetType.UpdateViewPosition{
		ChunkX: 6,
		ChunkZ: 6,
	})

	connection.SendPacket(&viewPos)

	// Region 0,0
	region := connection.Minecraft.DataStore.Map[0][0]

	for _, chunk := range region.Chunks {
		if chunk == nil {
			continue
		}
		if len(chunk.Sections) == 0 {
			fmt.Println("not sending blank chunk")
			continue
		}
		//fmt.Println("Sending chunk", chunk.XPos, chunk.YPos, chunk.ZPos)
		chunkRaw := dataTypes.WriteNetChunk(chunk, connection.Minecraft.DataStore.BlockData)

		lightMaskLength := (len(chunk.Sections) + 2) / 8
		lightMask := dataTypes.WriteVarInt(lightMaskLength)
		for i := 0; i < lightMaskLength; i++ {
			lightMask = append(lightMask, dataTypes.WriteLong(int64(0))...)
		}

		chunkData := packetType.Packet(&packetType.ChunkData{
			X: int(chunk.XPos),
			Z: int(chunk.ZPos),
			HeightMap: packetType.HeightMapOuter{Inner: packetType.HeightMap{
				MotionBlocking: chunk.Heightmaps.MotionBlocking,
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

	playerSpawn := packetType.Packet(&packetType.SpawnPosition{
		Location: 0,
		Angle:    0,
	})

	connection.SendPacket(&playerSpawn)

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
	connection.Joined = true

	connection.Minecraft.PlayerJoin(connection)

	//TODO Player info

}
