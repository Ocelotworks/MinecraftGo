package packet

import (
	"../entity"
	"encoding/json"
	"fmt"
)

type Minecraft struct {
	Connections      []*Connection
	ServerName       entity.ChatMessageComponent
	ConnectedPlayers int
	MaxPlayers       int
	EnableEncryption bool
}

func CreateMinecraft() *Minecraft {
	purple := entity.Purple
	return &Minecraft{
		//Connections: make([]*packet.Connection, 0),
		ServerName: entity.ChatMessageComponent{
			Text:   "Petecraft",
			Colour: &purple,
		},
		MaxPlayers:       255,
		ConnectedPlayers: 0,
		EnableEncryption: false,
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

	playerDisplayName := entity.ChatMessageComponent{
		Text: connection.Player.Username,
	}

	playerDisplayNameJson, exception := json.Marshal(playerDisplayName)

	if exception != nil {
		fmt.Println("Marshalling player username")
		fmt.Println(exception)
	} else {
		currentPlayersPacket := Packet(&PlayerInfoAddPlayer{
			Action: 0,
			Players: []Player{
				{
					UUID:           connection.Player.UUID,
					Username:       connection.Player.Username,
					Properties:     []PlayerProperty{},
					Gamemode:       0,
					Ping:           0,
					HasDisplayname: true,
					DisplayName:    string(playerDisplayNameJson),
				},
			},
		})

		mc.SendToAllExcept(connection, &currentPlayersPacket)
	}

	currentPlayers := make([]Player, 0)
	for _, con := range mc.Connections {
		if con.Player == nil {
			continue
		}
		playerDisplayName := entity.ChatMessageComponent{
			Text: con.Player.Username,
		}

		playerDisplayNameJson, exception := json.Marshal(playerDisplayName)

		if exception != nil {
			fmt.Println(exception)
			continue
		}
		currentPlayers = append(currentPlayers, Player{
			UUID:           con.Player.UUID,
			Username:       con.Player.Username,
			Properties:     []PlayerProperty{},
			Gamemode:       0,
			Ping:           0,
			HasDisplayname: true,
			DisplayName:    string(playerDisplayNameJson),
		})
	}
	currentPlayersPacket := Packet(&PlayerInfoAddPlayer{
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
		fmt.Println("Spawning player ", con.Player.Username, "Entity ID ", con.Player.EntityID, " -- Our entity ID is ", connection.Player.EntityID)
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
