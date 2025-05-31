package controller

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type ClientSettings struct {
	CurrentPacket *packetType.ClientSettings
}

func (cs *ClientSettings) GetPacketStruct() packetType.Packet {
	return &packetType.ClientSettings{}
}

func (cs *ClientSettings) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
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

	//// Region 0,0
	//region := connection.Minecraft.DataStore.Map[0][0]
	//for _, chunk := range region.Chunks {
	//    if chunk == nil {
	//        fmt.Println("chunks not here")
	//        continue
	//    }
	//    if len(chunk.Sections) == 0 {
	//        fmt.Println("not sending blank chunk")
	//        continue
	//    }
	//    //fmt.Println("Sending chunk", chunk.XPos, chunk.YPos, chunk.ZPos)
	//    chunkRaw := dataTypes.WriteNetChunk(chunk, connection.Minecraft.DataStore.BlockData)
	//
	//    lightMaskLength := (len(chunk.Sections) + 2) / 8
	//    lightMask := dataTypes.WriteVarInt(lightMaskLength)
	//    for i := 0; i < lightMaskLength; i++ {
	//        lightMask = append(lightMask, dataTypes.WriteLong(int64(0))...)
	//    }
	//
	//    chunkData := packetType.Packet(&packetType.ChunkData{
	//        X:                    int(chunk.XPos),
	//        Z:                    int(chunk.ZPos),
	//        HeightMapsLength:     0,
	//        DataSize:             len(chunkRaw),
	//        Data:                 nil,
	//        SkyLightMask:         lightMask,
	//        BlockLightMask:       lightMask,
	//        EmptySkyLightMask:    lightMask,
	//        EmptyBlockLightMask:  lightMask,
	//        SkyLightArrayCount:   0,
	//        SkyLightArrays:       []byte{},
	//        BlockLightArrayCount: 0,
	//        BlockLightArrays:     []byte{},
	//    })
	//
	//    connection.SendPacket(&chunkData)
	//}

	playerSpawn := packetType.Packet(&packetType.SpawnPosition{
		DimensionType:    0,
		DimensionName:    "minecraft:overworld",
		HashedSeed:       0,
		GameMode:         1,
		PreviousGameMode: 1,
		IsDebug:          false,
		IsFlat:           false,
		HasDeathLocation: false,
		PortalCooldown:   0,
		SeaLevel:         64,
		DataKept:         0,
	})

	connection.SendPacket(&playerSpawn)

	playerPos := packetType.Packet(&packetType.PlayerPositionAndLook{
		TeleportID: 12345,
		X:          connection.Player.X,
		Y:          connection.Player.Y,
		Z:          connection.Player.Z,
		VelX:       0,
		VelY:       0,
		VelZ:       0,
		Yaw:        connection.Player.Yaw,
		Pitch:      connection.Player.Pitch,
		Flags1:     0,
		Flags2:     0,
		Flags3:     0,
		Flags4:     0,
	})

	connection.SendPacket(&playerPos)
	connection.Joined = true

	connection.Minecraft.PlayerJoin(connection)

	//TODO Player info

}
