package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"io/ioutil"
	"math"
	"time"

	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/gofrs/uuid"
)

type Minecraft struct {
	Connections          []*Connection
	ServerName           entity.ChatMessageComponent
	ConnectedPlayers     int
	MaxPlayers           int
	EnableEncryption     bool
	CompressionThreshold int
	GlobalEntityCounter  int
	BlockData            map[string]entity.BlockData
	WorldAge             int64
	TimeOfDay            int64
}

func CreateMinecraft() *Minecraft {
	purple := entity.Purple

	mc := &Minecraft{
		//Connections: make([]*packet.Connection, 0),
		ServerName: entity.ChatMessageComponent{
			Text:   "Petecraft",
			Colour: &purple,
		},
		MaxPlayers:           255,
		ConnectedPlayers:     0,
		EnableEncryption:     true,
		CompressionThreshold: -1,
		GlobalEntityCounter:  1,
		BlockData:            make(map[string]entity.BlockData),
	}

	blockFile, exception := ioutil.ReadFile("data/blocks.json")
	exception = json.Unmarshal(blockFile, &mc.BlockData)

	if exception != nil {
		fmt.Println("Error reading block data", exception)
		return mc
	}

	//go mc.timeTracker()

	return mc
}

func calculateDeltas(player *entity.Player, newX float64, newY float64, newZ float64) (int16, int16, int16, float64, float64, float64) {
	return int16((newX*32 - player.X*32) * 128), int16((newY*32 - player.Y*32) * 128), int16((newZ*32 - player.Z*32) * 128), newX - player.X, newY - player.Y, newZ - player.Z
}

func calculateRotation(angle float32) byte {
	rotation := byte(0)

	if angle != 0 {
		rotation = byte((angle / 360) * 254)
	}

	return rotation
}

func (mc *Minecraft) UpdatePlayerPosition(connection *Connection, newX float64, newY float64, newZ float64, newYaw float32, newPitch float32) {
	player := connection.Player

	yawRotation := calculateRotation(newYaw)
	pitchRotation := calculateRotation(newPitch)

	if newX == 0 && newY == 0 && newZ == 0 && newYaw != 00 && newPitch != 0 {
		// Convert degrees from n/260 to n/254

		packet := packetType.Packet(&packetType.EntityRotation{
			EntityID: player.EntityID,
			Yaw:      yawRotation,
			Pitch:    pitchRotation,
			OnGround: true,
		})

		mc.SendToAllExcept(connection, &packet)

		headLookPacket := packetType.Packet(&packetType.EntityHeadLook{
			EntityID: player.EntityID,
			Yaw:      yawRotation,
		})

		mc.SendToAllExcept(connection, &headLookPacket)

		return
	}

	mc.CalculateChunkBoundaryCrossing(connection, newX, newZ)

	deltaX, deltaY, deltaZ, blockDeltaX, blockDeltaY, blockDeltaZ := calculateDeltas(connection.Player, newX, newY, newZ)
	if deltaX != 0 || deltaY != 0 || deltaZ != 0 {
		player.X = newX
		player.Y = newY
		player.Z = newZ
		if newYaw != 0 {
			player.Yaw = newYaw
			player.Pitch = newPitch
		}
		var packet packetType.Packet
		if blockDeltaX > 8 || blockDeltaY > 8 || blockDeltaZ > 8 || blockDeltaX < -8 || blockDeltaY < -8 || blockDeltaZ < -8 {
			packet = packetType.Packet(&packetType.EntityTeleport{
				EntityID: player.EntityID,
				X:        newX,
				Y:        newY,
				Z:        newZ,
				Yaw:      yawRotation,
				Pitch:    pitchRotation,
				OnGround: true,
			})
		} else {
			if player.Yaw != 0 {
				packet = packetType.Packet(&packetType.EntityPositionAndRotation{
					EntityID: player.EntityID,
					DeltaX:   deltaX,
					DeltaY:   deltaY,
					DeltaZ:   deltaZ,
					Yaw:      yawRotation,
					Pitch:    pitchRotation,
					OnGround: true,
				})

				headLookPacket := packetType.Packet(&packetType.EntityHeadLook{
					EntityID: player.EntityID,
					Yaw:      yawRotation,
				})

				mc.SendToAllExcept(connection, &headLookPacket)

			} else {
				packet = packetType.Packet(&packetType.EntityPosition{
					EntityID: player.EntityID,
					DeltaX:   deltaX,
					DeltaY:   deltaY,
					DeltaZ:   deltaZ,
					OnGround: true,
				})
			}
		}
		mc.SendToAllExcept(connection, &packet)
	}
}

/**
* Calculate crossing chunk boundary on player movement
 */
