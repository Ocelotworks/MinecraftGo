package controller

import (
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
)

type AcknowledgeFinishConfiguration struct {
	CurrentPacket *packetType.AcknowledgeFinishConfiguration
	Minecraft     *Minecraft
}

func (lpr *AcknowledgeFinishConfiguration) GetPacketStruct() packetType.Packet {
	return &packetType.AcknowledgeFinishConfiguration{}
}

func (lpr *AcknowledgeFinishConfiguration) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	lpr.CurrentPacket = currentPacket.(*packetType.AcknowledgeFinishConfiguration)
	lpr.Minecraft = minecraft
}

func (lpr *AcknowledgeFinishConfiguration) Handle(packet []byte, connection *Connection) {
	connection.State = PLAY

	//os.WriteFile("ours2.nbt", compound.Write(), 0644)
	//os.WriteFile("theirs2.nbt",, 0644)

	joinGame := packetType.Packet(&packetType.JoinGame{
		EntityID:            connection.Player.EntityID,
		IsHardcore:          false,
		DimensionNames:      []string{"minecraft:overworld"}, // TODO: this but properly
		MaxPlayers:          connection.Minecraft.MaxPlayers,
		ViewDistance:        32,
		SimulationDistance:  32,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		DoLimitedCrafting:   false,
		DimensionType:       0,
		DimensionName:       "minecraft:overworld",
		HashedSeed:          71495747907944700,
		Gamemode:            1,
		PreviousGamemode:    1,
		IsDebug:             false,
		IsFlat:              false,
		HasDeathLocation:    false,
		PortalCooldown:      0,
		SeaLevel:            64,
		EnforcesSecureChat:  false,
	})

	connection.SendPacket(&joinGame)

	pluginMessage := packetType.Packet(&packetType.PluginMessage{
		IsServer:   false,
		Identifier: "minecraft:brand",
		ByteArray:  dataTypes.WriteString("BigPMC"),
	})

	connection.SendPacket(&pluginMessage)

	difficulty := packetType.Packet(&packetType.ServerDifficulty{
		Difficulty:       0,
		DifficultyLocked: false,
	})

	connection.SendPacket(&difficulty)

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

	waitingForChunksEvent := packetType.Packet(&packetType.GameEvent{
		Event: packetType.EventTypeWaitingForChunks,
		Value: 0,
	})
	connection.SendPacket(&waitingForChunksEvent)

	for chunkX := range lpr.Minecraft.ChunkManager.Chunks {
		for chunkZ := range lpr.Minecraft.ChunkManager.Chunks[chunkX] {
			chunk := lpr.Minecraft.ChunkManager.Chunks[chunkX][chunkZ]
			if chunk.BlocksInChunk == 0 {
				fmt.Printf("Chunk %d, %d is empty\n", chunkX, chunkZ)
				continue
			}
			lightMaskLength := (16 + 2) / 8
			lightMask := dataTypes.WriteVarInt(lightMaskLength)
			emptyLightMask := dataTypes.WriteVarInt(lightMaskLength)
			for i := 0; i < lightMaskLength; i++ {
				lightMask = append(lightMask, dataTypes.WriteLong(int64(0))...)
				emptyLightMask = append(emptyLightMask, dataTypes.WriteLong(int64(0))...)
			}

			chunkData := make([]byte, 0)

			for i := 0; i < 24; i++ {

				chunkData = append(chunkData, dataTypes.WriteShort(int16(chunk.BlocksInChunk))...) // Block count
				//chunkData = append(chunkData, dataTypes.WriteShort(int16(4096))...) // Block count
				chunkData = append(chunkData, 15) // Bits per block
				//chunkData = append(chunkData, dataTypes.WriteVarInt(1)...)

				for y := range chunk.Blocks {
					for z := range chunk.Blocks[y] {
						for _, block := range chunk.Blocks[y][z] {
							blockData, ok := lpr.Minecraft.DataStore.BlockData[block.Type]
							if !ok {
								fmt.Println("unable to find block of name", block.Type)
								chunkData = append(chunkData, dataTypes.WriteLong(int64(0))...)
							} else {
								chunkData = append(chunkData, dataTypes.WriteLong(int64(blockData.ID))...)
							}

						}
					}
				}

				// Biome data
				chunkData = append(chunkData, 0)
				chunkData = append(chunkData, dataTypes.WriteVarInt(1)...)
			}

			chunkPacket := packetType.Packet(&packetType.ChunkData{
				X:                    int(chunkX),
				Z:                    int(chunkZ),
				HeightMapsLength:     0,
				HeightMapData:        []byte{},
				DataSize:             len(chunkData),
				Data:                 chunkData,
				BlockEntityCount:     0,
				BlockEntities:        []byte{},
				SkyLightMask:         lightMask,
				BlockLightMask:       lightMask,
				EmptySkyLightMask:    emptyLightMask,
				EmptyBlockLightMask:  emptyLightMask,
				SkyLightArrayCount:   0,
				SkyLightArrays:       []byte{},
				BlockLightArrayCount: 0,
				BlockLightArrays:     []byte{},
			})

			connection.SendPacket(&chunkPacket)

		}
	}

	//playerSpawn := packetType.Packet(&packetType.SpawnPosition{
	//    DimensionType:    0,
	//    DimensionName:    "minecraft:overworld",
	//    HashedSeed:       0,
	//    GameMode:         1,
	//    PreviousGameMode: 1,
	//    IsDebug:          false,
	//    IsFlat:           false,
	//    HasDeathLocation: false,
	//    PortalCooldown:   0,
	//    SeaLevel:         64,
	//    DataKept:         0,
	//})
	//
	//connection.SendPacket(&playerSpawn)

	playerPos := packetType.Packet(&packetType.PlayerPositionAndLook{
		TeleportID: 12345,
		X:          connection.Player.X,
		Y:          connection.Player.Y,
		Z:          connection.Player.Z,
		VelX:       1,
		VelY:       2,
		VelZ:       3,
		Yaw:        connection.Player.Yaw,
		Pitch:      connection.Player.Pitch,
		Flags1:     1,
		Flags2:     2,
		Flags3:     3,
		Flags4:     4,
	})

	connection.SendPacket(&playerPos)
	connection.Joined = true

	entityEffect := packetType.Packet(&packetType.EntityEffect{
		EntityId:  connection.Player.EntityID,
		EffectId:  15, // TODO https://github.com/PrismarineJS/minecraft-data/blob/abd43b3b5f6627cfc6abc66f7eee9598cc00e44f/data/pc/1.21.5/effects.json
		Amplifier: 0,
		Duration:  -1,
		Flags:     0,
	})

	connection.SendPacket(&entityEffect)

	connection.Minecraft.PlayerJoin(connection)
}
