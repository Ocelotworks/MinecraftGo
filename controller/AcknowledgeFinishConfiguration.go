package controller

import (
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

	//pluginMessage := packetType.Packet(&packetType.PluginMessage{
	//	IsServer:   false,
	//	Identifier: "minecraft:brand",
	//	ByteArray:  dataTypes.WriteString("BigPMC"),
	//})
	//
	//connection.SendPacket(&pluginMessage)
	//
	//difficulty := packetType.Packet(&packetType.ServerDifficulty{
	//	Difficulty:       0,
	//	DifficultyLocked: false,
	//})
	//
	//connection.SendPacket(&difficulty)

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

	// Place the player before the chunks go out, as PlayerList.placeNewPlayer
	// does: what the client waits on to leave the loading screen is the chunk
	// the player is standing in, so it has to know where that is.
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
		Flags:      0,
	})

	connection.SendPacket(&playerPos)

	playerSpawn := packetType.Packet(&packetType.SpawnPosition{
		DimensionName: "minecraft:overworld",
		Location:      &dataTypes.Position{X: 0, Y: 64, Z: 0},
		Yaw:           0,
		Pitch:         0,
	})

	connection.SendPacket(&playerSpawn)

	// The client keeps its chunk cache centred here and discards chunks that
	// fall outside its view distance of it, so this has to be the chunk the
	// player is actually in or the chunks below get thrown away on arrival.
	viewPos := packetType.Packet(&packetType.UpdateViewPosition{
		ChunkX: connection.Minecraft.BlockCoordToChunkCoord(connection.Player.X),
		ChunkZ: connection.Minecraft.BlockCoordToChunkCoord(connection.Player.Z),
	})

	connection.SendPacket(&viewPos)

	// Until this arrives the client sits on "Loading terrain..." forever: its
	// LevelLoadTracker starts in WaitingForServer, and nothing else moves it
	// out of that state, which has no timeout of its own.  This is what
	// PlayerList.sendLevelInfo ends with, and it moves the client on to
	// waiting for the chunk it is standing in to render — or for thirty
	// seconds, whichever comes first — and then into the world.
	chunksLoadStart := packetType.Packet(&packetType.GameEvent{
		Event: packetType.EventTypeWaitingForChunks,
		Value: 0,
	})

	connection.SendPacket(&chunksLoadStart)

	// The view distance this uses came from the client information sent during
	// configuration.  The play state has a copy of that packet, but the client
	// only sends it again when its options change, so a join that waited for
	// one would usually wait forever.
	connection.Minecraft.SendChunks(connection)

	connection.Joined = true

	//entityEffect := packetType.Packet(&packetType.EntityEffect{
	//	EntityId:  connection.Player.EntityID,
	//	EffectId:  15, // TODO https://github.com/PrismarineJS/minecraft-data/blob/abd43b3b5f6627cfc6abc66f7eee9598cc00e44f/data/pc/1.21.5/effects.json
	//	Amplifier: 0,
	//	Duration:  -1,
	//	Flags:     0,
	//})
	//
	//connection.SendPacket(&entityEffect)

	connection.Minecraft.PlayerJoin(connection)
}