func (mc *Minecraft) CalculateChunkBoundaryCrossing(connection *Connection, newX float64, newZ float64) {
	currentXChunk := mc.BlockCoordToChunkCoord(connection.Player.X)
	currentZChunk := mc.BlockCoordToChunkCoord(connection.Player.Z)

	newXChunk := mc.BlockCoordToChunkCoord(newX)
	newZChunk := mc.BlockCoordToChunkCoord(newZ)
	// If we have crossed a chunk boundary send a chunk boundary update
	if currentXChunk != newXChunk || currentZChunk != newZChunk {
		//fmt.Printf("Updating view position: %d-%d \r\n", newXChunk, newZChunk)
		updateViewPositionPacket := packetType.Packet(&packetType.UpdateViewPosition{
			ChunkX: newXChunk,
			ChunkZ: newZChunk,
		})

		connection.SendPacket(&updateViewPositionPacket)
	}
}

func (mc *Minecraft) BlockCoordToChunkCoord(coord float64) int {
	intcoord := int(math.Floor(coord))
	return intcoord >> 4
}

func (mc *Minecraft) PlayerJoin(connection *Connection) {
	mc.ConnectedPlayers++

	currentPlayersPacket := packetType.Packet(&packetType.PlayerInfoAddPlayer{
		Action:  0,
		Players: []entity.Player{*connection.Player},
	})

	mc.SendToAllExcept(connection, &currentPlayersPacket)

	currentPlayers := make([]entity.Player, 0)
	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		currentPlayers = append(currentPlayers, *con.Player)
	}

	currentPlayersPacket = packetType.Packet(&packetType.PlayerInfoAddPlayer{
		Action:  0,
		Players: currentPlayers,
	})

	connection.SendPacket(&currentPlayersPacket)

	for _, con := range mc.Connections {
		if con.Player == nil || con == connection {
			continue
		}
		packet := packetType.Packet(&packetType.SpawnPlayer{
			EntityID: con.Player.EntityID,
			UUID:     con.Player.UUID,
			X:        con.Player.X,
			Y:        con.Player.Y,
			Z:        con.Player.Z,
			Yaw:      byte(con.Player.Yaw),
			Pitch:    byte(con.Player.Pitch),
		})
		connection.SendPacket(&packet)
	}

	packet := packetType.Packet(&packetType.SpawnPlayer{
		EntityID: connection.Player.EntityID,
		UUID:     connection.Player.UUID,
		X:        connection.Player.X,
		Y:        connection.Player.Y,
		Z:        connection.Player.Z,
		Yaw:      byte(connection.Player.Yaw),
		Pitch:    byte(connection.Player.Pitch),
	})
	mc.SendToAllExcept(connection, &packet)

	yellow := entity.Yellow
	chatMessageComponents := []entity.ChatMessageComponent{
		{
			Text:   connection.Player.Username,
			Colour: &yellow,
		},
	}

	chatMessage := entity.ChatMessage{
		Translate: "multiplayer.player.joined",
		With:      &chatMessageComponents,
	}

	go mc.SendMessage(1, chatMessage)

	bold := true
	headerComponents := []entity.ChatMessageComponent{
		{
			Text: "Big P MC",
			Bold: &bold,
		},
	}

	listHeader := entity.ChatMessage{
		Translate: "%s",
		With:      &headerComponents,
	}

	listFooter := entity.ChatMessage{
		Translate: "",
	}

	header, _ := json.Marshal(listHeader)
	footer, _ := json.Marshal(listFooter)

	packet = &packetType.PlayerListHeaderAndFooter{
		Header: string(header),
		Footer: string(footer),
	}

	connection.SendPacket(&packet)

	go connection.sendKeepAlive()
}

func (mc *Minecraft) PlayerLeave(connection *Connection) {
	// Send remove entity if player.entityID != 0
	if connection.Player != nil && connection.Player.EntityID != 0 {
		mc.ConnectedPlayers--
		// Send player list update
		currentPlayersPacket := packetType.Packet(&packetType.PlayerInfoRemovePlayer{
			Action:          4,
			NumberOfPlayers: 1,
			UUID:            connection.Player.UUID,
		})

		mc.SendToAllExcept(connection, &currentPlayersPacket)

		// Send chat message
		yellow := entity.Yellow
		chatMessageComponents := []entity.ChatMessageComponent{
			{
				Text:   connection.Player.Username,
				Colour: &yellow,
			},
		}

		chatMessage := entity.ChatMessage{
			Translate: "multiplayer.player.left",
			With:      &chatMessageComponents,
		}

		go mc.SendMessage(1, chatMessage)

		fmt.Printf("Destroying player entity id: %d\n", connection.Player.EntityID)
		destroyEntityIDs := []int{
			connection.Player.EntityID,
		}

		destroyEntityPacket := packetType.Packet(&packetType.DestroyEntity{
			Count:     1,
			EntityIDs: destroyEntityIDs,
		})

		go mc.SendToAllExcept(connection, &destroyEntityPacket)
	}
}

func (mc *Minecraft) SendToAllExcept(connection *Connection, packet *packetType.Packet) {
	for _, con := range mc.Connections {
		if con == connection || con.Player == nil {
			continue
		}
		con.SendPacket(packet)
	}
}

func (mc *Minecraft) SendToAll(packet *packetType.Packet) {
	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		con.SendPacket(packet)
	}
}

func (mc *Minecraft) SendToAllInPlay(packet *packetType.Packet) {
	for _, con := range mc.Connections {
		if con.Player == nil || con.State != PLAY {
			continue
		}
		con.SendPacket(packet)
	}
}

