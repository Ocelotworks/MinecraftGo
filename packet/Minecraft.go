package packet

import (
	"encoding/json"
	"fmt"
	"github.com/Ocelotworks/MinecraftGo/dataTypes"
	"github.com/Ocelotworks/MinecraftGo/entity"
	"github.com/gofrs/uuid"
)

type Minecraft struct {
	Connections          []*Connection
	ServerName           entity.ChatMessageComponent
	ConnectedPlayers     int
	MaxPlayers           int
	EnableEncryption     bool
	CompressionThreshold int
}

func CreateMinecraft() *Minecraft {
	purple := entity.Purple
	return &Minecraft{
		//Connections: make([]*packet.Connection, 0),
		ServerName: entity.ChatMessageComponent{
			Text:   "Petecraft",
			Colour: &purple,
		},
		MaxPlayers:           255,
		ConnectedPlayers:     0,
		EnableEncryption:     true,
		CompressionThreshold: -1,
	}
}

func calculateDeltas(player *entity.Player, newX float64, newY float64, newZ float64) (int16, int16, int16, float64, float64, float64) {
	return int16((newX*32 - player.X*32) * 128), int16((newY*32 - player.Y*32) * 128), int16((newZ*32 - player.Z*32) * 128), newX - player.X, newY - player.Y, newZ - player.Z
}

func (mc *Minecraft) UpdatePlayerPosition(connection *Connection, newX float64, newY float64, newZ float64, newYaw float32, newPitch float32) {
	player := connection.Player
	if newX == 0 && newY == 0 && newZ == 0 && newYaw != 00 && newPitch != 0 {
		player.Yaw = newYaw
		player.Pitch = newPitch
		packet := Packet(&EntityRotation{
			EntityID: player.EntityID,
			Yaw:      byte(player.Yaw),
			Pitch:    byte(player.Pitch),
			OnGround: true,
		})
		mc.SendToAllExcept(connection, &packet)

		headLookPacket := Packet(&EntityHeadLook{
			EntityID: player.EntityID,
			Yaw:      byte(player.Yaw),
		})

		mc.SendToAllExcept(connection, &headLookPacket)

		return
	}

	deltaX, deltaY, deltaZ, blockDeltaX, blockDeltaY, blockDeltaZ := calculateDeltas(connection.Player, newX, newY, newZ)
	if deltaX != 0 || deltaY != 0 || deltaZ != 0 {
		fmt.Println("Deltas ", deltaX, deltaY, deltaZ)
		player.X = newX
		player.Y = newY
		player.Z = newZ
		if newYaw != 0 {
			player.Yaw = newYaw
			player.Pitch = newPitch
		}
		var packet Packet
		if blockDeltaX > 8 || blockDeltaY > 8 || blockDeltaZ > 8 || blockDeltaX < -8 || blockDeltaY < -8 || blockDeltaZ < -8 {
			packet = Packet(&EntityTeleport{
				EntityID: player.EntityID,
				X:        newX,
				Y:        newY,
				Z:        newZ,
				Yaw:      byte(player.Yaw),
				Pitch:    byte(player.Pitch),
				OnGround: true,
			})
		} else {
			if player.Yaw != 0 {
				packet = Packet(&EntityPositionAndRotation{
					EntityID: player.EntityID,
					DeltaX:   deltaX,
					DeltaY:   deltaY,
					DeltaZ:   deltaZ,
					Yaw:      byte(player.Yaw),
					Pitch:    byte(player.Pitch),
					OnGround: true,
				})

				headLookPacket := Packet(&EntityHeadLook{
					EntityID: player.EntityID,
					Yaw:      byte(player.Yaw),
				})

				mc.SendToAllExcept(connection, &headLookPacket)

			} else {
				packet = Packet(&EntityPosition{
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

func (mc *Minecraft) PlayerJoin(connection *Connection) {
	mc.ConnectedPlayers++

	currentPlayersPacket := Packet(&PlayerInfoAddPlayer{
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

	currentPlayersPacket = Packet(&PlayerInfoAddPlayer{
		Action:  0,
		Players: currentPlayers,
	})

	connection.SendPacket(&currentPlayersPacket)

	for _, con := range mc.Connections {
		if con.Player == nil || con == connection {
			continue
		}
		packet := Packet(&SpawnPlayer{
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

	packet := Packet(&SpawnPlayer{
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

	go connection.sendKeepAlive()
}

func (mc *Minecraft) SendToAllExcept(connection *Connection, packet *Packet) {
	for _, con := range mc.Connections {
		if con == connection || con.Player == nil {
			continue
		}
		con.SendPacket(packet)
	}
}

func (mc *Minecraft) SendToAll(packet *Packet) {
	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		con.SendPacket(packet)
	}
}

func (mc *Minecraft) StartPlayerJoin(connection *Connection) {
	if connection.Minecraft.CompressionThreshold > 0 {
		compressionPacket := Packet(&SetCompression{
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

	returnPacket := Packet(&LoginSuccess{
		UUID:     stringUUID.String(),
		Username: connection.Player.Username,
	})

	connection.SendPacket(&returnPacket)

	connection.State = PLAY

	joinGame := Packet(&JoinGame{
		EntityID:            connection.Player.EntityID,
		Gamemode:            0,
		Dimension:           0,
		HashedSeed:          71495747907944700,
		MaxPlayers:          byte(connection.Minecraft.MaxPlayers),
		LevelType:           "default",
		ViewDistance:        32,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
	})

	connection.SendPacket(&joinGame)

	pluginMessage := Packet(&PluginMessage{
		IsServer:   false,
		Identifier: "minecraft:brand",
		ByteArray:  dataTypes.WriteString("BigPMC"),
	})

	connection.SendPacket(&pluginMessage)

	difficulty := Packet(&ServerDifficulty{
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

	chatPacket := Packet(&ChatMessage{
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
