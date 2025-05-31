package controller

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"github.com/Ocelotworks/MinecraftGo/entity"
	packetType "github.com/Ocelotworks/MinecraftGo/packet"
	"github.com/gofrs/uuid"
)

type LoginStart struct {
	CurrentPacket *packetType.LoginStart
}

func (ls *LoginStart) GetPacketStruct() packetType.Packet {
	return &packetType.LoginStart{}
}

func (ls *LoginStart) Init(currentPacket packetType.Packet, minecraft *Minecraft) {
	ls.CurrentPacket = currentPacket.(*packetType.LoginStart)
}

func (ls *LoginStart) Handle(packet []byte, connection *Connection) {
	fmt.Println("Username ", ls.CurrentPacket.Username)

	// Initialise the player entity
	connection.Player = &entity.Player{
		Username:       ls.CurrentPacket.Username,
		Properties:     []entity.PlayerProperty{},
		Gamemode:       1,
		Ping:           0,
		HasDisplayName: true,
		DisplayName: entity.ChatMessageComponent{
			Text: ls.CurrentPacket.Username,
		},
		EntityID: connection.Minecraft.GlobalEntityCounter,
		X:        0,
		Y:        100,
		Z:        0,
		Yaw:      0,
		Pitch:    0,
	}

	// Increment the global entity counter
	connection.Minecraft.GlobalEntityCounter++

	if connection.Minecraft.EnableEncryption {
		fmt.Println(connection.Key.PublicKey)

		publicKey, exception := x509.MarshalPKIXPublicKey(&connection.Key.PublicKey)

		if exception != nil {
			fmt.Println("Marshalling public key", exception)
			return
		}

		connection.VerifyToken = make([]byte, 4)
		rand.Read(connection.VerifyToken)

		encryptionPacket := packetType.Packet(&packetType.EncryptionRequest{
			ServerID:           "",
			PublicKeyLength:    len(publicKey),
			PublicKey:          publicKey,
			VerifyTokenLength:  len(connection.VerifyToken),
			VerifyToken:        connection.VerifyToken,
			ShouldAuthenticate: true,
		})

		connection.SendPacket(&encryptionPacket)
	} else {
		connection.Player.UUID = uuid.NewV3(uuid.Nil, "OfflinePlayer:"+ls.CurrentPacket.Username).Bytes()
		connection.Minecraft.StartPlayerJoin(connection)
	}
}
