package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/dataTypes/nbt"
	"github.com/Ocelotworks/MinecraftGo/helpers"
	"github.com/Ocelotworks/MinecraftGo/world"
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
	WorldAge             int64
	TimeOfDay            int64
	DataStore            *DataStore
	Registries           map[string]map[string]nbt.NBTValue
	ChunkManager         *world.ChunkManager
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
		EnableEncryption:     false,
		CompressionThreshold: -1,
		GlobalEntityCounter:  1,
		DataStore:            NewDataStore(),
		Registries:           make(map[string]map[string]nbt.NBTValue),
	}

	mc.LoadRegistries()
	mc.LoadChunks()

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

	//currentPlayersPacket := packetType.Packet(&packetType.PlayerInfoAddPlayer{
	//	Action:  0,
	//	Players: []entity.Player{*connection.Player},
	//})
	//
	//mc.SendToAllExcept(connection, &currentPlayersPacket)

	currentPlayers := make([]entity.Player, 0)
	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		currentPlayers = append(currentPlayers, *con.Player)
	}

	//currentPlayersPacket = packetType.Packet(&packetType.PlayerInfoAddPlayer{
	//	Action:  0,
	//	Players: currentPlayers,
	//})
	//
	//connection.SendPacket(&currentPlayersPacket)

	// TODO: why is this duplicated like this?
	for _, con := range mc.Connections {
		if con.Player == nil || con == connection {
			continue
		}
		packet := packetType.Packet(&packetType.SpawnEntity{
			EntityID: con.Player.EntityID,
			UUID:     con.Player.UUID,
			Type:     constants.EntityTypePlayer,
			X:        con.Player.X,
			Y:        con.Player.Y,
			Z:        con.Player.Z,
			Yaw:      byte(con.Player.Yaw),
			Pitch:    byte(con.Player.Pitch),
		})
		connection.SendPacket(&packet)
	}

	packet := packetType.Packet(&packetType.SpawnEntity{
		EntityID: connection.Player.EntityID,
		UUID:     connection.Player.UUID,
		Type:     constants.EntityTypePlayer,
		X:        connection.Player.X,
		Y:        connection.Player.Y,
		Z:        connection.Player.Z,
		Yaw:      byte(connection.Player.Yaw),
		Pitch:    byte(connection.Player.Pitch),
	})
	mc.SendToAllExcept(connection, &packet)

	//yellow := entity.Yellow
	//chatMessageComponents := []entity.ChatMessageComponent{
	//    {
	//        Text:   connection.Player.Username,
	//        Colour: &yellow,
	//    },
	//}

	//chatMessage := entity.ChatMessage{
	//    Translate: "multiplayer.player.joined",
	//    With:      &chatMessageComponents,
	//}
	//
	//go mc.SendMessage(1, chatMessage)

	//bold := true
	//headerComponents := []entity.ChatMessageComponent{
	//    {
	//        Text: "Big P MC",
	//        Bold: &bold,
	//    },
	//}

	//listHeader := entity.ChatMessage{
	//    Translate: "%s",
	//    With:      &headerComponents,
	//}
	//
	//listFooter := entity.ChatMessage{
	//    Translate: "",
	//}

	//header, _ := json.Marshal(listHeader)
	//footer, _ := json.Marshal(listFooter)

	//packet = &packetType.PlayerListHeaderAndFooter{
	//    Header: string(header),
	//    Footer: string(footer),
	//}
	//
	//connection.SendPacket(&packet)

	go connection.sendKeepAlive()
}

func (mc *Minecraft) PlayerLeave(connection *Connection) {
	// Send remove entity if player.entityID != 0
	if connection.Player != nil && connection.Player.EntityID != 0 {
		mc.ConnectedPlayers--
		// Send player list update
		//currentPlayersPacket := packetType.Packet(&packetType.PlayerInfoRemovePlayer{
		//    Action:          4,
		//    NumberOfPlayers: 1,
		//    UUID:            connection.Player.UUID,
		//})
		//
		//mc.SendToAllExcept(connection, &currentPlayersPacket)

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
		Properties: []packetType.LoginProperty{
			{
				Key:   "test",
				Value: "test2",
			},
		},
	})

	connection.SendPacket(&returnPacket)

}

func (mc *Minecraft) SendMessage(messageType byte, message entity.ChatMessage) {
	chatMessageJson, exception := json.Marshal(message)

	if exception != nil {
		fmt.Println(exception)
		return
	}

	blankUuid := make([]byte, 16)
	for i := range blankUuid {
		blankUuid[i] = 0
	}

	chatPacket := packetType.Packet(&packetType.ChatMessage{
		ChatData: string(chatMessageJson),
		Position: messageType,
		Sender:   blankUuid,
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

func (mc *Minecraft) LoadRegistries() {
	fmt.Println("Loading registries...")
	registryData := helpers.LoadAllRegistries()
	for registryName, registryEntries := range registryData {

		_, ok := mc.Registries[registryName]
		if !ok {
			mc.Registries[registryName] = make(map[string]nbt.NBTValue)
		}

		for entryFile, entryData := range registryEntries {
			nbtValue, _ := nbt.JSONToNBT([]byte(entryData))
			mc.Registries[registryName][entryFile] = &nbtValue
		}
	}
}

func (mc *Minecraft) LoadChunks() {

	mc.ChunkManager = world.NewChunkManager()
	mc.ChunkManager.LoadRegion(0, 0)

}