func (mc *Minecraft) StartPlayerJoin(connection *Connection) {
	if connection.Minecraft.CompressionThreshold > 0 {
		compressionPacket := packetType.Packet(&packetType.SetCompression{
			Threshold: connection.Minecraft.CompressionThreshold,
		})
		connection.SendPacket(&compressionPacket)
		connection.EnableCompression = true
	}

	stringUUID, exception := uuid.FromBytes(connection.Player.UUID)

	if exception != nil {
		fmt.Println("Malformed UUID? ", exception)
		return
	}

	returnPacket := packetType.Packet(&packetType.LoginSuccess{
		UUID:     stringUUID.Bytes(),
		Username: connection.Player.Username,
	})

	connection.SendPacket(&returnPacket)

	connection.State = PLAY

	//inData, exception := ioutil.ReadFile("data/codec.nbt")
	//
	//compressed := bytes.NewReader(inData)
	//zr, exception := gzip.NewReader(compressed)
	//
	//if exception != nil {
	//    fmt.Println(exception)
	//}
	//
	//uncompressed, exception := ioutil.ReadAll(zr)
	//
	//if exception != nil {
	//    fmt.Println(exception)
	//}
	//
	//compound, _ := nbt.ReadNBT(uncompressed)
	//
	//codec := dataTypes.CodecOuterCompound{}
	//
	//structScanner.NBTStructScan(&codec, &compound)

	joinGame := packetType.Packet(&packetType.JoinGame{
		EntityID:         connection.Player.EntityID,
		IsHardcore:       false,
		Gamemode:         1,
		PreviousGamemode: 1,
		WorldNames:       []string{"minecraft:overworld"},
		DimensionCodec: dataTypes.DimensionCodec{
			DimensionType: dataTypes.DimensionTypeRegistry{
				Type: "minecraft:dimension_type",
				Value: []dataTypes.DimensionTypeRegistryEntry{{
					Name: "minecraft:overworld",
					Id:   0,
					Element: dataTypes.DimensionType{
						PiglinSafe:         0,
						Natural:            1,
						AmbientLight:       0,
						Infiniburn:         "#minecraft:infiniburn_overworld",
						RespawnAnchorWorks: 1,
						HasSkylight:        1,
						BedWorks:           1,
						Effects:            "minecraft:overworld",
						HasRaids:           1,
						MinY:               -64,
						Height:             384,
						LogicalHeight:      256,
						CoordinateScale:    1,
						Ultrawarm:          0,
						HasCeiling:         0,
					},
				}},
			},
			BiomeRegistry: dataTypes.BiomeRegistry{
				Type: "minecraft:worldgen/biome",
				Value: []dataTypes.BiomeRegistryEntry{{
					Name: "minecraft:plains",
					Id:   0,
					Element: dataTypes.Biome{
						Precipitation: "rain",
						Depth:         1,
						Temperature:   1,
						Scale:         1,
						Downfall:      0.4,
						Category:      "plains",
						Effects: dataTypes.BiomeEffect{
							SkyColor:      7907327,
							WaterFogColor: 329011,
							FogColor:      12638463,
							WaterColor:    4159204,
						},
					},
				}},
			},
		},
		Dimension: dataTypes.DimensionType{
			PiglinSafe:         0,
			Natural:            1,
			AmbientLight:       1,
			Infiniburn:         "",
			RespawnAnchorWorks: 1,
			HasSkylight:        1,
			BedWorks:           1,
			Effects:            "minecraft:overworld",
			HasRaids:           1,
			MinY:               0,
			Height:             256,
			LogicalHeight:      256,
			CoordinateScale:    1,
			Ultrawarm:          0,
			HasCeiling:         0,
		},
		DimensionName:       "minecraft:overworld",
		HashedSeed:          71495747907944700,
		MaxPlayers:          byte(connection.Minecraft.MaxPlayers),
		ViewDistance:        32,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		IsDebug:             false,
		IsFlat:              false,
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
}

func (mc *Minecraft) SendMessage(messageType byte, message entity.ChatMessage) {
	chatMessageJson, exception := json.Marshal(message)

	if exception != nil {
		fmt.Println(exception)
		return
	}

	chatPacket := packetType.Packet(&packetType.ChatMessage{
		ChatData: string(chatMessageJson),
		Position: messageType,
	})

	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		chatMode := con.Player.Settings.ChatMode
		if messageType == 2 || chatMode == 0 || (chatMode == 1 && messageType == 1) {
			con.SendPacket(&chatPacket)
		}
	}

}

func (mc *Minecraft) timeTracker() {
	for {
		<-time.After(time.Second)

		mc.TimeOfDay += 20
		mc.WorldAge += 20

		if mc.TimeOfDay > 24000 {
			mc.TimeOfDay = 0
		}

		packet := packetType.Packet(&packetType.TimeUpdate{
			WorldAge:  mc.WorldAge,
			TimeOfDay: mc.TimeOfDay,
		})

		mc.SendToAllInPlay(&packet)
	}
}
